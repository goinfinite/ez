package dto

import "github.com/speedianet/sfm/src/domain/valueObject"

type UpdateContainer struct {
	AccountId         valueObject.AccountId          `json:"accountId"`
	ContainerId       valueObject.ContainerId        `json:"id"`
	Status            *bool                          `json:"status"`
	ResourceProfileId *valueObject.ResourceProfileId `json:"resourceProfileId"`
}

func NewUpdateContainer(
	accountId valueObject.AccountId,
	containerId valueObject.ContainerId,
	status *bool,
	resourceProfileId *valueObject.ResourceProfileId,
) UpdateContainer {
	return UpdateContainer{
		AccountId:         accountId,
		ContainerId:       containerId,
		Status:            status,
		ResourceProfileId: resourceProfileId,
	}
}
