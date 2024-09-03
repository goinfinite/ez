package dto

import "github.com/speedianet/control/src/domain/valueObject"

type CreateContainerSessionToken struct {
	AccountId         valueObject.AccountId   `json:"accountId"`
	ContainerId       valueObject.ContainerId `json:"containerId"`
	SessionIpAddress  valueObject.IpAddress   `json:"sessionIpAddress"`
	OperatorAccountId valueObject.AccountId   `json:"-"`
	OperatorIpAddress valueObject.IpAddress   `json:"-"`
}

func NewCreateContainerSessionToken(
	accountId valueObject.AccountId,
	containerId valueObject.ContainerId,
	sessionIpAddress valueObject.IpAddress,
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
) CreateContainerSessionToken {
	return CreateContainerSessionToken{
		AccountId:         accountId,
		ContainerId:       containerId,
		SessionIpAddress:  sessionIpAddress,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
	}
}
