package voHelper

import (
	"errors"
	"reflect"
	"strconv"
)

func InterfaceToFloat(input interface{}) (float64, error) {
	var output float64
	var err error
	var defaultErr error = errors.New("InvalidInput")
	switch v := input.(type) {
	case string:
		output, err = strconv.ParseFloat(v, 64)
	case int, int8, int16, int32, int64:
		output = float64(reflect.ValueOf(v).Int())
	case uint, uint8, uint16, uint32, uint64:
		output = float64(reflect.ValueOf(v).Uint())
	case float32, float64:
		output = float64(reflect.ValueOf(v).Float())
	default:
		err = defaultErr
	}

	if err != nil {
		return float64(0), defaultErr
	}

	return output, nil
}
