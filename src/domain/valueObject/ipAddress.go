package valueObject

import (
	"errors"
	"net"
	"strings"

	voHelper "github.com/speedianet/control/src/domain/valueObject/helper"
)

type IpAddress string

func NewIpAddress(value interface{}) (IpAddress, error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return "", errors.New("IpAddressMustBeString")
	}

	stringValue = strings.TrimSpace(stringValue)
	isValid := net.ParseIP(stringValue) != nil
	if !isValid {
		return "", errors.New("InvalidIpAddress")
	}

	return IpAddress(stringValue), nil
}

func NewLocalhostIpAddress() IpAddress {
	return IpAddress("127.0.0.1")
}

func (vo IpAddress) String() string {
	return string(vo)
}
