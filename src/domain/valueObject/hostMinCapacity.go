package valueObject

import (
	"errors"
	"strconv"

	voHelper "github.com/speedianet/control/src/domain/valueObject/helper"
)

type HostMinCapacity float64

func NewHostMinCapacity(value interface{}) (HostMinCapacity, error) {
	hmc, err := voHelper.InterfaceToFloat(value)
	if err != nil {
		return 0, errors.New("InvalidHostMinCapacity")
	}

	if hmc < 0 || hmc > 100 {
		return 0, errors.New("HostMinCapacityInvalidRange")
	}

	return HostMinCapacity(hmc), nil
}

func DefaultHostMinCapacity() HostMinCapacity {
	return HostMinCapacity(20)
}

func NewHostMinCapacityPanic(value interface{}) HostMinCapacity {
	hmc, err := NewHostMinCapacity(value)
	if err != nil {
		panic(err)
	}
	return hmc
}

func (hmc HostMinCapacity) Get() float64 {
	return float64(hmc)
}

func (hmc HostMinCapacity) String() string {
	return strconv.FormatFloat(float64(hmc), 'f', -1, 64)
}
