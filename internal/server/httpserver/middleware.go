package httpserver

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"io"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

func Decrypted(pathPrivateKey string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if pathPrivateKey == "" {
				return next(c)
			}

			bodyBytes, err := io.ReadAll(c.Request().Body)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "cannot read request body")
			}

			privateKeyPEM, err := os.ReadFile(pathPrivateKey)
			if err != nil {
				panic(err)
			}
			privateKeyBlock, _ := pem.Decode(privateKeyPEM)
			privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
			if err != nil {
				panic(err)
			}

			encryptedBytes, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, bodyBytes, nil)
			if err != nil {
				panic(err)
			}

			c.Request().Body = io.NopCloser(bytes.NewReader(encryptedBytes))

			return next(c)
		}
	}
}
