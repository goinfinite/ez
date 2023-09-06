package dbModel

import (
	"github.com/speedianet/sfm/src/domain/entity"
	"github.com/speedianet/sfm/src/domain/valueObject"
	"gorm.io/gorm"
)

type Account struct {
	gorm.Model
	GroupID           uint   `gorm:"not null"`
	Username          string `gorm:"not null"`
	KeyHash           *string
	AccountQuota      AccountQuota
	AccountQuotaUsage AccountQuotaUsage
}

func (Account) TableName() string {
	return "accounts"
}

func (Account) ToEntity(model Account) (entity.Account, error) {
	accId, err := valueObject.NewAccountId(model.ID)
	if err != nil {
		return entity.Account{}, err
	}

	groupId, err := valueObject.NewGroupId(model.GroupID)
	if err != nil {
		return entity.Account{}, err
	}

	username, err := valueObject.NewUsername(model.Username)
	if err != nil {
		return entity.Account{}, err
	}

	quota, err := model.AccountQuota.ToValueObject(model.AccountQuota)
	if err != nil {
		return entity.Account{}, err
	}

	quotaUsage, err := model.AccountQuotaUsage.ToValueObject(model.AccountQuotaUsage)
	if err != nil {
		return entity.Account{}, err
	}

	return entity.NewAccount(
		accId,
		groupId,
		username,
		quota,
		quotaUsage,
		valueObject.UnixTime(model.CreatedAt.Unix()),
		valueObject.UnixTime(model.UpdatedAt.Unix()),
	), nil
}

func (Account) ToModel(entity entity.Account) (Account, error) {
	quota, err := AccountQuota{}.ToModel(entity.Quota)
	if err != nil {
		return Account{}, err
	}

	quotaUsage, err := AccountQuotaUsage{}.ToModel(entity.QuotaUsage)
	if err != nil {
		return Account{}, err
	}

	return Account{
		Model:             gorm.Model{ID: uint(entity.Id.Get())},
		GroupID:           uint(entity.GroupId.Get()),
		Username:          entity.Username.String(),
		KeyHash:           nil,
		AccountQuota:      quota,
		AccountQuotaUsage: quotaUsage,
	}, nil
}
