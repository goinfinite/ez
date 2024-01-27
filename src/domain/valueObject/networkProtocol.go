package valueObject

import (
	"errors"
	"strings"

	"golang.org/x/exp/slices"
)

type NetworkProtocol string

var ValidNetworkProtocols = []string{
	"http",
	"https",
	"ws",
	"wss",
	"grpc",
	"grpcs",
	"tcp",
	"udp",
}

func NewNetworkProtocol(value string) (NetworkProtocol, error) {
	value = strings.TrimSpace(value)
	value = strings.ToLower(value)
	if !slices.Contains(ValidNetworkProtocols, value) {
		return "", errors.New("InvalidNetworkProtocol")
	}
	return NetworkProtocol(value), nil
}

func NewNetworkProtocolPanic(value string) NetworkProtocol {
	np, err := NewNetworkProtocol(value)
	if err != nil {
		panic(err)
	}
	return np
}

func (np NetworkProtocol) String() string {
	return string(np)
}

func GuessNetworkProtocolByPort(port NetworkPort) NetworkProtocol {
	protocolStr := "tcp"
	switch port.Get() {
	case 53, 123, 514:
		protocolStr = "udp"
	case 80, 2368, 3000, 5000, 5601, 8000, 8001, 8002, 8065:
		protocolStr = "http"
	case 443, 1618, 8443, 8444, 8445:
		protocolStr = "https"
	}

	networkProtocol, _ := NewNetworkProtocol(protocolStr)
	return networkProtocol
}
