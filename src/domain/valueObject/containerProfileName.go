package valueObject

import (
	"errors"
	"regexp"

	voHelper "github.com/goinfinite/ez/src/domain/valueObject/helper"
)

const containerProfileNameRegex string = `^\w[\p{L}0-9\ \-]{1,64}\w$`

type ContainerProfileName string

func NewContainerProfileName(value any) (name ContainerProfileName, err error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return name, errors.New("ContainerProfileNameMustBeString")
	}

	re := regexp.MustCompile(containerProfileNameRegex)
	if !re.MatchString(stringValue) {
		return name, errors.New("InvalidContainerProfileName")
	}

	return ContainerProfileName(stringValue), nil
}

func (vo ContainerProfileName) String() string {
	return string(vo)
}
