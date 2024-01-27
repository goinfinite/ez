package dbModel

import (
	"errors"

	"github.com/speedianet/control/src/domain/valueObject"
	"gorm.io/gorm"
)

type ContainerPortBinding struct {
	ID            uint   `gorm:"primarykey"`
	ContainerID   string `gorm:"not null"`
	ServiceName   string `gorm:"not null"`
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
		ServiceName:   vo.ServiceName.String(),
		PublicPort:    uint(vo.PublicPort.Get()),
		ContainerPort: uint(vo.ContainerPort.Get()),
		Protocol:      vo.Protocol.String(),
		PrivatePort:   uint(vo.PrivatePort.Get()),
	}
}

func (model ContainerPortBinding) ToValueObject() (valueObject.PortBinding, error) {
	var portBinding valueObject.PortBinding

	serviceName, err := valueObject.NewServiceName(model.ServiceName)
	if err != nil {
		return portBinding, err
	}

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
		serviceName,
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

	if len(portsToIgnore) > 0 {
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

func (model ContainerPortBinding) GetNextAvailablePublicPort(
	ormSvc *gorm.DB,
	portBinding valueObject.PortBinding,
	portsToIgnore []valueObject.NetworkPort,
) (valueObject.NetworkPort, error) {
	usedPublicPorts := []uint{}

	err := ormSvc.Model(model).
		Select("public_port").
		Order("public_port asc").
		Find(&usedPublicPorts).Error
	if err != nil {
		return 0, err
	}

	if len(portsToIgnore) > 0 {
		portsToIgnoreUint := []uint{}
		for _, port := range portsToIgnore {
			portsToIgnoreUint = append(portsToIgnoreUint, uint(port.Get()))
		}
		usedPublicPorts = append(usedPublicPorts, portsToIgnoreUint...)
	}

	publicPortInterval, err := portBinding.GetPublicPortInterval()
	if err != nil {
		return portBinding.ContainerPort, nil
	}
	if publicPortInterval.Max == nil {
		return publicPortInterval.Min, nil
	}

	initialPort := uint(publicPortInterval.Min.Get())

	nextPort := uint(initialPort)
	for _, port := range usedPublicPorts {
		if port == nextPort {
			nextPort++
			continue
		}
		break
	}

	if nextPort < initialPort {
		return 0, errors.New("PublicPortTooLow")
	}

	maxPort := uint(publicPortInterval.Max.Get())
	if nextPort > maxPort {
		return 0, errors.New("NoAvailablePublicPort")
	}

	return valueObject.NewNetworkPort(nextPort)
}
