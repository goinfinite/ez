package valueObject

import (
	"errors"
	"regexp"
	"strings"
)

const accessTokenStrRegex string = `^[a-zA-Z0-9\-_=+/.]{22,444}$`

type AccessTokenStr string

func NewAccessTokenStr(value string) (AccessTokenStr, error) {
	value = strings.TrimSpace(value)

	re := regexp.MustCompile(accessTokenStrRegex)
	isValid := re.MatchString(value)
	if !isValid {
		return "", errors.New("InvalidAccessTokenStr")
	}
	return AccessTokenStr(value), nil
}

func (ats AccessTokenStr) String() string {
	return string(ats)
}
