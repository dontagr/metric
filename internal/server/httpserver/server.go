package httpserver

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"go.uber.org/fx"

	"github.com/dontagr/metric/internal/server/config"
)

type HttpServer struct {
	Master *http.Server
}

func NewServer(cfg *config.Config, mux *http.ServeMux, lc fx.Lifecycle, shutdowner fx.Shutdowner) *HttpServer {
	mainServer := &http.Server{Addr: cfg.HttpServing.BindAddress, Handler: mux}

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			ln, err := net.Listen("tcp", mainServer.Addr)
			if err != nil {
				return err
			}
			fmt.Printf("Starting HTTP server. Bind: %s\n", mainServer.Addr)
			go func() {
				err := mainServer.Serve(ln)
				if err != nil {
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

	return &HttpServer{
		Master: mainServer,
	}
}
