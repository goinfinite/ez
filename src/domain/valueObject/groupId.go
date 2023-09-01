package valueObject

import (
	"errors"
	"strconv"
)

type GroupId uint64

func NewGroupId(value interface{}) (GroupId, error) {
	var gId uint64
	var err error
	switch v := value.(type) {
	case string:
		gId, err = strconv.ParseUint(v, 10, 64)
		if err != nil {
			return 0, errors.New("InvalidGroupId")
		}
	case int, int8, int16, int32, int64:
		gId = uint64(v.(int64))
	case uint, uint8, uint16, uint32, uint64:
		gId = uint64(v.(uint64))
	case float32, float64:
		gId = uint64(v.(float64))
	default:
		return 0, errors.New("InvalidGroupId")
	}

	return GroupId(gId), nil
}

func NewGroupIdPanic(value interface{}) GroupId {
	gId, err := NewGroupId(value)
	if err != nil {
		panic(err)
	}
	return gId
}

func (id GroupId) Get() uint64 {
	return uint64(id)
}

func (id GroupId) String() string {
	return strconv.FormatUint(uint64(id), 10)
}
