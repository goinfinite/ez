package dto

import (
	"strings"

	"github.com/speedianet/control/src/domain/valueObject"
)

type AddMappingTargetWithoutMappingId struct {
	ContainerId valueObject.ContainerId      `json:"containerId"`
	Port        *valueObject.NetworkPort     `json:"port"`
	Protocol    *valueObject.NetworkProtocol `json:"protocol"`
}

func NewAddMappingTargetWithoutMappingId(
	containerId valueObject.ContainerId,
	port *valueObject.NetworkPort,
	protocol *valueObject.NetworkProtocol,
) AddMappingTargetWithoutMappingId {
	return AddMappingTargetWithoutMappingId{
		ContainerId: containerId,
		Port:        port,
		Protocol:    protocol,
	}
}

func NewAddMappingTargetWithoutMappingIdFromString(
	value string,
) (AddMappingTargetWithoutMappingId, error) {
	var target AddMappingTargetWithoutMappingId

	targetParts := strings.Split(value, ":")
	containerId, err := valueObject.NewContainerId(targetParts[0])
	if err != nil {
		return target, err
	}

	if len(targetParts) == 1 {
		return NewAddMappingTargetWithoutMappingId(
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

	port, err := valueObject.NewNetworkPort(hostPortStr)
	if err != nil {
		return target, err
	}

	protocol, err := valueObject.NewNetworkProtocol(hostProtocolStr)
	if err != nil {
		return target, err
	}

	return NewAddMappingTargetWithoutMappingId(
		containerId,
		&port,
		&protocol,
	), nil
}
