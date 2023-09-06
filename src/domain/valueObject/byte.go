package valueObject

import (
	"errors"
	"reflect"
	"strconv"
)

type Byte int64

func NewByte(value interface{}) (Byte, error) {
	var byte int64
	var err error
	switch v := value.(type) {
	case string:
		byte, err = strconv.ParseInt(v, 10, 64)
	case int, int8, int16, int32, int64:
		byte = int64(reflect.ValueOf(v).Int())
	case uint, uint8, uint16, uint32, uint64:
		byte = int64(reflect.ValueOf(v).Uint())
	case float32, float64:
		byte = int64(reflect.ValueOf(v).Float())
	default:
		err = errors.New("InvalidByte")
	}

	if err != nil {
		return 0, errors.New("InvalidByte")
	}

	return Byte(byte), nil
}

func NewBytePanic(value interface{}) Byte {
	byte, err := NewByte(value)
	if err != nil {
		panic(err)
	}
	return byte
}

func (b Byte) Get() int64 {
	return int64(b)
}

func (b Byte) ToKiB() int64 {
	return b.Get() / 1024
}

func (b Byte) ToMiB() int64 {
	return b.ToKiB() / 1024
}

func (b Byte) ToGiB() int64 {
	return b.ToMiB() / 1024
}

func (b Byte) ToTiB() int64 {
	return b.ToGiB() / 1024
}
