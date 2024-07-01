package valueObject

import (
	"errors"
	"slices"
	"strings"

	voHelper "github.com/speedianet/control/src/domain/valueObject/helper"
)

type SecurityEventType string

var ValidSecurityEventTypes = []string{
	"failed-login",
	"successful-login",
	"account-created",
	"account-deleted",
	"account-password-updated",
	"account-api-key-updated",
	"account-quota-updated",
	"unauthorized-access",
}

func NewSecurityEventType(value interface{}) (SecurityEventType, error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return "", errors.New("SecurityEventTypeMustBeString")
	}

	stringValue = strings.TrimSpace(stringValue)
	stringValue = strings.ToLower(stringValue)

	if !slices.Contains(ValidSecurityEventTypes, stringValue) {
		return "", errors.New("InvalidSecurityEventType")
	}
	return SecurityEventType(stringValue), nil
}

func (vo SecurityEventType) String() string {
	return string(vo)
}
