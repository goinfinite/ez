package valueObject

import (
	"errors"
	"strconv"

	voHelper "github.com/speedianet/control/src/domain/valueObject/helper"
)

type UnixGroupId uint64

func NewUnixGroupId(value interface{}) (UnixGroupId, error) {
	gid, err := voHelper.InterfaceToUint64(value)
	if err != nil {
		return 0, errors.New("InvalidGroupId")
	}

	return UnixGroupId(gid), nil
}

func (vo UnixGroupId) Uint64() uint64 {
	return uint64(vo)
}

func (vo UnixGroupId) String() string {
	return strconv.FormatUint(uint64(vo), 10)
}
