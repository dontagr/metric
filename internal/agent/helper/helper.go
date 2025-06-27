package helper

import (
	"fmt"
	"reflect"
)

func ConvertReflectValueToInt64(val reflect.Value) (int64, error) {
	underlyingValue := val.Interface()

	switch v := underlyingValue.(type) {
	case int:
		return int64(v), nil
	case int64:
		return v, nil
	default:
		return 0, fmt.Errorf("не поддерживаемый тип: %T", v)
	}
}

func ConvertReflectValueToFloat64(val reflect.Value) (float64, error) {
	underlyingValue := val.Interface()

	switch v := underlyingValue.(type) {
	case float64:
		return v, nil
	case float32:
		return float64(v), nil
	case uint64:
		return float64(v), nil
	case uint32:
		return float64(v), nil
	default:
		return 0, fmt.Errorf("не поддерживаемый тип: %T", v)
	}
}
