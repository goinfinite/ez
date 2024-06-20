package valueObject

import "time"

type UnixTime int64

func (vo UnixTime) Get() int64 {
	return time.Unix(int64(vo), 0).UTC().Unix()
}

func (vo UnixTime) GetRfcDate() string {
	return time.Unix(int64(vo), 0).UTC().Format(time.RFC3339)
}

func (vo UnixTime) GetDateOnly() string {
	return time.Unix(int64(vo), 0).UTC().Format("2006-01-02")
}

func (vo UnixTime) GetTimeOnly() string {
	return time.Unix(int64(vo), 0).UTC().Format("15:04:05")
}

func (vo UnixTime) String() string {
	return time.Unix(int64(vo), 0).UTC().String()
}
