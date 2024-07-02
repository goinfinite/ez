package valueObject

import (
	"errors"
	"strconv"

	voHelper "github.com/speedianet/control/src/domain/valueObject/helper"
)

type ContainerProfileId uint64

func NewContainerProfileId(value interface{}) (ContainerProfileId, error) {
	rpId, err := voHelper.InterfaceToUint64(value)
	if err != nil {
		return 0, errors.New("InvalidContainerProfileId")
	}

	return ContainerProfileId(rpId), nil
}

func NewContainerProfileIdPanic(value interface{}) ContainerProfileId {
	rpId, err := NewContainerProfileId(value)
	if err != nil {
		panic(err)
	}
	return rpId
}

func (id ContainerProfileId) Read() uint64 {
	return uint64(id)
}

func (id ContainerProfileId) String() string {
	return strconv.FormatUint(uint64(id), 10)
}
