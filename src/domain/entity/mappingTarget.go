package entity

import (
	"github.com/speedianet/control/src/domain/valueObject"
)

type MappingTarget struct {
	Id          valueObject.MappingTargetId `json:"id"`
	MappingId   valueObject.MappingId       `json:"mappingId"`
	ContainerId valueObject.ContainerId     `json:"containerId"`
}

func NewMappingTarget(
	id valueObject.MappingTargetId,
	mappingId valueObject.MappingId,
	containerId valueObject.ContainerId,
) MappingTarget {
	return MappingTarget{
		Id:          id,
		MappingId:   mappingId,
		ContainerId: containerId,
	}
}
