package valueObject

import (
	"errors"
	"regexp"
)

const containerProfileNameRegex string = `^\w[\w -]{1,30}\w$`

type ContainerProfileName string

func NewContainerProfileName(value string) (ContainerProfileName, error) {
	user := ContainerProfileName(value)
	if !user.isValid() {
		return "", errors.New("InvalidContainerProfileName")
	}
	return user, nil
}

func NewContainerProfileNamePanic(value string) ContainerProfileName {
	user := ContainerProfileName(value)
	if !user.isValid() {
		panic("InvalidContainerProfileName")
	}
	return user
}

func (user ContainerProfileName) isValid() bool {
	re := regexp.MustCompile(containerProfileNameRegex)
	return re.MatchString(string(user))
}

func (user ContainerProfileName) String() string {
	return string(user)
}
