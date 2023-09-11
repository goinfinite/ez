package dto

import "github.com/speedianet/sfm/src/domain/valueObject"

type UpdateAccount struct {
	AccountId          valueObject.AccountId     `json:"accountId"`
	Password           *valueObject.Password     `json:"password"`
	ShouldUpdateApiKey *bool                     `json:"shouldUpdateApiKey"`
	Quota              *valueObject.AccountQuota `json:"quota"`
}

func NewUpdateAccount(
	accountId valueObject.AccountId,
	password *valueObject.Password,
	shouldUpdateApiKey *bool,
	quota *valueObject.AccountQuota,
) UpdateAccount {
	return UpdateAccount{
		AccountId:          accountId,
		Password:           password,
		ShouldUpdateApiKey: shouldUpdateApiKey,
		Quota:              quota,
	}
}
