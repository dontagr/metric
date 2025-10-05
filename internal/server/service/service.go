package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"

	"github.com/dontagr/metric/internal/common/hash"
	"github.com/dontagr/metric/internal/server/metric/factory"
	"github.com/dontagr/metric/internal/server/metric/validator"
	"github.com/dontagr/metric/internal/server/service/backup"
	"github.com/dontagr/metric/internal/server/service/interfaces"
	serviceModels "github.com/dontagr/metric/internal/server/service/models"
	"github.com/dontagr/metric/models"
)

type Service struct {
	Store         interfaces.Store
	MetricFactory *factory.MetricFactory
	Backup        *backup.Service
	HashKey       string
}

func (s *Service) GetMetric(requestMetric serviceModels.RequestMetric) (*models.Metrics, *echo.HTTPError) {
	if !validator.IsValidMType(requestMetric.MType) {
		return nil, &echo.HTTPError{Code: http.StatusBadRequest, Message: "Invalid type"}
	}

	oldMetric, err := s.Store.LoadMetric(requestMetric.MName, requestMetric.MType)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, &echo.HTTPError{Code: http.StatusNotFound, Message: fmt.Sprintf("Not found %s %s", requestMetric.MName, requestMetric.MType)}
	} else if err != nil {
		return nil, &echo.HTTPError{Code: http.StatusInternalServerError, Message: fmt.Sprintf("загрузка не удалась для (id: %s, mtype: %s): %v", requestMetric.MName, requestMetric.MType, err)}
	}

	return oldMetric, nil
}

func (s *Service) GetStringValue(metrics *models.Metrics) (string, *echo.HTTPError) {
	metricProcessor, err := s.MetricFactory.GetMetric(metrics.MType)
	if err != nil {
		return "", &echo.HTTPError{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	return metricProcessor.ReturnValue(metrics), nil
}

func (s *Service) GetAllMetricHTML() (string, *echo.HTTPError) {
	collection, err := s.Store.ListMetric()
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return "", &echo.HTTPError{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	html := ""
	for _, metrics := range collection {
		metricProcessor, err := s.MetricFactory.GetMetric(metrics.MType)
		if err != nil {
			return "", &echo.HTTPError{Code: http.StatusInternalServerError, Message: err.Error()}
		}

		html += "<li>" + metrics.ID + ": " + metricProcessor.ReturnValue(metrics) + "</li>\n"
	}

	if html != "" {
		html = "<ul>\n" + html + "</ul>\n"
	} else {
		html = "<p>there are no metrics yet</p>"
	}

	return "<!DOCTYPE html>\n<html>\n<body>\n" + html + "</body>\n</html>", nil
}

func (s *Service) UpdateMetrics(requestArrayMetric serviceModels.RequestArrayMetric) (map[string]*models.Metrics, *echo.HTTPError) {
	metrics := make(map[string]*models.Metrics, len(requestArrayMetric))
	for _, requestMetric := range requestArrayMetric {
		var previousData *models.Metrics
		key := fmt.Sprintf("%s_%s", requestMetric.MType, requestMetric.MName)
		metric, ok := metrics[key]
		if ok {
			previousData = metric
		}

		newMetric, err := s.processUpdateData(&requestMetric, previousData)
		if err != nil {
			return nil, err
		}

		metrics[key] = newMetric
	}

	err := s.Store.BulkSaveMetric(metrics)
	if err != nil {
		return nil, &echo.HTTPError{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	s.Backup.Process()

	return metrics, nil
}

func (s *Service) UpdateMetric(requestMetric serviceModels.RequestMetric) (*models.Metrics, *echo.HTTPError) {
	newMetric, echoErr := s.processUpdateData(&requestMetric, nil)
	if echoErr != nil {
		return nil, echoErr
	}

	err := s.Store.SaveMetric(newMetric)
	if err != nil {
		return nil, &echo.HTTPError{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	s.Backup.Process()

	return newMetric, nil
}

func (s *Service) processUpdateData(requestData *serviceModels.RequestMetric, oldMetric *models.Metrics) (*models.Metrics, *echo.HTTPError) {
	if !validator.IsValidMType(requestData.MType) {
		return nil, &echo.HTTPError{Code: http.StatusBadRequest, Message: "Invalid type"}
	}

	metricProcessor, err := s.MetricFactory.GetMetric(requestData.MType)
	if err != nil {
		return nil, &echo.HTTPError{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	var newMetric *models.Metrics
	if requestData.Delta != nil {
		newMetric, err = metricProcessor.GetMetricsByData(requestData.MName, *requestData.Delta)
	} else if requestData.Value != nil {
		newMetric, err = metricProcessor.GetMetricsByData(requestData.MName, *requestData.Value)
	} else {
		newMetric, err = metricProcessor.ConvertToMetrics(requestData.MName, requestData.MValue)
	}
	if err != nil {
		return nil, &echo.HTTPError{Code: http.StatusBadRequest, Message: err.Error()}
	}
	if requestData.Hash != nil {
		newMetric.Hash = *requestData.Hash
	}

	if s.HashKey != "" {
		hashManager := hash.NewHashManager()
		hashManager.SetKey(s.HashKey)
		hashManager.SetMetrics(newMetric)
		if requestData.Hash != nil && hashManager.GetHash() != *requestData.Hash {
			return nil, &echo.HTTPError{Code: http.StatusBadRequest, Message: "Хеш не совпадает"}
		}
	}

	if oldMetric == nil {
		oldMetric, err = s.Store.LoadMetric(requestData.MName, requestData.MType)
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return nil, &echo.HTTPError{Code: http.StatusInternalServerError, Message: fmt.Sprintf("загрузка не удалась для (id: %s, mtype: %s): %v", requestData.MName, requestData.MType, err)}
		}

	}
	err = metricProcessor.Process(oldMetric, newMetric)
	if err != nil {
		return nil, &echo.HTTPError{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	if s.HashKey != "" {
		hashManager := hash.NewHashManager()
		hashManager.SetKey(s.HashKey)
		hashManager.SetMetrics(newMetric)

		newMetric.Hash = hashManager.GetHash()
	}

	return newMetric, nil
}

func (s *Service) Ping(ctx context.Context) error {
	err := s.Store.Ping(ctx)
	if err != nil {
		return err
	}

	return nil
}
