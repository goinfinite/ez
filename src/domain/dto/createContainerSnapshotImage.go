package dto

import (
	"github.com/speedianet/control/src/domain/valueObject"
)

type CreateContainerSnapshotImage struct {
	AccountId         valueObject.AccountId   `json:"accountId"`
	ContainerId       valueObject.ContainerId `json:"containerId"`
	OperatorAccountId valueObject.AccountId   `json:"-"`
	OperatorIpAddress valueObject.IpAddress   `json:"-"`
}

func NewCreateContainerSnapshotImage(
	accountId valueObject.AccountId,
	containerId valueObject.ContainerId,
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
) CreateContainerSnapshotImage {
	return CreateContainerSnapshotImage{
		AccountId:         accountId,
		ContainerId:       containerId,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
	}
}
