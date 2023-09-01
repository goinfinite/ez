package valueObject

import (
	"errors"
	"strconv"
)

type AccountId uint64

func NewAccountId(value interface{}) (AccountId, error) {
	var accId uint64
	var err error
	switch v := value.(type) {
	case string:
		accId, err = strconv.ParseUint(v, 10, 64)
		if err != nil {
			return 0, errors.New("InvalidAccountId")
		}
	case int, int8, int16, int32, int64:
		accId = uint64(v.(int64))
	case uint, uint8, uint16, uint32, uint64:
		accId = uint64(v.(uint64))
	case float32, float64:
		accId = uint64(v.(float64))
	default:
		return 0, errors.New("InvalidAccountId")
	}

	return AccountId(accId), nil
}

func NewAccountIdPanic(value interface{}) AccountId {
	accId, err := NewAccountId(value)
	if err != nil {
		panic(err)
	}
	return accId
}

func (id AccountId) Get() uint64 {
	return uint64(id)
}

func (id AccountId) String() string {
	return strconv.FormatUint(uint64(id), 10)
}
