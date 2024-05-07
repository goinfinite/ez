package dto

import "github.com/speedianet/control/src/domain/valueObject"

type CreateMappingTarget struct {
	MappingId   valueObject.MappingId    `json:"mappingId"`
	ContainerId valueObject.ContainerId  `json:"containerId"`
	Path        *valueObject.MappingPath `json:"path"`
}

func NewCreateMappingTarget(
	mappingId valueObject.MappingId,
	containerId valueObject.ContainerId,
	path *valueObject.MappingPath,
) CreateMappingTarget {
	return CreateMappingTarget{
		MappingId:   mappingId,
		ContainerId: containerId,
		Path:        path,
	}
}
