package valueObject

import (
	"errors"
	"strings"

	voHelper "github.com/speedianet/control/src/domain/valueObject/helper"
)

type SecurityEventDetails string

func NewSecurityEventDetails(value interface{}) (SecurityEventDetails, error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return "", errors.New("SecurityEventDetailsMustBeString")
	}

	stringValue = strings.TrimSpace(stringValue)
	valueLength := len(stringValue)
	if valueLength > 2048 {
		stringValue = stringValue[:2048]
	}

	return SecurityEventDetails(stringValue), nil
}

func (vo SecurityEventDetails) String() string {
	return string(vo)
}
