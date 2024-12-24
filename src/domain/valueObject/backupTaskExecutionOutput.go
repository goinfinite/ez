package valueObject

import (
	"errors"

	voHelper "github.com/goinfinite/ez/src/domain/valueObject/helper"
)

type BackupTaskExecutionOutput string

func NewBackupTaskExecutionOutput(value interface{}) (BackupTaskExecutionOutput, error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return "", errors.New("BackupTaskExecutionOutputMustBeString")
	}

	valueLength := len(stringValue)
	if valueLength > 2048 {
		stringValue = stringValue[:2048]
	}

	return BackupTaskExecutionOutput(stringValue), nil
}

func (vo BackupTaskExecutionOutput) String() string {
	return string(vo)
}
