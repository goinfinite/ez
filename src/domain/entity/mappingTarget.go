package entity

import (
	"github.com/goinfinite/ez/src/domain/valueObject"
)

type MappingTarget struct {
	Id                   valueObject.MappingTargetId `json:"id"`
	MappingId            valueObject.MappingId       `json:"mappingId"`
	ContainerId          valueObject.ContainerId     `json:"containerId"`
	ContainerHostname    valueObject.Fqdn            `json:"containerHostname"`
	ContainerPrivatePort valueObject.NetworkPort     `json:"containerPrivatePort"`
}

func NewMappingTarget(
	id valueObject.MappingTargetId,
	mappingId valueObject.MappingId,
	containerId valueObject.ContainerId,
	containerHostname valueObject.Fqdn,
	containerPrivatePort valueObject.NetworkPort,
) MappingTarget {
	return MappingTarget{
		Id:                   id,
		MappingId:            mappingId,
		ContainerId:          containerId,
		ContainerHostname:    containerHostname,
		ContainerPrivatePort: containerPrivatePort,
	}
}
