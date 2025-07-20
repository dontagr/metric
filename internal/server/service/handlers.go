package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"

	"github.com/dontagr/metric/internal/server/service/event"
	"github.com/dontagr/metric/internal/server/service/intersaces"
	"github.com/dontagr/metric/models"
)

type UpdateHandler struct {
	MetricFactory  *MetricFactory
	Store          intersaces.Store
	Event          *event.Event
	IsDirectBackup bool
	Key            string
}

type requestArrayMetric []requestMetric

type requestMetric struct {
	MType  string   `param:"mType" json:"type"`
	MName  string   `param:"mName" json:"id"`
	MValue string   `param:"mValue"`
	Delta  *int64   `json:"delta,omitempty"`
	Value  *float64 `json:"value,omitempty"`
	Hash   *string  `json:"hash,omitempty"`
}

func (h *UpdateHandler) GetMetric(c echo.Context) error {
	var requestMetric requestMetric
	err := c.Bind(&requestMetric)
	if err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "Bad request"}
	}

	if !isValidMetricType(requestMetric.MType) {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "Invalid type"}
	}

	oldMetric, err := h.Store.LoadMetric(requestMetric.MName, requestMetric.MType)
	if errors.Is(err, pgx.ErrNoRows) {
		return &echo.HTTPError{Code: http.StatusNotFound, Message: "Not found"}
	} else if err != nil {
		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: fmt.Sprintf("загрузка не удалась для (id: %s, mtype: %s): %v", requestMetric.MName, requestMetric.MType, err)}
	}

	metricProcessor, err := h.MetricFactory.GetMetric(requestMetric.MType)
	if err != nil {
		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	if h.Key != "" {
		c.Response().Header().Set("HashSHA256", oldMetric.Hash)
	}

	contentType := c.Request().Header.Get(echo.HeaderContentType)
	if contentType == "application/json" {
		return c.JSON(200, oldMetric)
	}

	return c.HTML(200, metricProcessor.ReturnValue(oldMetric))
}

func (h *UpdateHandler) GetAllMetric(c echo.Context) error {
	collection, err := h.Store.ListMetric()
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	html := ""
	for _, metrics := range collection {
		metricProcessor, err := h.MetricFactory.GetMetric(metrics.MType)
		if err != nil {
			return &echo.HTTPError{Code: http.StatusInternalServerError, Message: err.Error()}
		}

		html += "<li>" + metrics.ID + ": " + metricProcessor.ReturnValue(metrics) + "</li>\n"
	}

	if html != "" {
		html = "<ul>\n" + html + "</ul>\n"
	} else {
		html = "<p>there are no metrics yet</p>"
	}

	return c.HTML(200, "<!DOCTYPE html>\n<html>\n<body>\n"+html+"</body>\n</html>")
}

func (h *UpdateHandler) UpdatesMetric(c echo.Context) error {
	//bodyBytes, err1 := ioutil.ReadAll(c.Request().Body)
	//if err1 != nil {
	//	return &echo.HTTPError{Code: http.StatusInternalServerError, Message: "Ошибка чтения тела запроса"}
	//}
	//// Вывод прочитанных данных в лог (для отладки)
	//fmt.Println("Прочитанное тело запроса:", string(bodyBytes))
	//
	//// Сбрасываем ридер до начала, чтобы последующий вызов c.Bind мог снова прочитать данные.
	//c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	var requestArrayMetric requestArrayMetric
	err := c.Bind(&requestArrayMetric)
	if err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "Bad request"}
	}

	metrics := make(map[string]*models.Metrics)
	for _, requestMetric := range requestArrayMetric {
		var previousData *models.Metrics
		key := fmt.Sprintf("%s_%s", requestMetric.MType, requestMetric.MName)
		metric, ok := metrics[key]
		if ok {
			previousData = metric
		}

		newMetric, echoError := h.processUpdateData(&requestMetric, previousData)
		if echoError != nil {
			return echoError
		}

		metrics[key] = newMetric
	}

	err = h.Store.BulkSaveMetric(metrics)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	if h.IsDirectBackup {
		metric, err := h.Store.ListMetric()
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		h.Event.Metrics <- metric
	}

	return c.JSON(200, metrics)
}

func (h *UpdateHandler) UpdateMetric(c echo.Context) error {
	//bodyBytes, err1 := ioutil.ReadAll(c.Request().Body)
	//if err1 != nil {
	//	return &echo.HTTPError{Code: http.StatusInternalServerError, Message: "Ошибка чтения тела запроса"}
	//}
	//// Вывод прочитанных данных в лог (для отладки)
	//fmt.Println("Прочитанное тело запроса:", string(bodyBytes))
	//
	//// Сбрасываем ридер до начала, чтобы последующий вызов c.Bind мог снова прочитать данные.
	//c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	var requestMetric requestMetric
	err := c.Bind(&requestMetric)
	if err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "Bad request"}
	}

	newMetric, echoError := h.processUpdateData(&requestMetric, nil)
	if echoError != nil {
		return echoError
	}

	err = h.Store.SaveMetric(newMetric)
	if err != nil {
		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: err.Error()}
	}
	if h.IsDirectBackup {
		metric, err := h.Store.ListMetric()
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return &echo.HTTPError{Code: http.StatusInternalServerError, Message: err.Error()}
		}
		h.Event.Metrics <- metric
	}

	contentType := c.Request().Header.Get(echo.HeaderContentType)
	if contentType == "application/json" {
		return c.JSON(200, newMetric)
	}

	return c.String(200, "")
}

func (h *UpdateHandler) BadRequest(_ echo.Context) error {
	return &echo.HTTPError{Code: http.StatusBadRequest, Message: ""}
}

func (h *UpdateHandler) Ping(c echo.Context) error {
	ctx := context.Background()

	err := h.Store.Ping(ctx)
	if err != nil {
		return err
	}

	return c.String(http.StatusOK, "")
}
