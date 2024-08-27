package valueObject

import (
	"errors"
	"regexp"

	voHelper "github.com/speedianet/control/src/domain/valueObject/helper"
)

const containerEnvRegex string = `^\w{1,1000}=.{1,1000}$`

type ContainerEnv string

func NewContainerEnv(value interface{}) (env ContainerEnv, err error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return env, errors.New("ContainerEnvMustBeString")
	}

	re := regexp.MustCompile(containerEnvRegex)
	if !re.MatchString(stringValue) {
		return env, errors.New("InvalidContainerEnv")
	}

	return ContainerEnv(stringValue), nil
}

func (vo ContainerEnv) String() string {
	return string(vo)
}
