package dto

import (
	"github.com/speedianet/control/src/domain/valueObject"
)

type DeleteContainer struct {
	AccountId         valueObject.AccountId   `json:"accountId"`
	ContainerId       valueObject.ContainerId `json:"containerId"`
	OperatorAccountId valueObject.AccountId   `json:"-"`
	OperatorIpAddress valueObject.IpAddress   `json:"-"`
}

func NewDeleteContainer(
	accountId valueObject.AccountId,
	containerId valueObject.ContainerId,
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
) DeleteContainer {
	return DeleteContainer{
		AccountId:         accountId,
		ContainerId:       containerId,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
	}
}
