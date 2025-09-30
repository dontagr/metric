package handler

import (
	"net/http/pprof"

	"github.com/labstack/echo/v4"
)

func (h *UpdateHandler) Pprof(c echo.Context) error {
	pprof.Index(c.Response().Writer, c.Request())
	return nil
}

func (h *UpdateHandler) PprofHeap(c echo.Context) error {
	pprof.Handler("heap").ServeHTTP(c.Response(), c.Request())
	return nil
}

func (h *UpdateHandler) PprofGoroutine(c echo.Context) error {
	pprof.Handler("goroutine").ServeHTTP(c.Response().Writer, c.Request())
	return nil
}

func (h *UpdateHandler) PprofBlock(c echo.Context) error {
	pprof.Handler("block").ServeHTTP(c.Response().Writer, c.Request())
	return nil
}

func (h *UpdateHandler) PprofThreadCreate(c echo.Context) error {
	pprof.Handler("threadcreate").ServeHTTP(c.Response().Writer, c.Request())
	return nil
}

func (h *UpdateHandler) PprofCmdline(c echo.Context) error {
	pprof.Cmdline(c.Response().Writer, c.Request())
	return nil
}

func (h *UpdateHandler) PprofProfile(c echo.Context) error {
	pprof.Profile(c.Response().Writer, c.Request())
	return nil
}

func (h *UpdateHandler) PprofSymbol(c echo.Context) error {
	pprof.Symbol(c.Response().Writer, c.Request())
	return nil
}

func (h *UpdateHandler) PprofTrace(c echo.Context) error {
	pprof.Trace(c.Response().Writer, c.Request())
	return nil
}

func (h *UpdateHandler) PprofMutex(c echo.Context) error {
	pprof.Handler("mutex").ServeHTTP(c.Response().Writer, c.Request())
	return nil
}
