package dto

import "github.com/speedianet/sfm/src/domain/valueObject"

type UpdateContainer struct {
	AccountId   valueObject.AccountId       `json:"accountId"`
	ContainerId valueObject.ContainerId     `json:"id"`
	Status      bool                        `json:"status"`
	BaseSpecs   *valueObject.ContainerSpecs `json:"baseSpecs"`
	MaxSpecs    *valueObject.ContainerSpecs `json:"maxSpecs"`
}

func NewUpdateContainer(
	accountId valueObject.AccountId,
	containerId valueObject.ContainerId,
	status bool,
	baseSpecs *valueObject.ContainerSpecs,
	maxSpecs *valueObject.ContainerSpecs,
) UpdateContainer {
	return UpdateContainer{
		AccountId:   accountId,
		ContainerId: containerId,
		Status:      status,
		BaseSpecs:   baseSpecs,
		MaxSpecs:    maxSpecs,
	}
}
