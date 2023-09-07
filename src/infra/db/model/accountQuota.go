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
	AccountID   uint
}

func (AccountQuota) TableName() string {
	return "accounts_quota"
}

func (AccountQuota) ToModel(
	vo valueObject.AccountQuota,
	accId uint,
) (AccountQuota, error) {
	return AccountQuota{
		CpuCores:    vo.CpuCores.Get(),
		MemoryBytes: uint64(vo.MemoryBytes.Get()),
		DiskBytes:   uint64(vo.DiskBytes.Get()),
		Inodes:      vo.Inodes.Get(),
		AccountID:   accId,
	}, nil
}

func (model AccountQuota) ToValueObject() (valueObject.AccountQuota, error) {
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
