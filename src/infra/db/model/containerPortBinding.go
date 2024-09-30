package dbModel

import (
	"errors"
	"slices"

	"github.com/goinfinite/ez/src/domain/valueObject"
	"gorm.io/gorm"
)

type ContainerPortBinding struct {
	ID            uint64 `gorm:"primarykey"`
	ContainerID   string `gorm:"not null"`
	ServiceName   string `gorm:"not null"`
	PublicPort    uint16 `gorm:"not null"`
	ContainerPort uint16 `gorm:"not null"`
	Protocol      string `gorm:"not null"`
	PrivatePort   uint16 `gorm:"not null"`
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
		PublicPort:    vo.PublicPort.Uint16(),
		ContainerPort: vo.ContainerPort.Uint16(),
		Protocol:      vo.Protocol.String(),
		PrivatePort:   vo.PrivatePort.Uint16(),
	}
}

func (model ContainerPortBinding) ToValueObject() (
	portBinding valueObject.PortBinding, err error,
) {
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
		serviceName, publicPort, containerPort, protocol, &privatePort,
	), nil
}

func (model ContainerPortBinding) getUnusablePorts(
	ormSvc *gorm.DB,
	portType string,
	portsToIgnore []valueObject.NetworkPort,
) (unusablePorts []uint16, err error) {
	err = ormSvc.Model(model).Select(portType).
		Order(portType + " asc").Find(&unusablePorts).Error
	if err != nil {
		return nil, err
	}

	if len(portsToIgnore) == 0 {
		return unusablePorts, nil
	}

	portsToIgnoreUint := []uint16{}
	for _, port := range portsToIgnore {
		portsToIgnoreUint = append(portsToIgnoreUint, port.Uint16())
	}
	unusablePorts = append(unusablePorts, portsToIgnoreUint...)
	unusablePorts = slices.Compact(unusablePorts)
	slices.Sort(unusablePorts)

	return unusablePorts, nil
}

func (model ContainerPortBinding) getNextAvailablePort(
	notUsablePorts []uint16,
	initialPort uint16,
	maxPort uint16,
) (nextAvailablePort valueObject.NetworkPort, err error) {
	nextPort := initialPort
	for _, port := range notUsablePorts {
		if port < initialPort {
			continue
		}

		if port == nextPort {
			nextPort++
			continue
		}

		break
	}

	if nextPort < initialPort {
		return nextAvailablePort, errors.New("PortTooLow")
	}

	if nextPort > maxPort {
		return nextAvailablePort, errors.New("NoAvailablePort")
	}

	return valueObject.NewNetworkPort(nextPort)
}

func (model ContainerPortBinding) GetNextAvailablePrivatePort(
	ormSvc *gorm.DB,
	portsToIgnore []valueObject.NetworkPort,
) (nextAvailablePort valueObject.NetworkPort, err error) {
	unusablePorts, err := model.getUnusablePorts(ormSvc, "private_port", portsToIgnore)
	if err != nil {
		return nextAvailablePort, err
	}

	return model.getNextAvailablePort(unusablePorts, 40000, 60000)
}

func (model ContainerPortBinding) GetNextAvailablePublicPort(
	ormSvc *gorm.DB,
	portBinding valueObject.PortBinding,
	portsToIgnore []valueObject.NetworkPort,
) (nextAvailablePort valueObject.NetworkPort, err error) {
	publicPortInterval, err := portBinding.GetPublicPortInterval()
	if err != nil {
		return portBinding.ContainerPort, nil
	}

	if publicPortInterval.Min == publicPortInterval.Max {
		return publicPortInterval.Min, nil
	}

	unusablePorts, err := model.getUnusablePorts(ormSvc, "public_port", portsToIgnore)
	if err != nil {
		return nextAvailablePort, err
	}

	return model.getNextAvailablePort(
		unusablePorts, publicPortInterval.Min.Uint16(), publicPortInterval.Max.Uint16(),
	)
}
