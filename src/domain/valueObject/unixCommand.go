package valueObject

import (
	"errors"
	"strings"

	voHelper "github.com/goinfinite/ez/src/domain/valueObject/helper"
)

type UnixCommand string

func NewUnixCommand(value interface{}) (UnixCommand, error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return "", errors.New("UnixCommandMustBeString")
	}

	stringValue = strings.TrimSpace(stringValue)

	valueLength := len(stringValue)
	switch {
	case valueLength == 0:
		return "", errors.New("UnixCommandEmpty")
	case valueLength < 3:
		return "", errors.New("UnixCommandTooShort")
	case valueLength > 4096:
		return "", errors.New("UnixCommandTooLong")
	}

	return UnixCommand(stringValue), nil
}

func (vo UnixCommand) String() string {
	return string(vo)
}
