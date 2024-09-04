package dto

import "github.com/speedianet/control/src/domain/valueObject"

type CreateAccount struct {
	Username          valueObject.Username      `json:"username"`
	Password          valueObject.Password      `json:"password"`
	Quota             *valueObject.AccountQuota `json:"quota,omitempty"`
	OperatorAccountId valueObject.AccountId     `json:"-"`
	OperatorIpAddress valueObject.IpAddress     `json:"-"`
}

func NewCreateAccount(
	username valueObject.Username,
	password valueObject.Password,
	quota *valueObject.AccountQuota,
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
) CreateAccount {
	return CreateAccount{
		Username:          username,
		Password:          password,
		Quota:             quota,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
	}
}
