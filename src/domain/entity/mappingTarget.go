package entity

import (
	"github.com/speedianet/control/src/domain/valueObject"
)

type MappingTarget struct {
	Id            valueObject.MappingTargetId  `json:"id"`
	MappingId     valueObject.MappingId        `json:"mappingId"`
	ContainerId   valueObject.ContainerId      `json:"containerId"`
	ContainerPort *valueObject.NetworkPort     `json:"containerPort"`
	Protocol      *valueObject.NetworkProtocol `json:"protocol"`
}

func NewMappingTarget(
	id valueObject.MappingTargetId,
	mappingId valueObject.MappingId,
	containerId valueObject.ContainerId,
	containerPort *valueObject.NetworkPort,
	protocol *valueObject.NetworkProtocol,
) MappingTarget {
	return MappingTarget{
		Id:            id,
		MappingId:     mappingId,
		ContainerId:   containerId,
		ContainerPort: containerPort,
		Protocol:      protocol,
	}
}
