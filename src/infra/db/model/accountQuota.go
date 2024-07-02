package dbModel

import (
	"time"

	"github.com/speedianet/control/src/domain/valueObject"
)

type AccountQuota struct {
	ID          uint    `gorm:"primarykey"`
	CpuCores    float64 `gorm:"not null"`
	MemoryBytes uint64  `gorm:"not null"`
	DiskBytes   uint64  `gorm:"not null"`
	Inodes      uint64  `gorm:"not null"`
	AccountID   uint
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (AccountQuota) TableName() string {
	return "accounts_quota"
}

func (AccountQuota) ToModel(
	vo valueObject.AccountQuota,
	accId uint,
) (AccountQuota, error) {
	return AccountQuota{
		CpuCores:    vo.CpuCores.Read(),
		MemoryBytes: uint64(vo.MemoryBytes.Read()),
		DiskBytes:   uint64(vo.DiskBytes.Read()),
		Inodes:      vo.Inodes.Read(),
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
