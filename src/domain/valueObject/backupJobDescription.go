package valueObject

import (
	"errors"

	voHelper "github.com/goinfinite/ez/src/domain/valueObject/helper"
)

type BackupJobDescription string

func NewBackupJobDescription(value interface{}) (
	description BackupJobDescription, err error,
) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return description, errors.New("BackupJobDescriptionMustBeString")
	}

	if len(stringValue) < 2 {
		return description, errors.New("BackupJobDescriptionTooSmall")
	}

	if len(stringValue) > 2048 {
		return description, errors.New("BackupJobDescriptionTooBig")
	}

	return BackupJobDescription(stringValue), nil
}

func (vo BackupJobDescription) String() string {
	return string(vo)
}
