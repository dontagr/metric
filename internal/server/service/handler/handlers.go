package handler

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/dontagr/metric/internal/common/hash"
	"github.com/dontagr/metric/internal/server/service/intersaces"
	serviceModels "github.com/dontagr/metric/internal/server/service/models"
	"github.com/dontagr/metric/models"
)

type UpdateHandler struct {
	Service intersaces.Service
	HashKey string
}

func (h *UpdateHandler) GetMetric(c echo.Context) error {
	var requestMetric serviceModels.RequestMetric
	err := c.Bind(&requestMetric)
	if err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "Bad request"}
	}

	oldMetric, err := h.Service.GetMetric(requestMetric)
	if err != nil {
		return err
	}

	contentType := c.Request().Header.Get(echo.HeaderContentType)
	if contentType == "application/json" {
		if h.HashKey != "" {
			c.Response().Header().Set(models.HashAlgKey, oldMetric.Hash)
		}

		return c.JSON(http.StatusOK, oldMetric)
	}

	value, err := h.Service.GetStringValue(oldMetric)
	if err != nil {
		return err
	}
	if h.HashKey != "" {
		c.Response().Header().Set(models.HashAlgKey, hash.StringHash(h.HashKey, value))
	}

	return c.HTML(http.StatusOK, value)
}

func (h *UpdateHandler) GetAllMetric(c echo.Context) error {
	html, err := h.Service.GetAllMetricHTML()
	if err != nil {
		return err
	}

	return c.HTML(http.StatusOK, html)
}

func (h *UpdateHandler) UpdatesMetric(c echo.Context) error {
	var requestArrayMetric serviceModels.RequestArrayMetric
	err := c.Bind(&requestArrayMetric)
	if err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "Bad request"}
	}

	metrics, errEcho := h.Service.UpdateMetrics(requestArrayMetric)
	if errEcho != nil {
		return errEcho
	}

	return c.JSON(http.StatusOK, metrics)
}

func (h *UpdateHandler) UpdateMetric(c echo.Context) error {
	var requestMetric serviceModels.RequestMetric
	err := c.Bind(&requestMetric)
	if err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "Bad request"}
	}

	metric, errEcho := h.Service.UpdateMetric(requestMetric)
	if errEcho != nil {
		return errEcho
	}

	contentType := c.Request().Header.Get(echo.HeaderContentType)
	if contentType == "application/json" {
		return c.JSON(http.StatusOK, metric)
	}

	return c.String(http.StatusOK, "")
}

func (h *UpdateHandler) BadRequest(_ echo.Context) error {
	return &echo.HTTPError{Code: http.StatusBadRequest, Message: ""}
}

func (h *UpdateHandler) Ping(c echo.Context) error {
	ctx := context.Background()

	err := h.Service.Ping(ctx)
	if err != nil {
		return err
	}

	return c.String(http.StatusOK, "")
}
