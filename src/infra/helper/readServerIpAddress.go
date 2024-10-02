package infraHelper

import (
	"errors"
	"strings"

	"github.com/goinfinite/ez/src/domain/valueObject"
)

func ReadServerPrivateIpAddress() (ipAddress valueObject.IpAddress, err error) {
	rawRecord, err := RunCmdWithSubShell("hostname -I")
	if err != nil {
		return ipAddress, err
	}

	rawRecord = strings.TrimSpace(rawRecord)
	if rawRecord == "" {
		return ipAddress, errors.New("PrivateIpAddressNotFound")
	}

	return valueObject.NewIpAddress(rawRecord)
}

func ReadServerPublicIpAddress() (ipAddress valueObject.IpAddress, err error) {
	digCmd := "dig +short TXT"
	rawRecord, err := RunCmdWithSubShell(
		digCmd + " o-o.myaddr.l.google.com @ns1.google.com",
	)
	if err != nil || rawRecord == "" {
		rawRecord, err = RunCmdWithSubShell(
			digCmd + "CH whoami.cloudflare @1.1.1.1",
		)
		if err != nil {
			return ipAddress, err
		}
	}

	rawRecord = strings.Trim(rawRecord, `"`)
	rawRecord = strings.TrimSpace(rawRecord)
	if rawRecord == "" {
		return ipAddress, errors.New("PublicIpAddressNotFound")
	}

	return valueObject.NewIpAddress(rawRecord)
}
