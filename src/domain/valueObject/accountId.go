package valueObject

import (
	"errors"
	"strconv"
)

type AccountId uint64

func NewAccountId(value interface{}) (AccountId, error) {
	var accId uint64
	var err error
	switch value := value.(type) {
	case string:
		accId, err = strconv.ParseUint(value, 10, 64)
	case float64:
		accId = uint64(value)
	case int64:
		accId = uint64(value)
	case uint64:
		accId = value
	default:
		return 0, errors.New("InvalidAccountId")
	}

	if err != nil {
		return 0, errors.New("InvalidAccountId")
	}

	return AccountId(accId), nil
}

func NewAccountIdPanic(value string) AccountId {
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
