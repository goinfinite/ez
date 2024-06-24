package valueObject

import (
	"errors"
	"strconv"

	voHelper "github.com/speedianet/control/src/domain/valueObject/helper"
)

type MappingTargetId uint64

func NewMappingTargetId(value interface{}) (MappingTargetId, error) {
	targetId, err := voHelper.InterfaceToUint64(value)
	if err != nil {
		return 0, errors.New("InvalidMappingTargetId")
	}

	return MappingTargetId(targetId), nil
}

func NewMappingTargetIdPanic(value interface{}) MappingTargetId {
	targetId, err := NewMappingTargetId(value)
	if err != nil {
		panic(err)
	}
	return targetId
}

func (id MappingTargetId) Get() uint64 {
	return uint64(id)
}

func (id MappingTargetId) String() string {
	return strconv.FormatUint(uint64(id), 10)
}
