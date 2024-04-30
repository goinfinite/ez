package valueObject

import (
	"errors"
	"regexp"
	"strings"
)

const serviceNameRegex string = `^[a-z][\w\-]{0,128}$`

type ServiceName string

func NewServiceName(value string) (ServiceName, error) {
	value = strings.TrimSpace(value)
	value = strings.ToLower(value)
	specialCharReplacer := strings.NewReplacer(" ", "-", "_", "-", ".", "-")
	value = specialCharReplacer.Replace(value)

	re := regexp.MustCompile(serviceNameRegex)
	isValid := re.MatchString(value)
	if !isValid {
		return "", errors.New("InvalidServiceName")
	}

	return ServiceName(value), nil
}

func NewServiceNamePanic(value string) ServiceName {
	name, err := NewServiceName(value)
	if err != nil {
		panic(err)
	}
	return name
}

func (name ServiceName) String() string {
	return string(name)
}
