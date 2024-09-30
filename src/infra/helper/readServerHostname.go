package infraHelper

import (
	"errors"
	"strings"

	"github.com/goinfinite/ez/src/domain/valueObject"
)

func ReadServerHostname() (hostname valueObject.Fqdn, err error) {
	rawHostname, err := RunCmd("hostname")
	if err != nil || rawHostname == "" {
		return hostname, errors.New("ServerHostnameNotFound")
	}

	rawHostname = strings.TrimSpace(rawHostname)

	return valueObject.NewFqdn(rawHostname)
}
