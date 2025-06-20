package service

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/dontagr/metric/models"
)

type UpdateHandler struct {
	metricFactory *MetricFactory
	store         Store
}

func NewUpdateHandler(mf *MetricFactory, st Store) *UpdateHandler {
	return &UpdateHandler{
		metricFactory: mf,
		store:         st,
	}
}

type requestData struct {
	MType  string   `param:"mType" json:"type"`
	MName  string   `param:"mName" json:"id"`
	MValue string   `param:"mValue"`
	Delta  *int64   `json:"delta,omitempty"`
	Value  *float64 `json:"value,omitempty"`
	Hash   string   `json:"hash,omitempty"`
}

func (h *UpdateHandler) GetMetric(c echo.Context) error {
	var requestData requestData
	err := c.Bind(&requestData)
	if err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "Bad request"}
	}

	if !isValidMetricType(requestData.MType) {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "Invalid type"}
	}

	oldMetric := h.store.LoadMetric(requestData.MName, requestData.MType)
	if oldMetric.ID == "" {
		return &echo.HTTPError{Code: http.StatusNotFound, Message: "Not found"}
	}

	metricProcessor, err := h.metricFactory.GetMetric(requestData.MType)
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
	collection := h.store.ListMetric()

	html := ""
	for _, metrics := range collection {
		metricProcessor, err := h.metricFactory.GetMetric(metrics.MType)
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

func (h *UpdateHandler) UpdateMetric(c echo.Context) error {
	var requestData requestData
	err := c.Bind(&requestData)
	if err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "Bad request"}
	}
	if !isValidMetricType(requestData.MType) {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "Invalid type"}
	}

	metricProcessor, err := h.metricFactory.GetMetric(requestData.MType)
	if err != nil {
		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	var newMetric *models.Metrics
	if requestData.Delta != nil {
		newMetric, err = metricProcessor.GetMetricsByData(requestData.MName, *requestData.Delta)
	} else if requestData.Value != nil {
		newMetric, err = metricProcessor.GetMetricsByData(requestData.MName, *requestData.Value)
	} else {
		newMetric, err = metricProcessor.ConvertToMetrics(requestData.MName, requestData.MValue)
		if err != nil {
			return &echo.HTTPError{Code: http.StatusBadRequest, Message: err.Error()}
		}
	}

	oldMetric := h.store.LoadMetric(requestData.MName, requestData.MType)
	err = metricProcessor.Process(oldMetric, newMetric)
	if err != nil {
		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	h.store.SaveMetric(newMetric)

	contentType := c.Request().Header.Get(echo.HeaderContentType)
	if contentType == "application/json" {
		return c.JSON(200, newMetric)
	}

	logs := fmt.Sprintf("POST %s %s %s\n", newMetric.ID, newMetric.MType, metricProcessor.ReturnValue(newMetric))

	return c.HTML(200, logs)
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

func (h *UpdateHandler) BadRequest(_ echo.Context) error {
	return &echo.HTTPError{Code: http.StatusBadRequest, Message: ""}
}
