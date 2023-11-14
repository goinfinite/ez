package valueObject

import (
	"errors"
	"strconv"

	voHelper "github.com/speedianet/control/src/domain/valueObject/helper"
)

type GroupId uint64

func NewGroupId(value interface{}) (GroupId, error) {
	gid, err := voHelper.InterfaceToUint(value)
	if err != nil {
		return 0, errors.New("InvalidGroupId")
	}

	return GroupId(gid), nil
}

func NewGroupIdPanic(value interface{}) GroupId {
	gid, err := NewGroupId(value)
	if err != nil {
		panic(err)
	}
	return gid
}

func (id GroupId) Get() uint64 {
	return uint64(id)
}

func (id GroupId) String() string {
	return strconv.FormatUint(uint64(id), 10)
}
