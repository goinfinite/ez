package valueObject

import (
	"errors"
	"strconv"

	voHelper "github.com/speedianet/control/src/domain/valueObject/helper"
)

type HostMinCapacity float64

func NewHostMinCapacity(value interface{}) (HostMinCapacity, error) {
	floatValue, err := voHelper.InterfaceToFloat(value)
	if err != nil {
		return 0, errors.New("InvalidHostMinCapacity")
	}

	if floatValue < 0 || floatValue > 100 {
		return 0, errors.New("HostMinCapacityInvalidRange")
	}

	return HostMinCapacity(floatValue), nil
}

func DefaultHostMinCapacity() HostMinCapacity {
	return HostMinCapacity(20)
}

func (vo HostMinCapacity) Float64() float64 {
	return float64(vo)
}

func (vo HostMinCapacity) String() string {
	return strconv.FormatFloat(float64(vo), 'f', -1, 64)
}
