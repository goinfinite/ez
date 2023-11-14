package dto

import "github.com/speedianet/control/src/domain/valueObject"

type AddAccount struct {
	Username valueObject.Username      `json:"username"`
	Password valueObject.Password      `json:"password"`
	Quota    *valueObject.AccountQuota `json:"quota"`
}

func NewAddAccount(
	username valueObject.Username,
	password valueObject.Password,
	quota *valueObject.AccountQuota,
) AddAccount {
	return AddAccount{
		Username: username,
		Password: password,
		Quota:    quota,
	}
}
