package valueObject

import (
	"errors"
	"reflect"
	"strconv"
)

type GroupId uint64

func NewGroupId(value interface{}) (GroupId, error) {
	var gid uint64
	var err error
	switch v := value.(type) {
	case string:
		gid, err = strconv.ParseUint(v, 10, 64)
	case int, int8, int16, int32, int64:
		intValue := reflect.ValueOf(v).Int()
		if intValue < 0 {
			err = errors.New("InvalidGroupId")
		}
		gid = uint64(intValue)
	case uint, uint8, uint16, uint32, uint64:
		gid = uint64(reflect.ValueOf(v).Uint())
	case float32, float64:
		floatValue := reflect.ValueOf(v).Float()
		if floatValue < 0 {
			err = errors.New("InvalidGroupId")
		}
		gid = uint64(floatValue)
	default:
		err = errors.New("InvalidGroupId")
	}

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
