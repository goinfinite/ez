package dbModel

import (
	"time"

	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/valueObject"
)

type Account struct {
	ID         uint64 `gorm:"primarykey"`
	GroupId    uint64 `gorm:"not null"`
	Username   string `gorm:"not null"`
	KeyHash    *string
	Quota      AccountQuota
	QuotaUsage AccountQuotaUsage
	CreatedAt  time.Time
	UpdatedAt  time.Time
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
		ID:         idUint64,
		GroupId:    entity.GroupId.Uint64(),
		Username:   entity.Username.String(),
		KeyHash:    nil,
		Quota:      quota,
		QuotaUsage: quotaUsage,
	}, nil
}

func (model Account) ToEntity() (entity.Account, error) {
	accountId, err := valueObject.NewAccountId(model.ID)
	if err != nil {
		return entity.Account{}, err
	}

	groupId, err := valueObject.NewUnixGroupId(model.GroupId)
	if err != nil {
		return entity.Account{}, err
	}

	username, err := valueObject.NewUsername(model.Username)
	if err != nil {
		return entity.Account{}, err
	}

	quota, err := model.Quota.ToValueObject()
	if err != nil {
		return entity.Account{}, err
	}

	quotaUsage, err := model.QuotaUsage.ToValueObject()
	if err != nil {
		return entity.Account{}, err
	}

	return entity.NewAccount(
		accountId,
		groupId,
		username,
		quota,
		quotaUsage,
		valueObject.NewUnixTimeWithGoTime(model.CreatedAt),
		valueObject.NewUnixTimeWithGoTime(model.UpdatedAt),
	), nil
}
