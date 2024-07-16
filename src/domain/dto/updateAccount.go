package dto

import "github.com/speedianet/control/src/domain/valueObject"

type UpdateAccount struct {
	AccountId          valueObject.AccountId     `json:"accountId"`
	Password           *valueObject.Password     `json:"password"`
	ShouldUpdateApiKey *bool                     `json:"shouldUpdateApiKey"`
	Quota              *valueObject.AccountQuota `json:"quota"`
	OperatorAccountId  valueObject.AccountId     `json:"-"`
	IpAddress          valueObject.IpAddress     `json:"-"`
}

func NewUpdateAccount(
	accountId valueObject.AccountId,
	password *valueObject.Password,
	shouldUpdateApiKey *bool,
	quota *valueObject.AccountQuota,
	operatorAccountId valueObject.AccountId,
	ipAddress valueObject.IpAddress,
) UpdateAccount {
	return UpdateAccount{
		AccountId:          accountId,
		Password:           password,
		ShouldUpdateApiKey: shouldUpdateApiKey,
		Quota:              quota,
		OperatorAccountId:  operatorAccountId,
		IpAddress:          ipAddress,
	}
}
