package valueObject

import (
	"errors"
	"strings"
)

type PortBinding struct {
	PublicPort    NetworkPort     `json:"publicPort"`
	ContainerPort NetworkPort     `json:"containerPort"`
	Protocol      NetworkProtocol `json:"protocol"`
	PrivatePort   *NetworkPort    `json:"privatePort"`
}

func NewPortBinding(
	publicPort NetworkPort,
	containerPort NetworkPort,
	protocol NetworkProtocol,
	privatePort *NetworkPort,
) PortBinding {
	return PortBinding{
		PublicPort:    publicPort,
		ContainerPort: containerPort,
		Protocol:      protocol,
		PrivatePort:   privatePort,
	}
}

// format: publicPort[:containerPort][/protocol][:privatePort]
func NewPortBindingFromString(value string) (PortBinding, error) {
	var portBinding PortBinding

	bindingParts := strings.Split(value, ":")
	if len(bindingParts) == 0 {
		return portBinding, errors.New("InvalidPortBinding")
	}

	publicPort, err := NewNetworkPort(bindingParts[0])
	if err != nil {
		return portBinding, err
	}

	protocol := GuessNetworkProtocolByPort(publicPort)
	if len(bindingParts) == 1 {
		return NewPortBinding(
			publicPort,
			publicPort,
			protocol,
			nil,
		), nil
	}

	containerPortProtocolParts := strings.Split(bindingParts[1], "/")

	containerPortStr := containerPortProtocolParts[0]
	containerPort, err := NewNetworkPort(containerPortStr)
	if err != nil {
		return portBinding, err
	}

	protocolStr := protocol.String()
	if len(containerPortProtocolParts) == 2 {
		protocolStr = containerPortProtocolParts[1]
	}
	protocol, err = NewNetworkProtocol(protocolStr)
	if err != nil {
		return portBinding, err
	}

	if len(bindingParts) == 2 {
		return NewPortBinding(
			publicPort,
			containerPort,
			protocol,
			nil,
		), nil
	}

	privatePort, err := NewNetworkPort(bindingParts[2])
	if err != nil {
		return portBinding, err
	}

	return NewPortBinding(
		publicPort,
		containerPort,
		protocol,
		&privatePort,
	), nil
}

func (portBinding PortBinding) GetPublicPort() NetworkPort {
	return portBinding.PublicPort
}

func (portBinding PortBinding) GetContainerPort() NetworkPort {
	return portBinding.ContainerPort
}

func (portBinding PortBinding) GetProtocol() NetworkProtocol {
	return portBinding.Protocol
}

func (portBinding PortBinding) String() string {
	containerPortWithProtocol := portBinding.ContainerPort.String() +
		"/" + portBinding.Protocol.String()

	stringParts := []string{
		portBinding.PublicPort.String(),
		containerPortWithProtocol,
	}

	if portBinding.PrivatePort != nil {
		stringParts = append(stringParts, portBinding.PrivatePort.String())
	}

	return strings.Join(stringParts, ":")
}
