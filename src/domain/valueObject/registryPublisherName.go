package valueObject

import (
	"errors"
	"regexp"
)

const registryPublisherNameRegex string = `^\w{1,128}$`

type RegistryPublisherName string

func NewRegistryPublisherName(value string) (RegistryPublisherName, error) {
	name := RegistryPublisherName(value)
	if !name.isValid() {
		return "", errors.New("InvalidRegistryPublisherName")
	}
	return name, nil
}

func NewRegistryPublisherNamePanic(value string) RegistryPublisherName {
	name, err := NewRegistryPublisherName(value)
	if err != nil {
		panic(err)
	}
	return name
}

func (name RegistryPublisherName) isValid() bool {
	re := regexp.MustCompile(registryPublisherNameRegex)
	return re.MatchString(string(name))
}

func (name RegistryPublisherName) String() string {
	return string(name)
}
