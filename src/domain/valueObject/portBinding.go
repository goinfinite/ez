package valueObject

import "strconv"

type PortBinding struct {
	Protocol      NetworkProtocol `json:"protocol"`
	HostPort      uint64          `json:"hostPort"`
	ContainerPort uint64          `json:"containerPort"`
}

func NewPortBinding(
	protocol NetworkProtocol,
	hostPort uint64,
	containerPort uint64,
) PortBinding {
	return PortBinding{
		Protocol:      protocol,
		HostPort:      hostPort,
		ContainerPort: containerPort,
	}
}

func (portBinding PortBinding) GetProtocol() NetworkProtocol {
	return portBinding.Protocol
}

func (portBinding PortBinding) GetHostPort() uint64 {
	return portBinding.HostPort
}

func (portBinding PortBinding) GetContainerPort() uint64 {
	return portBinding.ContainerPort
}

func (portBinding PortBinding) GetHostPortAsString() string {
	return strconv.FormatUint(portBinding.HostPort, 10)
}

func (portBinding PortBinding) GetContainerPortAsString() string {
	return strconv.FormatUint(portBinding.ContainerPort, 10)
}
