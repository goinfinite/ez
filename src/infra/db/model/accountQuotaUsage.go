package dbModel

import (
	"time"

	"github.com/speedianet/control/src/domain/valueObject"
)

type AccountQuotaUsage struct {
	ID          uint    `gorm:"primarykey"`
	CpuCores    float64 `gorm:"not null"`
	MemoryBytes uint64  `gorm:"not null"`
	DiskBytes   uint64  `gorm:"not null"`
	Inodes      uint64  `gorm:"not null"`
	AccountID   uint
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (AccountQuotaUsage) TableName() string {
	return "accounts_quota_usage"
}

func (AccountQuotaUsage) ToModel(
	vo valueObject.AccountQuota,
	accId uint,
) (AccountQuotaUsage, error) {
	return AccountQuotaUsage{
		CpuCores:    vo.CpuCores.Read(),
		MemoryBytes: uint64(vo.MemoryBytes.Read()),
		DiskBytes:   uint64(vo.DiskBytes.Read()),
		Inodes:      vo.Inodes.Read(),
		AccountID:   accId,
	}, nil
}

func (model AccountQuotaUsage) ToValueObject() (valueObject.AccountQuota, error) {
	cpuCores, err := valueObject.NewCpuCoresCount(model.CpuCores)
	if err != nil {
		return valueObject.AccountQuota{}, err
	}

	memoryBytes, err := valueObject.NewByte(model.MemoryBytes)
	if err != nil {
		return valueObject.AccountQuota{}, err
	}

	diskBytes, err := valueObject.NewByte(model.DiskBytes)
	if err != nil {
		return valueObject.AccountQuota{}, err
	}

	inodes, err := valueObject.NewInodesCount(model.Inodes)
	if err != nil {
		return valueObject.AccountQuota{}, err
	}

	return valueObject.NewAccountQuota(
		cpuCores,
		memoryBytes,
		diskBytes,
		inodes,
	), nil
}
