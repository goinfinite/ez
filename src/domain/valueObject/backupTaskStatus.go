package valueObject

import (
	"errors"
	"strings"

	voHelper "github.com/goinfinite/ez/src/domain/valueObject/helper"
)

type BackupTaskStatus string

const (
	BackupTaskStatusCompleted BackupTaskStatus = "completed"
	BackupTaskStatusFailed    BackupTaskStatus = "failed"
	BackupTaskStatusExecuting BackupTaskStatus = "executing"
	BackupTaskStatusPartial   BackupTaskStatus = "partial"
	BackupTaskStatusCanceled  BackupTaskStatus = "canceled"
	BackupTaskStatusCancelled BackupTaskStatus = "cancelled"
)

func NewBackupTaskStatus(value interface{}) (
	taskStatus BackupTaskStatus, err error,
) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return taskStatus, errors.New("BackupTaskStatusMustBeString")
	}
	stringValue = strings.ToLower(stringValue)

	stringValueVo := BackupTaskStatus(stringValue)
	switch stringValueVo {
	case BackupTaskStatusCompleted, BackupTaskStatusFailed,
		BackupTaskStatusExecuting, BackupTaskStatusPartial, BackupTaskStatusCancelled:
		return stringValueVo, nil
	case BackupTaskStatusCanceled:
		return BackupTaskStatusCancelled, nil
	default:
		return taskStatus, errors.New("InvalidBackupTaskStatus")
	}
}

func (vo BackupTaskStatus) String() string {
	return string(vo)
}
