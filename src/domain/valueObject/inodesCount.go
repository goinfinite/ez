package valueObject

import (
	"errors"
	"strconv"

	voHelper "github.com/speedianet/control/src/domain/valueObject/helper"
)

type InodesCount uint64

func NewInodesCount(value interface{}) (InodesCount, error) {
	ic, err := voHelper.InterfaceToUint(value)
	if err != nil {
		return 0, errors.New("InvalidInodesCount")
	}

	return InodesCount(ic), nil
}

func NewInodesCountPanic(value interface{}) InodesCount {
	ic, err := NewInodesCount(value)
	if err != nil {
		panic(err)
	}
	return ic
}

func (ic InodesCount) Get() uint64 {
	return uint64(ic)
}

func (ic InodesCount) String() string {
	return strconv.FormatUint(uint64(ic), 10)
}
