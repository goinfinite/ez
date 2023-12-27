package valueObject

type MappingTarget struct {
	ContainerId ContainerId     `json:"containerId"`
	Protocol    NetworkProtocol `json:"protocol"`
	Port        NetworkPort     `json:"port"`
}

func NewMappingTarget(
	containerId ContainerId,
	port NetworkPort,
	protocol NetworkProtocol,
) MappingTarget {
	return MappingTarget{
		ContainerId: containerId,
		Port:        port,
		Protocol:    protocol,
	}
}
