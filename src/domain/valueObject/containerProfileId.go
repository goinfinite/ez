package valueObject

import (
	"errors"
	"strconv"

	voHelper "github.com/speedianet/control/src/domain/valueObject/helper"
)

type ContainerProfileId uint64

func NewContainerProfileId(value interface{}) (ContainerProfileId, error) {
	uintValue, err := voHelper.InterfaceToUint64(value)
	if err != nil {
		return 0, errors.New("InvalidContainerProfileId")
	}

	return ContainerProfileId(uintValue), nil
}

func (vo ContainerProfileId) Uint64() uint64 {
	return uint64(vo)
}

func (vo ContainerProfileId) String() string {
	return strconv.FormatUint(uint64(vo), 10)
}
