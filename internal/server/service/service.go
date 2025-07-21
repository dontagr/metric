package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/dontagr/metric/internal/common/hash"
	"github.com/dontagr/metric/internal/server/metric/factory"
	"github.com/dontagr/metric/internal/server/service/event"
	"github.com/dontagr/metric/internal/server/service/intersaces"
	serviceModels "github.com/dontagr/metric/internal/server/service/models"
	"github.com/dontagr/metric/models"
)

type Service struct {
	MetricFactory  *factory.MetricFactory
	Store          intersaces.Store
	Event          *event.Event
	IsDirectBackup bool
	HashKey        string
}

func (s *Service) GetMetric(requestMetric serviceModels.RequestMetric) (*models.Metrics, *echo.HTTPError) {
	if !isValidMetricType(requestMetric.MType) {
		return nil, &echo.HTTPError{Code: http.StatusBadRequest, Message: "Invalid type"}
	}

	oldMetric, err := s.Store.LoadMetric(requestMetric.MName, requestMetric.MType)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, &echo.HTTPError{Code: http.StatusNotFound, Message: "Not found"}
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
	metrics := make(map[string]*models.Metrics)
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

	echoErr := s.backup()
	if echoErr != nil {
		return nil, echoErr
	}

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

	echoErr = s.backup()
	if echoErr != nil {
		return nil, echoErr
	}

	return newMetric, nil
}

func (s *Service) backup() *echo.HTTPError {
	if s.IsDirectBackup {
		metric, err := s.Store.ListMetric()
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return &echo.HTTPError{Code: http.StatusInternalServerError, Message: err.Error()}
		}
		s.Event.Metrics <- metric
	}

	return nil
}

func isValidMetricType(mType string) bool {
	if mType == "" {
		return false
	}
	if models.Counter != mType && models.Gauge != mType {
		return false
	}
	return true
}

func (s *Service) AutoBackUp(interval int, log *zap.SugaredLogger) {
	for {
		time.Sleep(time.Duration(interval) * time.Second)

		metric, err := s.Store.ListMetric()
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			log.Error(err)
		} else {
			s.Event.Metrics <- metric
		}
	}
}

func (s *Service) processUpdateData(requestData *serviceModels.RequestMetric, oldMetric *models.Metrics) (*models.Metrics, *echo.HTTPError) {
	if !isValidMetricType(requestData.MType) {
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
		computedHash := hash.ComputeHash(s.HashKey, newMetric)
		if requestData.Hash != nil && computedHash != *requestData.Hash {
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
		newMetric.Hash = hash.ComputeHash(s.HashKey, newMetric)
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
