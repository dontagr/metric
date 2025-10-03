package handler

import (
	"net/http/pprof"

	"github.com/labstack/echo/v4"
)

func (h *Handler) Pprof(c echo.Context) error {
	pprof.Index(c.Response().Writer, c.Request())
	return nil
}

func (h *Handler) PprofHeap(c echo.Context) error {
	pprof.Handler("heap").ServeHTTP(c.Response(), c.Request())
	return nil
}

func (h *Handler) PprofGoroutine(c echo.Context) error {
	pprof.Handler("goroutine").ServeHTTP(c.Response().Writer, c.Request())
	return nil
}

func (h *Handler) PprofBlock(c echo.Context) error {
	pprof.Handler("block").ServeHTTP(c.Response().Writer, c.Request())
	return nil
}

func (h *Handler) PprofThreadCreate(c echo.Context) error {
	pprof.Handler("threadcreate").ServeHTTP(c.Response().Writer, c.Request())
	return nil
}

func (h *Handler) PprofCmdline(c echo.Context) error {
	pprof.Cmdline(c.Response().Writer, c.Request())
	return nil
}

func (h *Handler) PprofProfile(c echo.Context) error {
	pprof.Profile(c.Response().Writer, c.Request())
	return nil
}

func (h *Handler) PprofSymbol(c echo.Context) error {
	pprof.Symbol(c.Response().Writer, c.Request())
	return nil
}

func (h *Handler) PprofTrace(c echo.Context) error {
	pprof.Trace(c.Response().Writer, c.Request())
	return nil
}

func (h *Handler) PprofMutex(c echo.Context) error {
	pprof.Handler("mutex").ServeHTTP(c.Response().Writer, c.Request())
	return nil
}
