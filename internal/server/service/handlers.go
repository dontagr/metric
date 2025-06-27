package service

import (
	"context"
	"fmt"
	"net/http"

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
}

type requestArrayMetric []requestMetric

type requestMetric struct {
	MType  string   `param:"mType" json:"type"`
	MName  string   `param:"mName" json:"id"`
	MValue string   `param:"mValue"`
	Delta  *int64   `json:"delta,omitempty"`
	Value  *float64 `json:"value,omitempty"`
	Hash   string   `json:"hash,omitempty"`
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

	oldMetric := h.Store.LoadMetric(requestMetric.MName, requestMetric.MType)
	if oldMetric.ID == "" {
		return &echo.HTTPError{Code: http.StatusNotFound, Message: "Not found"}
	}

	metricProcessor, err := h.MetricFactory.GetMetric(requestMetric.MType)
	if err != nil {
		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	contentType := c.Request().Header.Get(echo.HeaderContentType)
	if contentType == "application/json" {
		return c.JSON(200, oldMetric)
	}

	return c.HTML(200, metricProcessor.ReturnValue(oldMetric))
}

func (h *UpdateHandler) GetAllMetric(c echo.Context) error {
	collection := h.Store.ListMetric()

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

	h.Store.BulkSaveMetric(metrics)
	if h.IsDirectBackup {
		h.Event.Metrics <- h.Store.ListMetric()
	}

	return c.JSON(200, metrics)
}

func (h *UpdateHandler) UpdateMetric(c echo.Context) error {
	var requestMetric requestMetric
	err := c.Bind(&requestMetric)
	if err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "Bad request"}
	}

	newMetric, echoError := h.processUpdateData(&requestMetric, nil)
	if echoError != nil {
		return echoError
	}

	h.Store.SaveMetric(newMetric)
	if h.IsDirectBackup {
		h.Event.Metrics <- h.Store.ListMetric()
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
