package valueObject

import "errors"

type NetworkProtocol string

const (
	tcp NetworkProtocol = "tcp"
	udp NetworkProtocol = "udp"
)

func NewNetworkProtocol(value string) (NetworkProtocol, error) {
	np := NetworkProtocol(value)
	if !np.isValid() {
		return "", errors.New("InvalidNetworkProtocol")
	}
	return np, nil
}

func NewNetworkProtocolPanic(value string) NetworkProtocol {
	np := NetworkProtocol(value)
	if !np.isValid() {
		panic("InvalidNetworkProtocol")
	}
	return np
}

func (np NetworkProtocol) isValid() bool {
	switch np {
	case tcp, udp:
		return true
	default:
		return false
	}
}

func (np NetworkProtocol) String() string {
	return string(np)
}
