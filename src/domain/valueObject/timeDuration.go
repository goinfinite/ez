package valueObject

import (
	"errors"
	"strconv"

	voHelper "github.com/goinfinite/ez/src/domain/valueObject/helper"
)

type TimeDuration uint64

func NewTimeDuration(value interface{}) (TimeDuration, error) {
	timeDuration, err := voHelper.InterfaceToUint64(value)
	if err != nil {
		return 0, errors.New("InvalidTimeDuration")
	}

	return TimeDuration(timeDuration), nil
}

func (vo TimeDuration) Uint64() uint64 {
	return uint64(vo)
}

func (vo TimeDuration) Int64() int64 {
	return int64(vo)
}

func (vo TimeDuration) String() string {
	return strconv.FormatUint(uint64(vo), 10)
}

func (vo TimeDuration) StringWithSuffix() string {
	voUint64 := uint64(vo)
	switch {
	case voUint64 < 60:
		return strconv.FormatUint(voUint64, 10) + "s"
	case voUint64 < 3600:
		return strconv.FormatUint(voUint64/60, 10) + "m"
	case voUint64 < 86400:
		return strconv.FormatUint(voUint64/3600, 10) + "h"
	case voUint64 >= 86400:
		return strconv.FormatUint(voUint64/86400, 10) + "d"
	}

	return strconv.FormatUint(voUint64, 10) + "s"
}
