package valueObject

import (
	"errors"
	"regexp"
)

const serviceNameRegex string = `^[\w\.\_\-]{1,128}$`

type ServiceName string

func NewServiceName(value string) (ServiceName, error) {
	name := ServiceName(value)
	if !name.isValid() {
		return "", errors.New("InvalidServiceName")
	}
	return name, nil
}

func NewServiceNamePanic(value string) ServiceName {
	name, err := NewServiceName(value)
	if err != nil {
		panic(err)
	}
	return name
}

func (name ServiceName) isValid() bool {
	re := regexp.MustCompile(serviceNameRegex)
	return re.MatchString(string(name))
}

func (name ServiceName) String() string {
	return string(name)
}
