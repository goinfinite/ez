package valueObject

import (
	"errors"
	"net"
	"regexp"
	"strings"
)

const fqdnRegex string = `^((\*\.)?([a-zA-Z0-9_]+[\w-]*\.)*)?([a-zA-Z0-9_][\w-]*[a-zA-Z0-9])\.([a-zA-Z]{2,})$`

type Fqdn string

func NewFqdn(value interface{}) (Fqdn, error) {
	stringValue, assertOk := value.(string)
	if !assertOk {
		return "", errors.New("FqdnMustBeString")
	}

	stringValue = strings.TrimSpace(stringValue)
	stringValue = strings.ToLower(stringValue)

	isIpAddress := net.ParseIP(stringValue) != nil
	if isIpAddress {
		return "", errors.New("FqdnCannotBeIpAddress")
	}

	re := regexp.MustCompile(fqdnRegex)
	isValid := re.MatchString(stringValue)
	if !isValid {
		return "", errors.New("InvalidFqdn")
	}

	return Fqdn(stringValue), nil
}

func NewFqdnPanic(value string) Fqdn {
	fqdn, err := NewFqdn(value)
	if err != nil {
		panic(err)
	}
	return fqdn
}

func (fqdn Fqdn) String() string {
	return string(fqdn)
}
