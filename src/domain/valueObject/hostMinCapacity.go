package valueObject

import (
	"errors"
	"strconv"

	voHelper "github.com/speedianet/control/src/domain/valueObject/helper"
)

type HostMinCapacity uint8

func NewHostMinCapacity(value interface{}) (HostMinCapacity, error) {
	uint8Value, err := voHelper.InterfaceToUint8(value)
	if err != nil {
		return 0, errors.New("InvalidHostMinCapacity")
	}

	if uint8Value > 100 {
		return 0, errors.New("HostMinCapacityTooHigh")
	}

	return HostMinCapacity(uint8Value), nil
}

func DefaultHostMinCapacity() HostMinCapacity {
	return HostMinCapacity(20)
}

func (vo HostMinCapacity) Uint8() uint8 {
	return uint8(vo)
}

func (vo HostMinCapacity) String() string {
	return strconv.FormatUint(uint64(vo), 10)
}
