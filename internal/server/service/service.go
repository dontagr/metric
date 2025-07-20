package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/dontagr/metric/models"
)

func isValidMetricType(mType string) bool {
	if mType == "" {
		return false
	}
	if models.Counter != mType && models.Gauge != mType {
		return false
	}
	return true
}

func (h *UpdateHandler) AutoBackUp(interval int, log *zap.SugaredLogger) {
	for {
		time.Sleep(time.Duration(interval) * time.Second)

		metric, err := h.Store.ListMetric()
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			log.Error(err)
		} else {
			h.Event.Metrics <- metric
		}
	}
}

func (h *UpdateHandler) processUpdateData(requestData *requestMetric, oldMetric *models.Metrics) (*models.Metrics, *echo.HTTPError) {
	if !isValidMetricType(requestData.MType) {
		return nil, &echo.HTTPError{Code: http.StatusBadRequest, Message: "Invalid type"}
	}

	metricProcessor, err := h.MetricFactory.GetMetric(requestData.MType)
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

	fmt.Println("=======")
	fmt.Println(h.Key)
	fmt.Println("=======")

	if h.Key != "" {
		computedHash, err := h.computeHash(newMetric)
		fmt.Println(computedHash)
		if err != nil {
			return nil, &echo.HTTPError{Code: http.StatusInternalServerError, Message: "Ошибка вычисления хеша"}
		}

		if requestData.Hash != nil && computedHash != *requestData.Hash {
			return nil, &echo.HTTPError{Code: http.StatusBadRequest, Message: "Хеш не совпадает"}
		}
	}

	if oldMetric == nil {
		oldMetric, err = h.Store.LoadMetric(requestData.MName, requestData.MType)
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return nil, &echo.HTTPError{Code: http.StatusInternalServerError, Message: fmt.Sprintf("загрузка не удалась для (id: %s, mtype: %s): %v", requestData.MName, requestData.MType, err)}
		}

	}
	err = metricProcessor.Process(oldMetric, newMetric)
	if err != nil {
		return nil, &echo.HTTPError{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	if h.Key != "" {
		computedHash, err := h.computeHash(newMetric)
		if err != nil {
			return nil, &echo.HTTPError{Code: http.StatusInternalServerError, Message: "Ошибка вычисления хеша"}
		}

		newMetric.Hash = computedHash
	}

	return newMetric, nil
}

func (h *UpdateHandler) computeHash(metric *models.Metrics) (string, error) {
	hmacHasher := hmac.New(sha256.New, []byte(h.Key))

	hmacHasher.Write([]byte(metric.ID))
	hmacHasher.Write([]byte(metric.MType))
	if metric.Delta != nil {
		hmacHasher.Write([]byte(fmt.Sprintf("%d", *metric.Delta)))
	}
	if metric.Value != nil {
		hmacHasher.Write([]byte(fmt.Sprintf("%f", *metric.Value)))
	}

	hash := hex.EncodeToString(hmacHasher.Sum(nil))

	return hash, nil
}
