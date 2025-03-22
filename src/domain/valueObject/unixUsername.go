package valueObject

import (
	"errors"
	"regexp"
	"strings"

	voHelper "github.com/goinfinite/ez/src/domain/valueObject/helper"
)

const unixUsernameRegex string = `^[a-z_]([a-z0-9_-]{0,31}|[a-z0-9_-]{0,30}\$)$`

type UnixUsername string

func NewUnixUsername(value interface{}) (UnixUsername, error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return "", errors.New("UnixUsernameMustBeString")
	}

	stringValue = strings.ToLower(stringValue)

	re := regexp.MustCompile(unixUsernameRegex)
	if !re.MatchString(stringValue) {
		return "", errors.New("InvalidUnixUsername")
	}

	return UnixUsername(stringValue), nil
}

func (vo UnixUsername) String() string {
	return string(vo)
}
