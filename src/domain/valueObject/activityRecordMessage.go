package valueObject

import (
	"errors"

	voHelper "github.com/speedianet/control/src/domain/valueObject/helper"
)

type ActivityRecordMessage string

func NewActivityRecordMessage(value interface{}) (ActivityRecordMessage, error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return "", errors.New("ActivityRecordMessageMustBeString")
	}

	valueLength := len(stringValue)
	if valueLength > 2048 {
		stringValue = stringValue[:2048]
	}

	return ActivityRecordMessage(stringValue), nil
}

func (vo ActivityRecordMessage) String() string {
	return string(vo)
}
