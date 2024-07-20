package dbModel

import (
	"time"

	"github.com/speedianet/control/src/domain/valueObject"
)

type AccountQuota struct {
	ID                      uint64 `gorm:"primarykey"`
	Millicores              uint   `gorm:"not null"`
	MemoryBytes             uint64 `gorm:"not null"`
	StorageBytes            uint64 `gorm:"not null"`
	StorageInodes           uint64 `gorm:"not null"`
	StoragePerformanceUnits uint   `gorm:"not null"`
	AccountID               uint64
	CreatedAt               time.Time
	UpdatedAt               time.Time
}

func (AccountQuota) TableName() string {
	return "accounts_quota"
}

func (AccountQuota) ToModel(
	vo valueObject.AccountQuota, accountId uint64,
) (AccountQuota, error) {
	return AccountQuota{
		Millicores:              vo.Millicores.Uint(),
		MemoryBytes:             uint64(vo.MemoryBytes),
		StorageBytes:            uint64(vo.StorageBytes),
		StorageInodes:           vo.StorageInodes,
		StoragePerformanceUnits: vo.StoragePerformanceUnits.Uint(),
		AccountID:               accountId,
		CreatedAt:               time.Now(),
		UpdatedAt:               time.Now(),
	}, nil
}

func (model AccountQuota) ToValueObject() (vo valueObject.AccountQuota, err error) {
	millicores, err := valueObject.NewMillicores(model.Millicores)
	if err != nil {
		return vo, err
	}

	memoryBytes, err := valueObject.NewByte(model.MemoryBytes)
	if err != nil {
		return vo, err
	}

	storageBytes, err := valueObject.NewByte(model.StorageBytes)
	if err != nil {
		return vo, err
	}

	storagePerformanceUnits, err := valueObject.NewStoragePerformanceUnits(
		model.StoragePerformanceUnits,
	)
	if err != nil {
		return vo, err
	}

	return valueObject.NewAccountQuota(
		millicores, memoryBytes, storageBytes, model.StorageInodes, storagePerformanceUnits,
	), nil
}
