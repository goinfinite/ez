package valueObject

import "strings"

type MappingTarget struct {
	ContainerId ContainerId      `json:"containerId"`
	Protocol    *NetworkProtocol `json:"protocol"`
	Port        *NetworkPort     `json:"port"`
}

func NewMappingTarget(
	containerId ContainerId,
	port *NetworkPort,
	protocol *NetworkProtocol,
) MappingTarget {
	return MappingTarget{
		ContainerId: containerId,
		Port:        port,
		Protocol:    protocol,
	}
}

func NewMappingTargetFromString(value string) (MappingTarget, error) {
	var target MappingTarget

	targetParts := strings.Split(value, ":")
	containerId, err := NewContainerId(targetParts[0])
	if err != nil {
		return target, err
	}

	if len(targetParts) == 1 {
		return NewMappingTarget(
			containerId,
			nil,
			nil,
		), nil
	}

	portProtocolParts := strings.Split(targetParts[1], "/")
	hostPortStr := portProtocolParts[0]
	hostProtocolStr := "tcp"
	if len(portProtocolParts) == 2 {
		hostProtocolStr = portProtocolParts[1]
	}

	port, err := NewNetworkPort(hostPortStr)
	if err != nil {
		return target, err
	}

	protocol, err := NewNetworkProtocol(hostProtocolStr)
	if err != nil {
		return target, err
	}

	return NewMappingTarget(
		containerId,
		&port,
		&protocol,
	), nil
}
