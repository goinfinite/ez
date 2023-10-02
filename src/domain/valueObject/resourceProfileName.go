package valueObject

import (
	"errors"
	"regexp"
)

const resourceProfileNameRegex string = `^\w[\w -]{1,30}\w$`

type ResourceProfileName string

func NewResourceProfileName(value string) (ResourceProfileName, error) {
	user := ResourceProfileName(value)
	if !user.isValid() {
		return "", errors.New("InvalidResourceProfileName")
	}
	return user, nil
}

func NewResourceProfileNamePanic(value string) ResourceProfileName {
	user := ResourceProfileName(value)
	if !user.isValid() {
		panic("InvalidResourceProfileName")
	}
	return user
}

func (user ResourceProfileName) isValid() bool {
	re := regexp.MustCompile(resourceProfileNameRegex)
	return re.MatchString(string(user))
}

func (user ResourceProfileName) String() string {
	return string(user)
}
