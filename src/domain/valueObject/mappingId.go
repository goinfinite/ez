package valueObject

import (
	"errors"
	"strconv"

	voHelper "github.com/speedianet/control/src/domain/valueObject/helper"
)

type MappingId uint64

func NewMappingId(value interface{}) (MappingId, error) {
	uintValue, err := voHelper.InterfaceToUint64(value)
	if err != nil {
		return 0, errors.New("InvalidMappingId")
	}

	return MappingId(uintValue), nil
}

func (vo MappingId) Uint64() uint64 {
	return uint64(vo)
}

func (vo MappingId) String() string {
	return strconv.FormatUint(uint64(vo), 10)
}
