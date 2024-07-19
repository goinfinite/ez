package valueObject

import (
	"errors"
	"strconv"

	voHelper "github.com/speedianet/control/src/domain/valueObject/helper"
)

type Millicores uint

func NewMillicores(value interface{}) (Millicores, error) {
	uintValue, err := voHelper.InterfaceToUint(value)
	if err != nil {
		return 0, errors.New("InvalidMillicores")
	}

	return Millicores(uintValue), nil
}

func (vo Millicores) ReadAsCores() float64 {
	return float64(vo) / 1000
}

func (vo Millicores) Uint() uint {
	return uint(vo)
}

func (vo Millicores) String() string {
	return strconv.FormatUint(uint64(vo), 10)
}
