// Модуль handler представляет собой сервисный слой где описана бизнес логика обработки входящих запросов
package handler

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/dontagr/metric/internal/common/hash"
	"github.com/dontagr/metric/internal/server/service/interfaces"
	serviceModels "github.com/dontagr/metric/internal/server/service/models"
	"github.com/dontagr/metric/models"
)

// Handler обрабатывает входящие запросы
type Handler struct {
	Service interfaces.Service // сервис для реализации обмена данных с хранилищем
	HashKey string             // ключ хэширования
}

// GetMetric получить метрику
func (h *Handler) GetMetric(c echo.Context) error {
	var requestMetric serviceModels.RequestMetric
	err := c.Bind(&requestMetric)
	if err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "Bad request"}
	}

	oldMetric, errEcho := h.Service.GetMetric(requestMetric)
	if errEcho != nil {
		return errEcho
	}

	contentType := c.Request().Header.Get(echo.HeaderContentType)
	if contentType == "application/json" {
		if h.HashKey != "" {
			c.Response().Header().Set(models.HashAlgKey, oldMetric.Hash)
		}

		return c.JSON(http.StatusOK, oldMetric)
	}

	value, errEcho := h.Service.GetStringValue(oldMetric)
	if errEcho != nil {
		return errEcho
	}
	if h.HashKey != "" {
		hashManager := hash.NewHashManager()
		hashManager.SetKey(h.HashKey)
		hashManager.SetStringValue(value)

		c.Response().Header().Set(models.HashAlgKey, hashManager.GetHash())
	}

	return c.HTML(http.StatusOK, value)
}

// GetAllMetric получить все метрики
func (h *Handler) GetAllMetric(c echo.Context) error {
	html, errEcho := h.Service.GetAllMetricHTML()
	if errEcho != nil {
		return errEcho
	}

	return c.HTML(http.StatusOK, html)
}

// UpdatesMetric обновить пакет метрик
func (h *Handler) UpdatesMetric(c echo.Context) error {
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

// UpdateMetric обновить метрику
func (h *Handler) UpdateMetric(c echo.Context) error {
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

// BadRequest обработка ошибочных запросов
func (h *Handler) BadRequest(_ echo.Context) error {
	return &echo.HTTPError{Code: http.StatusBadRequest, Message: ""}
}

// Ping just ping
func (h *Handler) Ping(c echo.Context) error {
	ctx := context.Background()

	errEcho := h.Service.Ping(ctx)
	if errEcho != nil {
		return errEcho
	}

	return c.String(http.StatusOK, "")
}
