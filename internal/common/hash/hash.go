package hash

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/dontagr/metric/models"
)

func ComputeHash(key string, metric *models.Metrics) string {
	hmacHasher := hmac.New(sha256.New, []byte(key))

	hmacHasher.Write([]byte(metric.ID))
	hmacHasher.Write([]byte(metric.MType))
	if metric.Delta != nil {
		hmacHasher.Write([]byte(fmt.Sprintf("%d", *metric.Delta)))
	}
	if metric.Value != nil {
		hmacHasher.Write([]byte(fmt.Sprintf("%f", *metric.Value)))
	}

	hash := hex.EncodeToString(hmacHasher.Sum(nil))

	return hash
}

func StringHash(key string, value string) string {
	hmacHasher := hmac.New(sha256.New, []byte(key))
	hmacHasher.Write([]byte(value))
	hash := hex.EncodeToString(hmacHasher.Sum(nil))

	return hash
}
