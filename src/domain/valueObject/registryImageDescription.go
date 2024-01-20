package valueObject

import "errors"

type RegistryImageDescription string

func NewRegistryImageDescription(value string) (RegistryImageDescription, error) {
	description := RegistryImageDescription(value)
	if !description.isValid() {
		return "", errors.New("InvalidRegistryImageDescription")
	}
	return description, nil
}

func NewRegistryImageDescriptionPanic(value string) RegistryImageDescription {
	description, err := NewRegistryImageDescription(value)
	if err != nil {
		panic(err)
	}
	return description
}

func (description RegistryImageDescription) isValid() bool {
	isTooShort := len(string(description)) < 1
	isTooLong := len(string(description)) > 1024
	return !isTooShort && !isTooLong
}

func (description RegistryImageDescription) String() string {
	return string(description)
}
