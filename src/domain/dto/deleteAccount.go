package dto

import (
	"github.com/speedianet/control/src/domain/valueObject"
)

type DeleteAccount struct {
	AccountId         valueObject.AccountId `json:"accountId"`
	OperatorAccountId valueObject.AccountId `json:"-"`
	IpAddress         valueObject.IpAddress `json:"-"`
}

func NewDeleteAccount(
	accountId valueObject.AccountId,
	operatorAccountId valueObject.AccountId,
	ipAddress valueObject.IpAddress,
) DeleteAccount {
	return DeleteAccount{
		AccountId:         accountId,
		OperatorAccountId: operatorAccountId,
		IpAddress:         ipAddress,
	}
}
