package valueObject

import (
	"errors"
	"strconv"

	voHelper "github.com/speedianet/control/src/domain/valueObject/helper"
)

type MappingId uint64

func NewMappingId(value interface{}) (MappingId, error) {
	mappingId, err := voHelper.InterfaceToUint(value)
	if err != nil {
		return 0, errors.New("InvalidMappingId")
	}

	return MappingId(mappingId), nil
}

func NewMappingIdPanic(value interface{}) MappingId {
	mappingId, err := NewMappingId(value)
	if err != nil {
		panic(err)
	}
	return mappingId
}

func (id MappingId) Get() uint64 {
	return uint64(id)
}

func (id MappingId) String() string {
	return strconv.FormatUint(uint64(id), 10)
}
