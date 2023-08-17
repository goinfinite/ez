package valueObject

import (
	"errors"
	"regexp"
)

const containerIdRegex string = `^\w{12,64}$`

type ContainerId string

func NewContainerId(value string) (ContainerId, error) {
	containerId := ContainerId(value)
	if !containerId.isValid() {
		return "", errors.New("InvalidContainerId")
	}
	return containerId, nil
}

func NewContainerIdPanic(value string) ContainerId {
	containerId := ContainerId(value)
	if !containerId.isValid() {
		panic("InvalidContainerId")
	}
	return containerId
}

func (containerId ContainerId) isValid() bool {
	re := regexp.MustCompile(containerIdRegex)
	return re.MatchString(string(containerId))
}

func (containerId ContainerId) String() string {
	return string(containerId)
}
