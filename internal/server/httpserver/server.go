package httpserver

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"

	"github.com/dontagr/metric/internal/server/config"
)

type HTTPServer struct {
	Master *echo.Echo
}

func NewServer(cfg *config.Config, lc fx.Lifecycle, shutdowner fx.Shutdowner) *HTTPServer {
	mainServer := echo.New()

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			fmt.Printf("Starting HTTP server. Bind: %s\n", cfg.HTTPServer.BindAddress)
			go func() {
				if err := mainServer.Start(cfg.HTTPServer.BindAddress); err != nil && err != http.ErrServerClosed {
					fmt.Println("Failed to HTTP Serve")
					fmt.Println(err)
					_ = shutdowner.Shutdown()
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return mainServer.Shutdown(ctx)
		},
	})

	return &HTTPServer{
		Master: mainServer,
	}
}
