package repository

import (
	"github.com/speedianet/sfm/src/domain/dto"
	"github.com/speedianet/sfm/src/domain/valueObject"
)

type AccCmdRepo interface {
	Add(addAccount dto.AddAccount) error
	Delete(accountId valueObject.AccountId) error
	UpdatePassword(accountId valueObject.AccountId, password valueObject.Password) error
	UpdateApiKey(accountId valueObject.AccountId) (valueObject.AccessTokenStr, error)
	UpdateQuota(accountId valueObject.AccountId, quota valueObject.AccountQuota) error
}
