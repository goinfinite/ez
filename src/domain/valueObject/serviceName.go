package valueObject

import "errors"

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
	isTooShort := len(string(name)) < 3
	isTooLong := len(string(name)) > 128
	return !isTooShort && !isTooLong
}

func (name ServiceName) String() string {
	return string(name)
}
