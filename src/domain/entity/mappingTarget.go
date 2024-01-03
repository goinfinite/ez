package entity

import (
	"github.com/speedianet/control/src/domain/valueObject"
)

type MappingTarget struct {
	Id          valueObject.MappingTargetId  `json:"id"`
	MappingId   valueObject.MappingId        `json:"mappingId"`
	ContainerId valueObject.ContainerId      `json:"containerId"`
	Protocol    *valueObject.NetworkProtocol `json:"protocol"`
	Port        *valueObject.NetworkPort     `json:"port"`
}

func NewMappingTarget(
	id valueObject.MappingTargetId,
	mappingId valueObject.MappingId,
	containerId valueObject.ContainerId,
	port *valueObject.NetworkPort,
	protocol *valueObject.NetworkProtocol,
) MappingTarget {
	return MappingTarget{
		Id:          id,
		MappingId:   mappingId,
		ContainerId: containerId,
		Port:        port,
		Protocol:    protocol,
	}
}
