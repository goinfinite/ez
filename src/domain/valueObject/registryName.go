package valueObject

import (
	"errors"
	"regexp"
)

const registryNameRegex string = `^\w{1,128}$`

type RegistryName string

func NewRegistryName(value string) (RegistryName, error) {
	name := RegistryName(value)
	if !name.isValid() {
		return "", errors.New("InvalidRegistryName")
	}
	return name, nil
}

func NewRegistryNamePanic(value string) RegistryName {
	name, err := NewRegistryName(value)
	if err != nil {
		panic(err)
	}
	return name
}

func (name RegistryName) isValid() bool {
	re := regexp.MustCompile(registryNameRegex)
	return re.MatchString(string(name))
}

func (name RegistryName) String() string {
	return string(name)
}
