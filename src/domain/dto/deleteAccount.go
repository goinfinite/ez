package dto

import (
	"github.com/speedianet/control/src/domain/valueObject"
)

type DeleteAccount struct {
	AccountId valueObject.AccountId `json:"accountId"`
	IpAddress valueObject.IpAddress `json:"-"`
}

func NewDeleteAccount(
	accountId valueObject.AccountId,
	ipAddress valueObject.IpAddress,
) DeleteAccount {
	return DeleteAccount{
		AccountId: accountId,
		IpAddress: ipAddress,
	}
}
