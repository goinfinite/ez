package valueObject

import (
	"errors"
	"strings"

	voHelper "github.com/goinfinite/ez/src/domain/valueObject/helper"
)

type ContainerEntrypoint string

func NewContainerEntrypoint(value interface{}) (ContainerEntrypoint, error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return "", errors.New("ContainerEntrypointMustBeString")
	}

	stringValue = strings.TrimSpace(stringValue)
	stringValue = strings.ToLower(stringValue)

	if len(stringValue) < 6 {
		return "", errors.New("ContainerEntrypointIsTooShort")
	}

	if len(stringValue) > 1000 {
		return "", errors.New("ContainerEntrypointIsTooLong")
	}

	return ContainerEntrypoint(stringValue), nil
}

func NewContainerEntrypointPanic(value string) ContainerEntrypoint {
	ce, err := NewContainerEntrypoint(value)
	if err != nil {
		panic(err)
	}
	return ce
}

func (ce ContainerEntrypoint) String() string {
	return string(ce)
}
