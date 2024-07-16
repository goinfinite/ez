package valueObject

import (
	"errors"
	"slices"

	voHelper "github.com/speedianet/control/src/domain/valueObject/helper"
)

type ActivityRecordCode string

var ValidActivityRecordCodes = []string{
	"LoginFailed", "LoginSuccessful",
	"AccountCreated", "AccountDeleted",
	"AccountPasswordUpdated", "AccountApiKeyUpdated", "AccountQuotaUpdated",
	"UnauthorizedAccess",
}

func NewActivityRecordCode(value interface{}) (ActivityRecordCode, error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return "", errors.New("ActivityRecordCodeMustBeString")
	}

	if !slices.Contains(ValidActivityRecordCodes, stringValue) {
		return "", errors.New("InvalidActivityRecordCode")
	}

	return ActivityRecordCode(stringValue), nil
}

func (vo ActivityRecordCode) String() string {
	return string(vo)
}
