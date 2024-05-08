package dbModel

import (
	"time"

	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/valueObject"
)

type MappingTarget struct {
	ID                   uint `gorm:"primarykey"`
	MappingID            uint
	ContainerId          string
	ContainerHostname    string
	ContainerPrivatePort uint
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

func (MappingTarget) TableName() string {
	return "mapping_targets"
}

func NewMappingTarget(
	id uint,
	mappingId uint,
	containerId string,
	containerHostname string,
	containerPrivatePort uint,
) MappingTarget {
	targetModel := MappingTarget{
		MappingID:            mappingId,
		ContainerId:          containerId,
		ContainerHostname:    containerHostname,
		ContainerPrivatePort: containerPrivatePort,
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

	containerPrivatePort, err := valueObject.NewNetworkPort(model.ContainerPrivatePort)
	if err != nil {
		return mappingTarget, err
	}

	return entity.NewMappingTarget(
		id,
		mappingId,
		containerId,
		containerHostname,
		containerPrivatePort,
	), nil
}
