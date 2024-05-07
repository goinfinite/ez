package dbModel

import (
	"time"

	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/valueObject"
)

type MappingTarget struct {
	ID                uint `gorm:"primarykey"`
	MappingID         uint
	ContainerId       string
	ContainerHostname string
	Path              *string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func (MappingTarget) TableName() string {
	return "mapping_targets"
}

func NewMappingTarget(
	id uint,
	mappingId uint,
	containerId string,
	containerHostname string,
	path *string,
) MappingTarget {
	targetModel := MappingTarget{
		MappingID:         mappingId,
		ContainerId:       containerId,
		ContainerHostname: containerHostname,
		Path:              path,
	}

	if id != 0 {
		targetModel.ID = id
	}

	return targetModel
}

func (model MappingTarget) ToEntity() (entity.MappingTarget, error) {
	var mappingTarget entity.MappingTarget

	id, err := valueObject.NewMappingTargetId(model.ID)
	if err != nil {
		return mappingTarget, err
	}

	mappingId, err := valueObject.NewMappingId(model.MappingID)
	if err != nil {
		return mappingTarget, err
	}

	containerId, err := valueObject.NewContainerId(model.ContainerId)
	if err != nil {
		return mappingTarget, err
	}

	containerHostname, err := valueObject.NewFqdn(model.ContainerHostname)
	if err != nil {
		return mappingTarget, err
	}

	var pathPtr *valueObject.MappingPath
	if model.Path != nil {
		path, err := valueObject.NewMappingPath(*model.Path)
		if err != nil {
			return mappingTarget, err
		}
		pathPtr = &path
	}

	return entity.NewMappingTarget(
		id,
		mappingId,
		containerId,
		containerHostname,
		pathPtr,
	), nil
}
