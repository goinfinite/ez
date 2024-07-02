package valueObject

import (
	"errors"
	"strconv"

	voHelper "github.com/speedianet/control/src/domain/valueObject/helper"
)

type CpuCoresCount float64

func NewCpuCoresCount(value interface{}) (CpuCoresCount, error) {
	ccc, err := voHelper.InterfaceToFloat(value)
	if err != nil || ccc < 0 {
		return 0, errors.New("InvalidCpuCoresCount")
	}

	return CpuCoresCount(ccc), nil
}

func NewCpuCoresCountPanic(value interface{}) CpuCoresCount {
	ccc, err := NewCpuCoresCount(value)
	if err != nil {
		panic(err)
	}
	return ccc
}

func (ccc CpuCoresCount) Read() float64 {
	return float64(ccc)
}

func (ccc CpuCoresCount) String() string {
	return strconv.FormatFloat(float64(ccc), 'f', -1, 64)
}
