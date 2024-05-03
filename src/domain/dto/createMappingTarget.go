package dto

import "github.com/speedianet/control/src/domain/valueObject"

type CreateMappingTarget struct {
	MappingId   valueObject.MappingId   `json:"mappingId"`
	ContainerId valueObject.ContainerId `json:"containerId"`
}

func NewCreateMappingTarget(
	mappingId valueObject.MappingId,
	containerId valueObject.ContainerId,
) CreateMappingTarget {
	return CreateMappingTarget{
		MappingId:   mappingId,
		ContainerId: containerId,
	}
}
