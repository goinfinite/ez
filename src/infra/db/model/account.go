package dbModel

import (
	"time"

	"github.com/goinfinite/ez/src/domain/entity"
	"github.com/goinfinite/ez/src/domain/valueObject"
)

type Account struct {
	ID            uint64 `gorm:"primarykey"`
	GroupId       uint64 `gorm:"not null"`
	Username      string `gorm:"not null"`
	Quota         AccountQuota
	QuotaUsage    AccountQuotaUsage
	HomeDirectory string `gorm:"not null"`
	KeyHash       *string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (Account) TableName() string {
	return "accounts"
}

func (Account) ToModel(entity entity.Account) (model Account, err error) {
	idUint64 := entity.Id.Uint64()
	quota, err := AccountQuota{}.ToModel(entity.Quota, idUint64)
	if err != nil {
		return Account{}, err
	}

	quotaUsage, err := AccountQuotaUsage{}.ToModel(entity.QuotaUsage, idUint64)
	if err != nil {
		return Account{}, err
	}

	return Account{
		ID:            idUint64,
		GroupId:       entity.GroupId.Uint64(),
		Username:      entity.Username.String(),
		Quota:         quota,
		QuotaUsage:    quotaUsage,
		KeyHash:       nil,
		HomeDirectory: entity.HomeDirectory.String(),
	}, nil
}

func (model Account) ToEntity() (accountEntity entity.Account, err error) {
	accountId, err := valueObject.NewAccountId(model.ID)
	if err != nil {
		return accountEntity, err
	}

	groupId, err := valueObject.NewUnixGroupId(model.GroupId)
	if err != nil {
		return accountEntity, err
	}

	username, err := valueObject.NewUsername(model.Username)
	if err != nil {
		return accountEntity, err
	}

	quota, err := model.Quota.ToValueObject()
	if err != nil {
		return accountEntity, err
	}

	quotaUsage, err := model.QuotaUsage.ToValueObject()
	if err != nil {
		return accountEntity, err
	}

	homeDirectory, err := valueObject.NewUnixFilePath(model.HomeDirectory)
	if err != nil {
		return accountEntity, err
	}

	return entity.NewAccount(
		accountId, groupId, username, quota, quotaUsage, homeDirectory,
		valueObject.NewUnixTimeWithGoTime(model.CreatedAt),
		valueObject.NewUnixTimeWithGoTime(model.UpdatedAt),
	), nil
}
