package valueObject

import (
	"errors"
	"regexp"
)

const registryImageNameRegex string = `^[\w\_\-]{1,128}/?[\w\_\-]{0,128}$`

type RegistryImageName string

func NewRegistryImageName(value string) (RegistryImageName, error) {
	imgName := RegistryImageName(value)
	if !imgName.isValid() {
		return "", errors.New("InvalidRegistryImageName")
	}
	return imgName, nil
}

func NewRegistryImageNamePanic(value string) RegistryImageName {
	imgName, err := NewRegistryImageName(value)
	if err != nil {
		panic(err)
	}
	return imgName
}

func (imgName RegistryImageName) isValid() bool {
	re := regexp.MustCompile(registryImageNameRegex)
	return re.MatchString(string(imgName))
}

func (imgName RegistryImageName) String() string {
	return string(imgName)
}
