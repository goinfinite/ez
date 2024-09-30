package valueObject

import (
	"errors"
	"strconv"
	"strings"

	voHelper "github.com/goinfinite/ez/src/domain/valueObject/helper"
)

type AccountQuota struct {
	Millicores              Millicores              `json:"millicores"`
	CpuCores                float64                 `json:"cpuCores"`
	MemoryBytes             Byte                    `json:"memoryBytes"`
	MemoryMebibytes         int64                   `json:"memoryMebibytes"`
	MemoryGibibytes         int64                   `json:"memoryGibibytes"`
	StorageBytes            Byte                    `json:"storageBytes"`
	StorageMebibytes        int64                   `json:"storageMebibytes"`
	StorageGibibytes        int64                   `json:"storageGibibytes"`
	StorageInodes           uint64                  `json:"storageInodes"`
	StoragePerformanceUnits StoragePerformanceUnits `json:"storagePerformanceUnits"`
}

func NewAccountQuota(
	millicores Millicores,
	memoryBytes, storageBytes Byte,
	storageInodes uint64,
	storagePerformanceUnits StoragePerformanceUnits,
) AccountQuota {
	cpuCores := millicores.ReadAsCores()
	memoryMebibytes := memoryBytes.ToMiB()
	memoryGibibytes := memoryBytes.ToGiB()
	storageMebibytes := storageBytes.ToMiB()
	storageGibibytes := storageBytes.ToGiB()

	return AccountQuota{
		Millicores:              millicores,
		CpuCores:                cpuCores,
		MemoryBytes:             memoryBytes,
		MemoryMebibytes:         memoryMebibytes,
		MemoryGibibytes:         memoryGibibytes,
		StorageBytes:            storageBytes,
		StorageMebibytes:        storageMebibytes,
		StorageGibibytes:        storageGibibytes,
		StorageInodes:           storageInodes,
		StoragePerformanceUnits: storagePerformanceUnits,
	}
}

// format: [millicores][:memoryBytes][:storageBytes][:storageInodes][:storagePerformanceUnits]
func NewAccountQuotaFromString(value string) (quota AccountQuota, err error) {
	value = strings.TrimSpace(value)
	value = strings.ToLower(value)

	quotaRegex := `^(?:(?P<millicores>\d{1,19}))?(?::?(?P<memoryBytes>\d{1,19}))?(?::(?P<storageBytes>\d{1,19}))?(?::(?P<storageInodes>\d{1,19}))?(?::(?P<storagePerformanceUnits>\d{1,19}))?$`
	quotaParts := voHelper.FindNamedGroupsMatches(quotaRegex, value)
	if len(quotaParts) == 0 {
		return quota, errors.New("InvalidQuotaStructure")
	}

	quota = NewAccountQuotaWithDefaultValues()

	if quotaParts["millicores"] != "" {
		quota.Millicores, err = NewMillicores(quotaParts["millicores"])
		if err != nil {
			return quota, err
		}
	}

	if quotaParts["memoryBytes"] == "" {
		quota.MemoryBytes, err = NewByte(quotaParts["memoryBytes"])
		if err != nil {
			return quota, err
		}
	}

	if quotaParts["storageBytes"] == "" {
		quota.StorageBytes, err = NewByte(quotaParts["storageBytes"])
		if err != nil {
			return quota, err
		}
	}

	if quotaParts["storageInodes"] == "" {
		quota.StorageInodes, err = strconv.ParseUint(quotaParts["storageInodes"], 10, 64)
		if err != nil {
			return quota, errors.New("InvalidStorageInodes")
		}
	}

	if quotaParts["storagePerformanceUnits"] == "" {
		quota.StoragePerformanceUnits, err = NewStoragePerformanceUnits(
			quotaParts["storagePerformanceUnits"],
		)
		if err != nil {
			return quota, err
		}
	}

	return NewAccountQuota(
		quota.Millicores, quota.MemoryBytes, quota.StorageBytes,
		quota.StorageInodes, quota.StoragePerformanceUnits,
	), nil
}

func NewAccountQuotaWithDefaultValues() AccountQuota {
	millicores, _ := NewMillicores(1000)
	memoryBytes, _ := NewByte(2147483648)
	storageBytes, _ := NewByte(5368709120)
	storageInodes := uint64(500000)
	storagePerformanceUnits, _ := NewStoragePerformanceUnits(5)

	return NewAccountQuota(
		millicores, memoryBytes, storageBytes, storageInodes, storagePerformanceUnits,
	)
}

func NewAccountQuotaWithBlankValues() AccountQuota {
	millicores, _ := NewMillicores(0)
	memoryBytes, _ := NewByte(0)
	storageBytes, _ := NewByte(0)
	storageInodes := uint64(0)
	storagePerformanceUnits, _ := NewStoragePerformanceUnits(0)

	return NewAccountQuota(
		millicores, memoryBytes, storageBytes, storageInodes, storagePerformanceUnits,
	)
}
