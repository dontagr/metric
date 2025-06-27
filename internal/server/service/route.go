package service

import (
	"github.com/dontagr/metric/internal/server/httpserver"
)

func BindRoutes(server *httpserver.HTTPServer, h *UpdateHandler) {
	server.Master.POST("/update/:mType/:mName/:mValue", h.UpdateMetric)
	server.Master.POST("/update/:mType/:mName/:mValue/*", h.BadRequest)
	server.Master.GET("/value/:mType/:mName", h.GetMetric)
	server.Master.GET("/value/:mType/:mName/*", h.BadRequest)
	server.Master.GET("/", h.GetAllMetric)
	server.Master.POST("/update/", h.UpdateMetric)
	server.Master.POST("/value/", h.GetMetric)
}
