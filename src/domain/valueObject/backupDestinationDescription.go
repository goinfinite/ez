package valueObject

import (
	"errors"

	voHelper "github.com/goinfinite/ez/src/domain/valueObject/helper"
)

type BackupDestinationDescription string

func NewBackupDestinationDescription(value interface{}) (
	description BackupDestinationDescription, err error,
) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return description, errors.New("BackupDestinationDescriptionMustBeString")
	}

	if len(stringValue) < 2 {
		return description, errors.New("BackupDestinationDescriptionTooSmall")
	}

	if len(stringValue) > 2048 {
		return description, errors.New("BackupDestinationDescriptionTooBig")
	}

	return BackupDestinationDescription(stringValue), nil
}

func (vo BackupDestinationDescription) String() string {
	return string(vo)
}
