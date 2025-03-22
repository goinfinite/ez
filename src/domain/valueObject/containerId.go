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
	stringValue = strings.ToLower(stringValue)

	re := regexp.MustCompile(containerIdRegex)
	if !re.MatchString(stringValue) {
		return "", errors.New("InvalidContainerId")
	}

	if len(stringValue) > 12 {
		stringValue = stringValue[:12]
	}

	return ContainerId(stringValue), nil
}

func (vo ContainerId) String() string {
	return string(vo)
}
