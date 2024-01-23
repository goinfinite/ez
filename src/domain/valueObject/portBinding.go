package valueObject

import (
	"strings"

	voHelper "github.com/speedianet/control/src/domain/valueObject/helper"
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

	portBindingRegex := `^(?P<publicPort>\d{1,5})(?::(?P<containerPort>\d{1,5}))?(?:\/(?P<protocol>\w{1,5}))?(?::(?P<privatePort>\d+))?$`
	portBindingParts := voHelper.FindNamedGroupsMatches(portBindingRegex, string(value))

	publicPort, err := NewNetworkPort(portBindingParts["publicPort"])
	if err != nil {
		return portBinding, err
	}

	if portBindingParts["containerPort"] == "" {
		portBindingParts["containerPort"] = portBindingParts["publicPort"]
	}

	containerPort, err := NewNetworkPort(portBindingParts["containerPort"])
	if err != nil {
		return portBinding, err
	}

	protocol := GuessNetworkProtocolByPort(publicPort)
	if portBindingParts["protocol"] != "" && protocol.String() == "tcp" {
		protocol, err = NewNetworkProtocol(portBindingParts["protocol"])
		if err != nil {
			return portBinding, err
		}
	}

	var privatePortPtr *NetworkPort
	if portBindingParts["privatePort"] != "" {
		privatePort, err := NewNetworkPort(portBindingParts["privatePort"])
		if err != nil {
			return portBinding, err
		}
		privatePortPtr = &privatePort
	}

	return NewPortBinding(publicPort, containerPort, protocol, privatePortPtr), nil
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
