package dbModel

import (
	"time"

	"github.com/speedianet/control/src/domain/valueObject"
)

type MappingTarget struct {
	ID          uint `gorm:"primarykey"`
	MappingID   uint
	ContainerId string
	Port        *uint
	Protocol    *string
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
	port *uint,
	protocol *string,
) MappingTarget {
	mappingTargetStruct := MappingTarget{
		MappingID:   mappingId,
		ContainerId: containerId,
		Port:        port,
		Protocol:    protocol,
	}

	if id != 0 {
		mappingTargetStruct.ID = id
	}

	return mappingTargetStruct
}

func (model MappingTarget) ToValueObject() (valueObject.MappingTarget, error) {
	var mappingTarget valueObject.MappingTarget

	containerId, err := valueObject.NewContainerId(model.ContainerId)
	if err != nil {
		return mappingTarget, err
	}

	var portPtr *valueObject.NetworkPort
	if model.Port != nil {
		port, err := valueObject.NewNetworkPort(*model.Port)
		if err != nil {
			return mappingTarget, err
		}
		portPtr = &port
	}

	var protocolPtr *valueObject.NetworkProtocol
	if model.Protocol != nil {
		protocol, err := valueObject.NewNetworkProtocol(*model.Protocol)
		if err != nil {
			return mappingTarget, err
		}
		protocolPtr = &protocol
	}

	return valueObject.NewMappingTarget(
		containerId,
		portPtr,
		protocolPtr,
	), nil
}
