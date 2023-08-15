package valueObject

import (
	"errors"
	"regexp"
)

const containerEnvRegex string = `^\w{1,1000}=.{1,1000}$`

type ContainerEnv string

func NewContainerEnv(value string) (ContainerEnv, error) {
	containerEnv := ContainerEnv(value)
	if !containerEnv.isValid() {
		return "", errors.New("InvalidContainerEnv")
	}
	return containerEnv, nil
}

func NewContainerEnvPanic(value string) ContainerEnv {
	containerEnv := ContainerEnv(value)
	if !containerEnv.isValid() {
		panic("InvalidContainerEnv")
	}
	return containerEnv
}

func (containerEnv ContainerEnv) isValid() bool {
	re := regexp.MustCompile(containerEnvRegex)
	return re.MatchString(string(containerEnv))
}

func (containerEnv ContainerEnv) String() string {
	return string(containerEnv)
}
