package entity

import (
	"github.com/speedianet/control/src/domain/valueObject"
)

type MappingTarget struct {
	Id                valueObject.MappingTargetId `json:"id"`
	MappingId         valueObject.MappingId       `json:"mappingId"`
	ContainerId       valueObject.ContainerId     `json:"containerId"`
	ContainerHostname valueObject.Fqdn            `json:"containerHostname"`
	Path              *valueObject.MappingPath    `json:"path"`
}

func NewMappingTarget(
	id valueObject.MappingTargetId,
	mappingId valueObject.MappingId,
	containerId valueObject.ContainerId,
	containerHostname valueObject.Fqdn,
	path *valueObject.MappingPath,
) MappingTarget {
	return MappingTarget{
		Id:                id,
		MappingId:         mappingId,
		ContainerId:       containerId,
		ContainerHostname: containerHostname,
		Path:              path,
	}
}
