package valueObject

import (
	"errors"
	"reflect"
	"strconv"
)

type CpuCoresCount float64

func NewCpuCoresCount(value interface{}) (CpuCoresCount, error) {
	var ccc float64
	var err error
	switch v := value.(type) {
	case string:
		ccc, err = strconv.ParseFloat(v, 64)
	case int, int8, int16, int32, int64:
		ccc = float64(reflect.ValueOf(v).Int())
	case uint, uint8, uint16, uint32, uint64:
		ccc = float64(reflect.ValueOf(v).Uint())
	case float32, float64:
		ccc = float64(reflect.ValueOf(v).Float())
	default:
		err = errors.New("InvalidCpuCoresCount")
	}

	if err != nil || ccc < 0 {
		return 0, errors.New("InvalidCpuCoresCount")
	}

	return CpuCoresCount(ccc), nil
}

func NewCpuCoresCountPanic(value interface{}) CpuCoresCount {
	ccc, err := NewCpuCoresCount(value)
	if err != nil {
		panic(err)
	}
	return ccc
}

func (ccc CpuCoresCount) Get() float64 {
	return float64(ccc)
}

func (ccc CpuCoresCount) String() string {
	return strconv.FormatFloat(float64(ccc), 'f', -1, 64)
}
