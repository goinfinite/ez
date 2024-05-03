package valueObject

import (
	"errors"
	"regexp"
	"strings"
)

const accessTokenValueRegex string = `^[a-zA-Z0-9\-_=+/.]{22,444}$`

type AccessTokenValue string

func NewAccessTokenValue(value string) (AccessTokenValue, error) {
	value = strings.TrimSpace(value)

	re := regexp.MustCompile(accessTokenValueRegex)
	isValid := re.MatchString(value)
	if !isValid {
		return "", errors.New("InvalidAccessTokenValue")
	}
	return AccessTokenValue(value), nil
}

func (vo AccessTokenValue) String() string {
	return string(vo)
}
