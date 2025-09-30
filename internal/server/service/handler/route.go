package handler

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
	server.Master.GET("/ping", h.Ping)
	server.Master.POST("/updates/", h.UpdatesMetric)

	server.Master.GET("/debug/pprof/", h.Pprof)
	server.Master.GET("/debug/pprof/heap", h.PprofHeap)
	server.Master.GET("/debug/pprof/goroutine", h.PprofGoroutine)
	server.Master.GET("/debug/pprof/block", h.PprofBlock)
	server.Master.GET("/debug/pprof/threadcreate", h.PprofThreadCreate)
	server.Master.GET("/debug/pprof/cmdline", h.PprofCmdline)
	server.Master.GET("/debug/pprof/profile", h.PprofProfile)
	server.Master.GET("/debug/pprof/symbol", h.PprofSymbol)
	server.Master.POST("/debug/pprof/symbol", h.PprofSymbol)
	server.Master.GET("/debug/pprof/trace", h.PprofTrace)
	server.Master.GET("/debug/pprof/mutex", h.PprofMutex)
}
