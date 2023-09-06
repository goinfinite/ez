package dbModel

import (
	"github.com/speedianet/sfm/src/domain/valueObject"
	"gorm.io/gorm"
)

type AccountQuota struct {
	gorm.Model
	CpuCores    float64 `gorm:"not null"`
	MemoryBytes uint64  `gorm:"not null"`
	DiskBytes   uint64  `gorm:"not null"`
	Inodes      uint64  `gorm:"not null"`
	AccountID   uint    `gorm:"not null"`
}

func (AccountQuota) TableName() string {
	return "accounts_quota"
}

func (AccountQuota) ToModel(
	vo valueObject.AccountQuota,
) (AccountQuota, error) {
	return AccountQuota{
		CpuCores:    vo.CpuCores,
		MemoryBytes: uint64(vo.MemoryBytes.Get()),
		DiskBytes:   uint64(vo.DiskBytes.Get()),
		Inodes:      vo.Inodes,
	}, nil
}

func (AccountQuota) ToValueObject(
	model AccountQuota,
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
