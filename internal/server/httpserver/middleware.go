package httpserver

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"

	crypro "github.com/dontagr/metric/pkg/crypto"
)

func middlewareShutdowner(wg *sync.WaitGroup) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			defer wg.Done()
			wg.Add(1)

			return next(c)
		}
	}

}

func middlewareDecrypted(cmanager *crypro.CManager) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if cmanager.PrivateKey == nil {
				return next(c)
			}

			bodyBytes, err := io.ReadAll(c.Request().Body)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("cannot read request body: %v", err))
			}

			encryptedBytes, err := cmanager.Decrypt(bodyBytes)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("cannot decode base64: %v", err))
			}

			c.Request().Body = io.NopCloser(bytes.NewReader(encryptedBytes))

			return next(c)
		}
	}
}
