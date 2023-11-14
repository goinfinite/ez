package entity

import "github.com/speedianet/control/src/domain/valueObject"

type Account struct {
	Id         valueObject.AccountId    `json:"id"`
	GroupId    valueObject.GroupId      `json:"groupId"`
	Username   valueObject.Username     `json:"username"`
	Quota      valueObject.AccountQuota `json:"quota"`
	QuotaUsage valueObject.AccountQuota `json:"quotaUsage"`
	CreatedAt  valueObject.UnixTime     `json:"createdAt"`
	UpdatedAt  valueObject.UnixTime     `json:"updatedAt"`
}

func NewAccount(
	id valueObject.AccountId,
	groupId valueObject.GroupId,
	username valueObject.Username,
	quota valueObject.AccountQuota,
	quotaUsage valueObject.AccountQuota,
	createdAt valueObject.UnixTime,
	updatedAt valueObject.UnixTime,
) Account {
	return Account{
		Id:         id,
		GroupId:    groupId,
		Username:   username,
		Quota:      quota,
		QuotaUsage: quotaUsage,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	}
}
