package dto

import "github.com/speedianet/control/src/domain/valueObject"

type AddMappingTarget struct {
	MappingId     valueObject.MappingId        `json:"mappingId"`
	ContainerId   valueObject.ContainerId      `json:"containerId"`
	ContainerPort *valueObject.NetworkPort     `json:"containerPort"`
	Protocol      *valueObject.NetworkProtocol `json:"protocol"`
}

func NewAddMappingTarget(
	mappingId valueObject.MappingId,
	containerId valueObject.ContainerId,
	containerPort *valueObject.NetworkPort,
	protocol *valueObject.NetworkProtocol,
) AddMappingTarget {
	return AddMappingTarget{
		MappingId:     mappingId,
		ContainerId:   containerId,
		ContainerPort: containerPort,
		Protocol:      protocol,
	}
}
