package valueObject

import (
	"errors"
	"regexp"
	"strings"

	voHelper "github.com/speedianet/control/src/domain/valueObject/helper"
)

const containerImageIdRegex string = `^\w{12,64}$`

type ContainerImageId string

func NewContainerImageId(value interface{}) (ContainerImageId, error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return "", errors.New("ContainerImageIdMustBeString")
	}
	stringValue = strings.ToLower(stringValue)

	re := regexp.MustCompile(containerImageIdRegex)
	if !re.MatchString(stringValue) {
		return "", errors.New("InvalidContainerImageId")
	}

	return ContainerImageId(stringValue), nil
}

func (vo ContainerImageId) String() string {
	return string(vo)
}
