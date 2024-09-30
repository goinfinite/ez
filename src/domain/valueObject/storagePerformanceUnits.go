package valueObject

import (
	"errors"
	"strconv"

	voHelper "github.com/goinfinite/ez/src/domain/valueObject/helper"
)

type StoragePerformanceUnits uint

type StoragePerformanceUnitsLimits struct {
	ReadBytes  Byte   `json:"readBytes"`
	WriteBytes Byte   `json:"writeBytes"`
	ReadIops   uint64 `json:"readIops"`
	WriteIops  uint64 `json:"writeIops"`
}

func NewStoragePerformanceUnits(value interface{}) (StoragePerformanceUnits, error) {
	uintValue, err := voHelper.InterfaceToUint(value)
	if err != nil {
		return 0, errors.New("InvalidStoragePerformanceUnits")
	}

	return StoragePerformanceUnits(uintValue), nil
}

func (vo StoragePerformanceUnits) Uint() uint {
	return uint(vo)
}

func (vo StoragePerformanceUnits) String() string {
	return strconv.FormatUint(uint64(vo), 10)
}

func (vo StoragePerformanceUnits) ReadLimits() StoragePerformanceUnitsLimits {
	return StoragePerformanceUnitsLimits{
		ReadBytes:  Byte(int(vo) * 5000000),
		WriteBytes: Byte(int(vo) * 5000000),
		ReadIops:   uint64(vo) * 250,
		WriteIops:  uint64(vo) * 250,
	}
}
