package valueObject

import (
	"errors"
	"strconv"
	"time"

	voHelper "github.com/goinfinite/ez/src/domain/valueObject/helper"
)

type UnixTime int64

func NewUnixTime(value interface{}) (UnixTime, error) {
	unixTime, err := voHelper.InterfaceToInt64(value)
	if err != nil {
		return 0, errors.New("InvalidUnixTime")
	}

	return UnixTime(unixTime), nil
}

func NewUnixTimeNow() UnixTime {
	return UnixTime(time.Now().Unix())
}

func NewUnixTimeBeforeNow(duration time.Duration) UnixTime {
	return UnixTime(time.Now().Add(-duration).Unix())
}

func NewUnixTimeAfterNow(duration time.Duration) UnixTime {
	return UnixTime(time.Now().Add(duration).Unix())
}

func NewUnixTimeWithGoTime(goTime time.Time) UnixTime {
	return UnixTime(goTime.Unix())
}

func (vo UnixTime) Read() int64 {
	return time.Unix(int64(vo), 0).UTC().Unix()
}

func (vo UnixTime) ReadAsRfcDate() string {
	return time.Unix(int64(vo), 0).UTC().Format(time.RFC3339)
}

func (vo UnixTime) GetDateOnly() string {
	return time.Unix(int64(vo), 0).UTC().Format("2006-01-02")
}

func (vo UnixTime) GetTimeOnly() string {
	return time.Unix(int64(vo), 0).UTC().Format("15:04:05")
}

func (vo UnixTime) GetAsGoTime() time.Time {
	return time.Unix(int64(vo), 0).UTC()
}

func (vo UnixTime) String() string {
	return strconv.FormatInt(int64(vo), 10)
}
