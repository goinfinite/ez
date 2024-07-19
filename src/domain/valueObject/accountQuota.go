package valueObject

import (
	"errors"
	"strconv"
	"strings"
)

type AccountQuota struct {
	Millicores              Millicores              `json:"millicores"`
	MemoryBytes             Byte                    `json:"memoryBytes"`
	StorageBytes            Byte                    `json:"storageBytes"`
	StorageInodes           uint64                  `json:"storageInodes"`
	StoragePerformanceUnits StoragePerformanceUnits `json:"storagePerformanceUnits"`
}

func NewAccountQuota(
	millicores Millicores,
	memoryBytes, storageBytes Byte,
	storageInodes uint64,
	storagePerformanceUnits StoragePerformanceUnits,
) AccountQuota {
	return AccountQuota{
		Millicores:              millicores,
		MemoryBytes:             memoryBytes,
		StorageBytes:            storageBytes,
		StorageInodes:           storageInodes,
		StoragePerformanceUnits: storagePerformanceUnits,
	}
}

func NewAccountQuotaFromString(value string) (quota AccountQuota, err error) {
	if value == "" {
		return quota, errors.New("InvalidAccountQuotaValue")
	}

	if !strings.Contains(value, ":") {
		return quota, errors.New("InvalidAccountQuotaFormat")
	}

	quotaParts := strings.Split(value, ":")
	if len(quotaParts) != 5 {
		return quota, errors.New("InvalidAccountQuotaStructure")
	}

	millicores, err := NewMillicores(quotaParts[0])
	if err != nil {
		return quota, err
	}

	memoryBytes, err := NewByte(quotaParts[1])
	if err != nil {
		return quota, err
	}

	storageBytes, err := NewByte(quotaParts[2])
	if err != nil {
		return quota, err
	}

	storageInodes, err := strconv.ParseUint(quotaParts[3], 10, 64)
	if err != nil {
		return quota, errors.New("InvalidAccountQuotaInodes")
	}

	storagePerformanceUnits, _ := NewStoragePerformanceUnits(1)
	if len(quotaParts) == 5 {
		storagePerformanceUnits, err = NewStoragePerformanceUnits(quotaParts[4])
		if err != nil {
			return quota, err
		}
	}

	return NewAccountQuota(
		millicores, memoryBytes, storageBytes, storageInodes, storagePerformanceUnits,
	), nil
}

func NewAccountQuotaWithDefaultValues() AccountQuota {
	millicores, _ := NewMillicores("500")
	memoryBytes, _ := NewByte("1073741824")
	storageBytes, _ := NewByte("5368709120")
	storageInodes := uint64(500000)
	storagePerformanceUnits, _ := NewStoragePerformanceUnits("1")

	return NewAccountQuota(
		millicores, memoryBytes, storageBytes, storageInodes, storagePerformanceUnits,
	)
}

func NewAccountQuotaWithBlankValues() AccountQuota {
	millicores, _ := NewMillicores("0")
	memoryBytes, _ := NewByte("0")
	storageBytes, _ := NewByte("0")
	storageInodes := uint64(0)
	storagePerformanceUnits, _ := NewStoragePerformanceUnits("0")

	return NewAccountQuota(
		millicores, memoryBytes, storageBytes, storageInodes, storagePerformanceUnits,
	)
}
