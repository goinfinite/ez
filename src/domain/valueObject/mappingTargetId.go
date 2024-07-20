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

func (vo MappingTargetId) Uint64() uint64 {
	return uint64(vo)
}

func (vo MappingTargetId) String() string {
	return strconv.FormatUint(uint64(vo), 10)
}
