package dto

import "github.com/speedianet/control/src/domain/valueObject"

type UpdateContainer struct {
	AccountId         valueObject.AccountId           `json:"accountId"`
	ContainerId       valueObject.ContainerId         `json:"id"`
	Status            *bool                           `json:"status,omitempty"`
	ProfileId         *valueObject.ContainerProfileId `json:"profileId,omitempty"`
	OperatorAccountId valueObject.AccountId           `json:"-"`
	OperatorIpAddress valueObject.IpAddress           `json:"-"`
}

func NewUpdateContainer(
	accountId valueObject.AccountId,
	containerId valueObject.ContainerId,
	status *bool,
	profileId *valueObject.ContainerProfileId,
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
) UpdateContainer {
	return UpdateContainer{
		AccountId:         accountId,
		ContainerId:       containerId,
		Status:            status,
		ProfileId:         profileId,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
	}
}
