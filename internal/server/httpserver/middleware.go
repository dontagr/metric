package httpserver

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"hash"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Decrypted(privateKey *rsa.PrivateKey) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if privateKey == nil {
				return next(c)
			}

			bodyBytes, err := io.ReadAll(c.Request().Body)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("cannot read request body: %v", err))
			}

			data, err := base64.StdEncoding.DecodeString(string(bodyBytes))
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("cannot decode base64: %v", err))
			}

			encryptedBytes, err := decryptOAEP(sha256.New(), rand.Reader, privateKey, data, nil)
			if err != nil {
				panic(err)
			}

			c.Request().Body = io.NopCloser(bytes.NewReader(encryptedBytes))

			return next(c)
		}
	}
}

func decryptOAEP(hash hash.Hash, random io.Reader, private *rsa.PrivateKey, msg []byte, label []byte) ([]byte, error) {
	msgLen := len(msg)
	step := private.PublicKey.Size()
	var decryptedBytes []byte

	for start := 0; start < msgLen; start += step {
		finish := start + step
		if finish > msgLen {
			finish = msgLen
		}

		decryptedBlockBytes, err := rsa.DecryptOAEP(hash, random, private, msg[start:finish], label)
		if err != nil {
			return nil, err
		}

		decryptedBytes = append(decryptedBytes, decryptedBlockBytes...)
	}

	return decryptedBytes, nil
}
