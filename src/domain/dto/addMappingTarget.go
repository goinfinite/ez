package dto

import "github.com/speedianet/control/src/domain/valueObject"

type AddMappingTarget struct {
	MappingId   valueObject.MappingId        `json:"mappingId"`
	ContainerId valueObject.ContainerId      `json:"containerId"`
	Port        *valueObject.NetworkPort     `json:"port"`
	Protocol    *valueObject.NetworkProtocol `json:"protocol"`
}

func NewAddMappingTarget(
	mappingId valueObject.MappingId,
	containerId valueObject.ContainerId,
	port *valueObject.NetworkPort,
	protocol *valueObject.NetworkProtocol,
) AddMappingTarget {
	return AddMappingTarget{
		MappingId:   mappingId,
		ContainerId: containerId,
		Port:        port,
		Protocol:    protocol,
	}
}
