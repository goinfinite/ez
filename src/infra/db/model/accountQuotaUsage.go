package dbModel

import (
	"github.com/speedianet/sfm/src/domain/valueObject"
	"gorm.io/gorm"
)

type AccountQuotaUsage struct {
	gorm.Model
	CpuCores    float64 `gorm:"not null"`
	MemoryBytes uint64  `gorm:"not null"`
	DiskBytes   uint64  `gorm:"not null"`
	Inodes      uint64  `gorm:"not null"`
	AccountID   uint    `gorm:"not null"`
}

func (AccountQuotaUsage) TableName() string {
	return "accounts_quota_usage"
}

func (AccountQuotaUsage) ToModel(
	vo valueObject.AccountQuota,
) (AccountQuotaUsage, error) {
	return AccountQuotaUsage{
		CpuCores:    vo.CpuCores,
		MemoryBytes: uint64(vo.MemoryBytes.Get()),
		DiskBytes:   uint64(vo.DiskBytes.Get()),
		Inodes:      vo.Inodes,
	}, nil
}

func (AccountQuotaUsage) ToValueObject(
	model AccountQuotaUsage,
) (valueObject.AccountQuota, error) {
	memoryBytes, err := valueObject.NewByte(model.MemoryBytes)
	if err != nil {
		return valueObject.AccountQuota{}, err
	}

	diskBytes, err := valueObject.NewByte(model.DiskBytes)
	if err != nil {
		return valueObject.AccountQuota{}, err
	}

	return valueObject.NewAccountQuota(
		model.CpuCores,
		memoryBytes,
		diskBytes,
		model.Inodes,
	), nil
}
