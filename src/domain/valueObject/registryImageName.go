package valueObject

import (
	"errors"
	"regexp"
	"strings"

	voHelper "github.com/goinfinite/ez/src/domain/valueObject/helper"
)

const registryImageNameRegex string = `^[\w\_\-]{1,128}/?[\w\.\_\-]{0,128}$`

type RegistryImageName string

func NewRegistryImageName(value interface{}) (
	imageName RegistryImageName, err error,
) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return imageName, errors.New("RegistryImageNameMustBeString")
	}

	stringValue = strings.TrimPrefix(stringValue, "localhost/")
	stringValue = strings.TrimPrefix(stringValue, "docker.io/")
	stringValue = strings.TrimSuffix(stringValue, ":latest")

	re := regexp.MustCompile(registryImageNameRegex)
	if !re.MatchString(stringValue) {
		return imageName, errors.New("InvalidRegistryImageName")
	}

	return RegistryImageName(stringValue), nil
}

func (vo RegistryImageName) String() string {
	return string(vo)
}
