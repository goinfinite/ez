package valueObject

import (
	"errors"
	"slices"
	"strings"

	voHelper "github.com/goinfinite/ez/src/domain/valueObject/helper"
)

type NetworkProtocol string

var ValidNetworkProtocols = []string{
	"http", "https", "ws", "wss", "grpc", "grpcs", "tcp", "udp",
}

func NewNetworkProtocol(value interface{}) (NetworkProtocol, error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return "", errors.New("NetworkProtocolMustBeString")
	}

	stringValue = strings.ToLower(stringValue)
	if !slices.Contains(ValidNetworkProtocols, stringValue) {
		return "", errors.New("InvalidNetworkProtocol")
	}

	return NetworkProtocol(stringValue), nil
}

func (vo NetworkProtocol) String() string {
	return string(vo)
}

func GuessNetworkProtocolByPort(port NetworkPort) NetworkProtocol {
	protocolStr := "tcp"
	switch port.Uint16() {
	case 53, 123, 514:
		protocolStr = "udp"
	case 80, 2368, 3000, 3001, 3002, 3003, 3004, 3005, 3006, 3007, 3008, 3009, 5000, 5601, 8000, 8001, 8002, 8065, 8080, 8081:
		protocolStr = "http"
	case 443, 1618, 8443, 8444, 8445:
		protocolStr = "https"
	}

	networkProtocol, _ := NewNetworkProtocol(protocolStr)
	return networkProtocol
}
