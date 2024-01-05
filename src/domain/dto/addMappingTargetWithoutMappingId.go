package dto

import (
	"strings"

	"github.com/speedianet/control/src/domain/valueObject"
)

type AddMappingTargetWithoutMappingId struct {
	ContainerId   valueObject.ContainerId      `json:"containerId"`
	ContainerPort *valueObject.NetworkPort     `json:"containerPort"`
	Protocol      *valueObject.NetworkProtocol `json:"protocol"`
}

func NewAddMappingTargetWithoutMappingId(
	containerId valueObject.ContainerId,
	containerPort *valueObject.NetworkPort,
	protocol *valueObject.NetworkProtocol,
) AddMappingTargetWithoutMappingId {
	return AddMappingTargetWithoutMappingId{
		ContainerId:   containerId,
		ContainerPort: containerPort,
		Protocol:      protocol,
	}
}

// format: containerId:containerPort/protocol
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
	containerPortStr := portProtocolParts[0]
	protocolStr := "tcp"
	if len(portProtocolParts) == 2 {
		protocolStr = portProtocolParts[1]
	}

	port, err := valueObject.NewNetworkPort(containerPortStr)
	if err != nil {
		return target, err
	}

	protocol, err := valueObject.NewNetworkProtocol(protocolStr)
	if err != nil {
		return target, err
	}

	return NewAddMappingTargetWithoutMappingId(
		containerId,
		&port,
		&protocol,
	), nil
}
