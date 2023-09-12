package valueObject

type PortBinding struct {
	Protocol      NetworkProtocol `json:"protocol"`
	HostPort      NetworkPort     `json:"hostPort"`
	ContainerPort NetworkPort     `json:"containerPort"`
}

func NewPortBinding(
	protocol NetworkProtocol,
	hostPort NetworkPort,
	containerPort NetworkPort,
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

func (portBinding PortBinding) GetHostPort() NetworkPort {
	return portBinding.HostPort
}

func (portBinding PortBinding) GetContainerPort() NetworkPort {
	return portBinding.ContainerPort
}
