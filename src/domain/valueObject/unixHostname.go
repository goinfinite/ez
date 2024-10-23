package valueObject

import (
	"errors"
	"regexp"
	"strings"

	voHelper "github.com/goinfinite/ez/src/domain/valueObject/helper"
)

const unixHostnameRegex string = `^[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`

type UnixHostname string

func NewUnixHostname(value interface{}) (
	hostname UnixHostname, err error,
) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return hostname, errors.New("UnixHostnameMustBeString")
	}

	stringValue = strings.ToLower(stringValue)

	re := regexp.MustCompile(unixHostnameRegex)
	if !re.MatchString(stringValue) {
		return hostname, errors.New("InvalidUnixHostname")
	}

	return UnixHostname(stringValue), nil
}

func (vo UnixHostname) String() string {
	return string(vo)
}
