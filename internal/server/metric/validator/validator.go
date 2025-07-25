package validator

import "github.com/dontagr/metric/models"

func IsValidMType(mType string) bool {
	if mType == "" {
		return false
	}
	if models.Counter != mType && models.Gauge != mType {
		return false
	}
	return true
}
