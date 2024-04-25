package valueObject

import (
	"errors"
	"regexp"
	"strings"
)

const containerIdRegex string = `^\w{12,64}$`

type ContainerId string

func NewContainerId(value interface{}) (ContainerId, error) {
	stringValue, assertOk := value.(string)
	if !assertOk {
		return "", errors.New("ContainerIdMustBeString")
	}

	stringValue = strings.TrimSpace(stringValue)
	stringValue = strings.ToLower(stringValue)

	re := regexp.MustCompile(containerIdRegex)
	isValid := re.MatchString(stringValue)
	if !isValid {
		return "", errors.New("InvalidContainerId")
	}

	return ContainerId(stringValue), nil
}

func NewContainerIdPanic(value string) ContainerId {
	containerId, err := NewContainerId(value)
	if err != nil {
		panic(err)
	}
	return containerId
}

func (containerId ContainerId) String() string {
	return string(containerId)
}
