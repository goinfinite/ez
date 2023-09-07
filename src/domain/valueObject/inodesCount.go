package valueObject

import (
	"errors"
	"reflect"
	"strconv"
)

type InodesCount uint64

func NewInodesCount(value interface{}) (InodesCount, error) {
	var ic uint64
	var err error
	switch v := value.(type) {
	case string:
		ic, err = strconv.ParseUint(v, 10, 64)
	case int, int8, int16, int32, int64:
		intValue := reflect.ValueOf(v).Int()
		if intValue < 0 {
			err = errors.New("InvalidInodesCount")
		}
		ic = uint64(intValue)
	case uint, uint8, uint16, uint32, uint64:
		ic = uint64(reflect.ValueOf(v).Uint())
	case float32, float64:
		floatValue := reflect.ValueOf(v).Float()
		if floatValue < 0 {
			err = errors.New("InvalidInodesCount")
		}
		ic = uint64(floatValue)
	default:
		err = errors.New("InvalidInodesCount")
	}

	if err != nil {
		return 0, errors.New("InvalidInodesCount")
	}

	return InodesCount(ic), nil
}

func NewInodesCountPanic(value interface{}) InodesCount {
	ic, err := NewInodesCount(value)
	if err != nil {
		panic(err)
	}
	return ic
}

func (ic InodesCount) Get() uint64 {
	return uint64(ic)
}

func (ic InodesCount) String() string {
	return strconv.FormatUint(uint64(ic), 10)
}
