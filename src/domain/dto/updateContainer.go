package dto

import "github.com/goinfinite/fleet/src/domain/valueObject"

type UpdateContainer struct {
	AccountId   valueObject.AccountId           `json:"accountId"`
	ContainerId valueObject.ContainerId         `json:"id"`
	Status      *bool                           `json:"status"`
	ProfileId   *valueObject.ContainerProfileId `json:"profileId"`
}

func NewUpdateContainer(
	accountId valueObject.AccountId,
	containerId valueObject.ContainerId,
	status *bool,
	profileId *valueObject.ContainerProfileId,
) UpdateContainer {
	return UpdateContainer{
		AccountId:   accountId,
		ContainerId: containerId,
		Status:      status,
		ProfileId:   profileId,
	}
}
