package repository

import (
	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/valueObject"
)

type AccountCmdRepo interface {
	Add(addAccount dto.AddAccount) error
	Delete(accId valueObject.AccountId) error
	UpdatePassword(accId valueObject.AccountId, password valueObject.Password) error
	UpdateApiKey(accId valueObject.AccountId) (valueObject.AccessTokenValue, error)
	UpdateQuota(accId valueObject.AccountId, quota valueObject.AccountQuota) error
	UpdateQuotaUsage(accId valueObject.AccountId) error
}
