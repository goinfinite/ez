package valueObject

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
