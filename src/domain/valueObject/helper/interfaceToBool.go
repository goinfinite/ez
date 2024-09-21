package voHelper

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
)

func InterfaceToBool(input interface{}) (output bool, err error) {
	defaultErr := errors.New("InvalidBool")
	switch v := input.(type) {
	case bool:
		output = v
	case string:
		output, err = strconv.ParseBool(v)
		if err != nil {
			inputStr := strings.ToLower(v)
			output = inputStr == "on" || inputStr == "yes" || inputStr == "y"
			err = nil
		}
	case int, int8, int16, int32, int64:
		intValue := reflect.ValueOf(v).Int()
		output = intValue != 0
	case uint, uint8, uint16, uint32, uint64:
		uintValue := reflect.ValueOf(v).Uint()
		output = uintValue != 0
	case float32, float64:
		floatValue := reflect.ValueOf(v).Float()
		output = floatValue != 0
	default:
		err = defaultErr
	}

	if err != nil {
		return false, defaultErr
	}

	return output, nil
}
