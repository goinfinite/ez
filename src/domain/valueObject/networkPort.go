package valueObject

import (
	"errors"
	"reflect"
	"strconv"
)

type NetworkPort uint64

func NewNetworkPort(value interface{}) (NetworkPort, error) {
	var np uint64
	var err error
	switch v := value.(type) {
	case string:
		np, err = strconv.ParseUint(v, 10, 64)
	case int, int8, int16, int32, int64:
		intValue := reflect.ValueOf(v).Int()
		if intValue < 0 {
			err = errors.New("InvalidNetworkPort")
		}
		np = uint64(intValue)
	case uint, uint8, uint16, uint32, uint64:
		np = uint64(reflect.ValueOf(v).Uint())
	case float32, float64:
		floatValue := reflect.ValueOf(v).Float()
		if floatValue < 0 {
			err = errors.New("InvalidNetworkPort")
		}
		np = uint64(floatValue)
	default:
		err = errors.New("InvalidNetworkPort")
	}

	if err != nil {
		return 0, errors.New("InvalidNetworkPort")
	}

	return NetworkPort(np), nil
}

func NewNetworkPortPanic(value interface{}) NetworkPort {
	np, err := NewNetworkPort(value)
	if err != nil {
		panic(err)
	}
	return np
}

func (np NetworkPort) Get() uint64 {
	return uint64(np)
}

func (np NetworkPort) String() string {
	return strconv.FormatUint(uint64(np), 10)
}
