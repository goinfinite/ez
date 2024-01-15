package dbModel

import (
	"time"

	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/valueObject"
)

type MappingTarget struct {
	ID          uint `gorm:"primarykey"`
	MappingID   uint
	ContainerId string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (MappingTarget) TableName() string {
	return "mapping_targets"
}

func NewMappingTarget(
	id uint,
	mappingId uint,
	containerId string,
	containerPort *uint,
	protocol *string,
) MappingTarget {
	targetModel := MappingTarget{
		MappingID:   mappingId,
		ContainerId: containerId,
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

	return entity.NewMappingTarget(
		id,
		mappingId,
		containerId,
	), nil
}

func (MappingTarget) AddDtoToModel(addDto dto.AddMappingTarget) MappingTarget {
	model := NewMappingTarget(
		0,
		uint(addDto.MappingId),
		addDto.ContainerId.String(),
		nil,
		nil,
	)

	return model
}
