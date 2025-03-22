package valueObject

import (
	"errors"
	"net"
	"strings"

	voHelper "github.com/goinfinite/ez/src/domain/valueObject/helper"
)

type NetworkHost string

func NewNetworkHost(value interface{}) (NetworkHost, error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return "", errors.New("NetworkHostMustBeString")
	}
	stringValue = strings.ToLower(stringValue)

	_, err = NewFqdn(stringValue)
	if err != nil {
		if net.ParseIP(stringValue) == nil {
			return "", errors.New("InvalidNetworkHost")
		}
	}

	return NetworkHost(stringValue), nil
}

func (vo NetworkHost) String() string {
	return string(vo)
}
