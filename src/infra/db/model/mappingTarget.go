package dbModel

import (
	"time"

	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/valueObject"
)

type MappingTarget struct {
	ID            uint `gorm:"primarykey"`
	MappingID     uint
	ContainerId   string
	ContainerPort *uint
	Protocol      *string
	CreatedAt     time.Time
	UpdatedAt     time.Time
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
		MappingID:     mappingId,
		ContainerId:   containerId,
		ContainerPort: containerPort,
		Protocol:      protocol,
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

	var containerPortPtr *valueObject.NetworkPort
	if model.ContainerPort != nil {
		containerPort, err := valueObject.NewNetworkPort(*model.ContainerPort)
		if err != nil {
			return mappingTarget, err
		}
		containerPortPtr = &containerPort
	}

	var protocolPtr *valueObject.NetworkProtocol
	if model.Protocol != nil {
		protocol, err := valueObject.NewNetworkProtocol(*model.Protocol)
		if err != nil {
			return mappingTarget, err
		}
		protocolPtr = &protocol
	}

	return entity.NewMappingTarget(
		id,
		mappingId,
		containerId,
		containerPortPtr,
		protocolPtr,
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

	if addDto.ContainerPort != nil {
		containerPortUint := uint(addDto.ContainerPort.Get())
		model.ContainerPort = &containerPortUint
	}

	if addDto.Protocol != nil {
		protocolStr := addDto.Protocol.String()
		model.Protocol = &protocolStr
	}

	return model
}
