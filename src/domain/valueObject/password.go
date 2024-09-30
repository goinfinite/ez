package valueObject

import (
	"errors"
	"strings"

	voHelper "github.com/goinfinite/ez/src/domain/valueObject/helper"
)

type Password string

func NewPassword(value interface{}) (Password, error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return "", errors.New("PasswordMustBeString")
	}

	stringValue = strings.TrimSpace(stringValue)
	valueLength := len(stringValue)
	switch {
	case valueLength == 0:
		return "", errors.New("PasswordEmpty")
	case valueLength < 5:
		return "", errors.New("PasswordTooShort")
	case valueLength > 128:
		return "", errors.New("PasswordTooLong")
	}

	return Password(stringValue), nil
}

func (vo Password) String() string {
	return string(vo)
}
