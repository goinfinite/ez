package entity

import "github.com/speedianet/sfm/src/domain/valueObject"

type Account struct {
	Username   valueObject.Username     `json:"username"`
	AccountId  valueObject.AccountId    `json:"id"`
	GroupId    valueObject.GroupId      `json:"groupId"`
	Quota      valueObject.AccountQuota `json:"quota"`
	QuotaUsage valueObject.AccountQuota `json:"quotaUsage"`
}

func NewAccount(
	username valueObject.Username,
	accountId valueObject.AccountId,
	groupId valueObject.GroupId,
	quota valueObject.AccountQuota,
	quotaUsage valueObject.AccountQuota,
) Account {
	return Account{
		Username:   username,
		AccountId:  accountId,
		GroupId:    groupId,
		Quota:      quota,
		QuotaUsage: quotaUsage,
	}
}
