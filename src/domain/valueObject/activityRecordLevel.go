package valueObject

import (
	"errors"
	"slices"
	"strings"

	voHelper "github.com/goinfinite/ez/src/domain/valueObject/helper"
)

type ActivityRecordLevel string

var ValidActivityRecordLevels = []string{
	"DEBUG", "INFO", "WARN", "ERROR", "SEC",
}

func NewActivityRecordLevel(value interface{}) (ActivityRecordLevel, error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return "", errors.New("ActivityRecordLevelMustBeString")
	}

	stringValue = strings.ToUpper(stringValue)

	if !slices.Contains(ValidActivityRecordLevels, stringValue) {
		switch stringValue {
		case "SECURITY":
			stringValue = "SEC"
		case "WARNING":
			stringValue = "WARN"
		default:
			return "", errors.New("InvalidActivityRecordLevel")
		}
	}

	return ActivityRecordLevel(stringValue), nil
}

func (vo ActivityRecordLevel) String() string {
	return string(vo)
}
