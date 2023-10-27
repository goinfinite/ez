package valueObject

import (
	"errors"
	"strings"

	voHelper "github.com/goinfinite/fleet/src/domain/valueObject/helper"
)

type AccountQuota struct {
	CpuCores    CpuCoresCount `json:"cpuCores"`
	MemoryBytes Byte          `json:"memoryBytes"`
	DiskBytes   Byte          `json:"diskBytes"`
	Inodes      InodesCount   `json:"inodes"`
}

func NewAccountQuota(
	cpuCores CpuCoresCount,
	memoryBytes Byte,
	diskBytes Byte,
	inodes InodesCount,
) AccountQuota {
	return AccountQuota{
		CpuCores:    cpuCores,
		MemoryBytes: memoryBytes,
		DiskBytes:   diskBytes,
		Inodes:      inodes,
	}
}

func NewAccountQuotaFromString(value string) (AccountQuota, error) {
	if value == "" {
		return AccountQuota{}, errors.New("InvalidAccountQuotaValue")
	}

	if !strings.Contains(value, ":") {
		return AccountQuota{}, errors.New("InvalidAccountQuotaFormat")
	}

	specParts := strings.Split(value, ":")
	if len(specParts) != 4 {
		return AccountQuota{}, errors.New("InvalidAccountQuotaStructure")
	}

	cpuCores, err := NewCpuCoresCount(specParts[0])
	if err != nil {
		return AccountQuota{}, err
	}

	memory, err := voHelper.InterfaceToUint(specParts[1])
	if err != nil {
		return AccountQuota{}, errors.New("InvalidMemoryLimit")
	}

	disk, err := voHelper.InterfaceToUint(specParts[2])
	if err != nil {
		return AccountQuota{}, errors.New("InvalidDiskLimit")
	}

	inodes, err := NewInodesCount(specParts[3])
	if err != nil {
		return AccountQuota{}, err
	}

	return NewAccountQuota(
		cpuCores,
		Byte(int64(memory)),
		Byte(int64(disk)),
		inodes,
	), nil
}

func NewAccountQuotaWithDefaultValues() AccountQuota {
	return AccountQuota{
		CpuCores:    NewCpuCoresCountPanic(0.5),
		MemoryBytes: NewBytePanic(1073741824),
		DiskBytes:   NewBytePanic(5368709120),
		Inodes:      NewInodesCountPanic(500000),
	}
}

func NewAccountQuotaWithBlankValues() AccountQuota {
	return AccountQuota{
		CpuCores:    NewCpuCoresCountPanic(0),
		MemoryBytes: NewBytePanic(0),
		DiskBytes:   NewBytePanic(0),
		Inodes:      NewInodesCountPanic(0),
	}
}
