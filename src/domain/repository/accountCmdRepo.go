package repository

import (
	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/valueObject"
)

type AccountCmdRepo interface {
	Create(createAccount dto.CreateAccount) (valueObject.AccountId, error)
	Delete(accountId valueObject.AccountId) error
	UpdatePassword(accountId valueObject.AccountId, password valueObject.Password) error
	UpdateApiKey(accountId valueObject.AccountId) (valueObject.AccessTokenValue, error)
	UpdateQuota(accountId valueObject.AccountId, quota valueObject.AccountQuota) error
	UpdateQuotaUsage(accountId valueObject.AccountId) error
}
