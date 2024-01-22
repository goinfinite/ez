package valueObject

import (
	"errors"
	"regexp"
)

const registryImageTagRegex string = `^[\w\.\_\-]{1,128}$`

type RegistryImageTag string

func NewRegistryImageTag(value string) (RegistryImageTag, error) {
	tag := RegistryImageTag(value)
	if !tag.isValid() {
		return "", errors.New("InvalidRegistryImageTag")
	}
	return tag, nil
}

func NewRegistryImageTagPanic(value string) RegistryImageTag {
	tag, err := NewRegistryImageTag(value)
	if err != nil {
		panic(err)
	}
	return tag
}

func (tag RegistryImageTag) isValid() bool {
	re := regexp.MustCompile(registryImageTagRegex)
	return re.MatchString(string(tag))
}

func (tag RegistryImageTag) String() string {
	return string(tag)
}
