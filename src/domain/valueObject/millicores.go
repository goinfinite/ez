package valueObject

import (
	"errors"
	"strconv"

	voHelper "github.com/goinfinite/ez/src/domain/valueObject/helper"
)

type Millicores uint

func NewMillicores(value interface{}) (Millicores, error) {
	uintValue, err := voHelper.InterfaceToUint(value)
	if err != nil {
		return 0, errors.New("InvalidMillicores")
	}

	return Millicores(uintValue), nil
}

func NewCpuCores(value interface{}) (Millicores, error) {
	floatValue, err := voHelper.InterfaceToFloat64(value)
	if err != nil {
		return 0, errors.New("InvalidCpuCores")
	}

	return Millicores(floatValue * 1000), nil
}

func (vo Millicores) ToCores() float64 {
	return float64(vo) / 1000
}

func (vo Millicores) Uint() uint {
	return uint(vo)
}

func (vo Millicores) String() string {
	return strconv.FormatUint(uint64(vo), 10)
}

func (vo Millicores) ToCoresString() string {
	return strconv.FormatFloat(vo.ToCores(), 'f', -1, 64)
}
