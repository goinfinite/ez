package valueObject

import (
	"errors"
	"fmt"
	"math"

	voHelper "github.com/goinfinite/ez/src/domain/valueObject/helper"
)

type Byte int64

func NewByte(value interface{}) (Byte, error) {
	intValue, err := voHelper.InterfaceToInt64(value)
	if err != nil {
		return 0, errors.New("InvalidByte")
	}

	return Byte(intValue), nil
}

func NewMebibyte(value interface{}) (Byte, error) {
	intValue, err := voHelper.InterfaceToInt64(value)
	if err != nil {
		return 0, errors.New("InvalidMebibytes")
	}

	return Byte(intValue * 1048576), nil
}

func NewGibibyte(value interface{}) (Byte, error) {
	intValue, err := voHelper.InterfaceToInt64(value)
	if err != nil {
		return 0, errors.New("InvalidGibibytes")
	}

	return Byte(intValue * 1073741824), nil
}

func (vo Byte) Int64() int64 {
	return int64(vo)
}

func (vo Byte) Uint64() uint64 {
	return uint64(vo)
}

func (vo Byte) Float64() float64 {
	return float64(vo)
}

func (vo Byte) ToKiB() int64 {
	return int64(math.Round(vo.Float64() / 1024))
}

func (vo Byte) ToMiB() int64 {
	return int64(math.Round(vo.Float64() / 1048576))
}

func (vo Byte) ToGiB() int64 {
	return int64(math.Round(vo.Float64() / 1073741824))
}

func (vo Byte) ToTiB() int64 {
	return int64(math.Round(vo.Float64() / 1099511627776))
}

func (vo Byte) String() string {
	return fmt.Sprintf("%d", vo.Int64())
}

func (vo Byte) StringWithSuffix() string {
	voInt64 := vo.Int64()
	switch {
	case voInt64 < 1048576:
		return fmt.Sprintf("%d KiB", vo.ToKiB())
	case voInt64 < 1073741824:
		return fmt.Sprintf("%d MiB", vo.ToMiB())
	case voInt64 < 1099511627776:
		return fmt.Sprintf("%d GiB", vo.ToGiB())
	case voInt64 < 1125899906842624:
		return fmt.Sprintf("%d TiB", vo.ToTiB())
	default:
		return fmt.Sprintf("%d B", voInt64)
	}
}
