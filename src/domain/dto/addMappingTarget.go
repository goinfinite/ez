package dto

import "github.com/speedianet/control/src/domain/valueObject"

type AddMappingTarget struct {
	MappingId   valueObject.MappingId   `json:"mappingId"`
	ContainerId valueObject.ContainerId `json:"containerId"`
}

func NewAddMappingTarget(
	mappingId valueObject.MappingId,
	containerId valueObject.ContainerId,
) AddMappingTarget {
	return AddMappingTarget{
		MappingId:   mappingId,
		ContainerId: containerId,
	}
}
