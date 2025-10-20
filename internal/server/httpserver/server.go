package httpserver

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/dontagr/metric/internal/server/config"
)

type HTTPServer struct {
	Master *echo.Echo
}

func NewServer(cfg *config.Config, log *zap.SugaredLogger, lc fx.Lifecycle, shutdowner fx.Shutdowner) (*HTTPServer, error) {
	mainServer := echo.New()

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

	var privateKey *rsa.PrivateKey
	if cfg.CryptoKey != "" {

		privateKeyPEM, err := os.ReadFile(cfg.CryptoKey)
		if err != nil {
			return nil, fmt.Errorf("privateKey ReadFile: %v", err)
		}
		privateKeyBlock, _ := pem.Decode(privateKeyPEM)
		privateKey, err = x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
		if err != nil {
			return nil, fmt.Errorf("ParsePKCS1PrivateKey: %v", err)
		}
	}

	mainServer.Use(Decrypted(privateKey))
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
			return mainServer.Shutdown(ctx)
		},
	})

	return &HTTPServer{
		Master: mainServer,
	}, nil
}
