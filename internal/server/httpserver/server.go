package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"net/netip"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/dontagr/metric/internal/server/config"
	crypro "github.com/dontagr/metric/pkg/crypto"
)

type HTTPServer struct {
	Master *echo.Echo
}

func NewServer(cfg *config.Config, cmanager *crypro.CManager, log *zap.SugaredLogger, lc fx.Lifecycle, shutdowner fx.Shutdowner) (*HTTPServer, error) {
	mainServer := echo.New()

	workerWG := sync.WaitGroup{}
	mainServer.Use(middlewareShutdowner(&workerWG))
	mainServer.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:          true,
		LogMethod:       true,
		LogStatus:       true,
		LogError:        true,
		LogResponseSize: true,
		LogLatency:      true,
		HandleError:     true,
		LogHeaders:      []string{echo.HeaderContentType, echo.HeaderContentEncoding, echo.HeaderAcceptEncoding},
		LogValuesFunc: func(_ echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				log.Infow("Request", "Method", v.Method, "URI", v.URI, "Status", v.Status, "Duration", v.Latency, "ResponseSize", v.ResponseSize, "Headers", v.Headers)
			} else {
				log.Errorw(v.Error.Error(), "Method", v.Method, "URI", v.URI, "Status", v.Status, "Duration", v.Latency, "ResponseSize", v.ResponseSize, "Headers", v.Headers)
			}

			return nil
		},
	}))

	if cfg.TrustedSubnet != "" {
		network, err := netip.ParsePrefix(cfg.TrustedSubnet)
		if err != nil {
			panic(err)
		}

		mainServer.Use(middlewareIPDefender(network))
	}

	err := cmanager.InitPrivateKey(cfg.CryptoKey)
	if err != nil {
		return nil, fmt.Errorf("failed init private key: %v", err)
	}
	mainServer.Use(middlewareDecrypted(cmanager))
	mainServer.Use(middleware.Decompress())
	mainServer.Use(middleware.Gzip())

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			log.Infof("starting HTTP server. Bind: %s", cfg.HTTPServer)
			go func() {
				if err := mainServer.Start(cfg.HTTPServer); err != nil && err != http.ErrServerClosed {
					log.Errorf("failed to start HTTP Server: %v", err)
					_ = shutdowner.Shutdown()
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Infof("Получен сигнал для завершения работы. Жду оканчания отправки.")
			workerWG.Wait()

			log.Infof("Завершение.")
			return mainServer.Shutdown(ctx)
		},
	})

	log.Infof("Строковое представление бесклассовой адресации (CIDR) = '%v'", cfg.TrustedSubnet)

	return &HTTPServer{
		Master: mainServer,
	}, nil
}
