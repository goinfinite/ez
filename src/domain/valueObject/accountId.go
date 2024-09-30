package valueObject

import (
	"errors"
	"strconv"

	voHelper "github.com/goinfinite/ez/src/domain/valueObject/helper"
)

var SystemAccountId = AccountId(0)

type AccountId uint64

func NewAccountId(value interface{}) (AccountId, error) {
	accountId, err := voHelper.InterfaceToUint64(value)
	if err != nil {
		return 0, errors.New("InvalidAccountId")
	}

	return AccountId(accountId), nil
}

func (vo AccountId) Uint64() uint64 {
	return uint64(vo)
}

func (vo AccountId) String() string {
	return strconv.FormatUint(uint64(vo), 10)
}
