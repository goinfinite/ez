package valueObject

import (
	"errors"
	"regexp"
	"strings"

	voHelper "github.com/goinfinite/ez/src/domain/valueObject/helper"
)

const containerIdRegex string = `^\w{12,64}$`

type ContainerId string

func NewContainerId(value interface{}) (ContainerId, error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
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
