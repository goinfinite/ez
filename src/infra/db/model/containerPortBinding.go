package dbModel

import (
	"errors"

	"github.com/speedianet/control/src/domain/valueObject"
	"gorm.io/gorm"
)

type ContainerPortBinding struct {
	ID            uint   `gorm:"primarykey"`
	ContainerID   string `gorm:"not null"`
	PublicPort    uint   `gorm:"not null"`
	ContainerPort uint   `gorm:"not null"`
	Protocol      string `gorm:"not null"`
	PrivatePort   uint   `gorm:"not null"`
}

func (ContainerPortBinding) TableName() string {
	return "container_port_bindings"
}

func (ContainerPortBinding) ToModel(
	containerId valueObject.ContainerId,
	vo valueObject.PortBinding,
) ContainerPortBinding {
	return ContainerPortBinding{
		ContainerID:   containerId.String(),
		PublicPort:    uint(vo.PublicPort.Get()),
		ContainerPort: uint(vo.ContainerPort.Get()),
		Protocol:      vo.Protocol.String(),
		PrivatePort:   uint(vo.PrivatePort.Get()),
	}
}

func (model ContainerPortBinding) ToValueObject() (valueObject.PortBinding, error) {
	var portBinding valueObject.PortBinding

	publicPort, err := valueObject.NewNetworkPort(model.PublicPort)
	if err != nil {
		return portBinding, err
	}

	containerPort, err := valueObject.NewNetworkPort(model.ContainerPort)
	if err != nil {
		return portBinding, err
	}

	protocol, err := valueObject.NewNetworkProtocol(model.Protocol)
	if err != nil {
		return portBinding, err
	}

	privatePort, err := valueObject.NewNetworkPort(model.PrivatePort)
	if err != nil {
		return portBinding, err
	}

	return valueObject.NewPortBinding(
		publicPort,
		containerPort,
		protocol,
		&privatePort,
	), nil
}

func (model ContainerPortBinding) GetNextAvailablePrivatePort(
	ormSvc *gorm.DB,
	portsToIgnore []valueObject.NetworkPort,
) (valueObject.NetworkPort, error) {
	usedPrivatePorts := []uint{}

	err := ormSvc.Model(model).
		Select("private_port").
		Order("private_port asc").
		Find(&usedPrivatePorts).Error
	if err != nil {
		return 0, err
	}

	if len(usedPrivatePorts) > 0 {
		portsToIgnoreUint := []uint{}
		for _, port := range portsToIgnore {
			portsToIgnoreUint = append(portsToIgnoreUint, uint(port.Get()))
		}
		usedPrivatePorts = append(usedPrivatePorts, portsToIgnoreUint...)
	}

	initialPort := uint(40000)

	nextPort := initialPort
	for _, port := range usedPrivatePorts {
		if port == nextPort {
			nextPort++
			continue
		}
		break
	}

	if nextPort < initialPort {
		return 0, errors.New("PrivatePortTooLow")
	}

	if nextPort > 60000 {
		return 0, errors.New("NoAvailablePrivatePort")
	}

	return valueObject.NewNetworkPort(nextPort)
}
