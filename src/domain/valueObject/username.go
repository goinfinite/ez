package valueObject

import (
	"errors"
	"regexp"
	"strings"

	voHelper "github.com/speedianet/control/src/domain/valueObject/helper"
)

const usernameRegex string = `^[a-z_]([a-z0-9_-]{0,31}|[a-z0-9_-]{0,30}\$)$`

type Username string

func NewUsername(value interface{}) (Username, error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return "", errors.New("UsernameMustBeString")
	}

	stringValue = strings.TrimSpace(stringValue)
	stringValue = strings.ToLower(stringValue)

	re := regexp.MustCompile(usernameRegex)
	isValid := re.MatchString(stringValue)
	if !isValid {
		return "", errors.New("InvalidUsername")
	}

	return Username(stringValue), nil
}

func (vo Username) String() string {
	return string(vo)
}
