package valueObject

import "errors"

type ContainerEntrypoint string

func NewContainerEntrypoint(value string) (ContainerEntrypoint, error) {
	ce := ContainerEntrypoint(value)
	if !ce.isValid() {
		return "", errors.New("InvalidContainerEntrypoint")
	}
	return ce, nil
}

func NewContainerEntrypointPanic(value string) ContainerEntrypoint {
	ce := ContainerEntrypoint(value)
	if !ce.isValid() {
		panic("InvalidContainerEntrypoint")
	}
	return ce
}

func (ce ContainerEntrypoint) isValid() bool {
	isTooShort := len(string(ce)) < 6
	isTooLong := len(string(ce)) > 1000
	return !isTooShort && !isTooLong
}

func (ce ContainerEntrypoint) String() string {
	return string(ce)
}
