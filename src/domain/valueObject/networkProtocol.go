package valueObject

import (
	"errors"
	"slices"
	"strings"
)

type NetworkProtocol string

var validNetworkProtocols = []string{
	"tcp",
	"udp",
}

func NewNetworkProtocol(value string) (NetworkProtocol, error) {
	value = strings.TrimSpace(value)
	value = strings.ToLower(value)
	if !slices.Contains(validNetworkProtocols, value) {
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
