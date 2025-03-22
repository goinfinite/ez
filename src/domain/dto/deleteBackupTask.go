package dto

import (
	"github.com/goinfinite/ez/src/domain/valueObject"
)

type DeleteBackupTask struct {
	TaskId             valueObject.BackupTaskId `json:"taskId"`
	ShouldDiscardFiles bool                     `json:"shouldDiscardFiles"`
	OperatorAccountId  valueObject.AccountId    `json:"-"`
	OperatorIpAddress  valueObject.IpAddress    `json:"-"`
}

func NewDeleteBackupTask(
	taskId valueObject.BackupTaskId,
	shouldDiscardFiles bool,
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
) DeleteBackupTask {
	return DeleteBackupTask{
		TaskId:             taskId,
		ShouldDiscardFiles: shouldDiscardFiles,
		OperatorAccountId:  operatorAccountId,
		OperatorIpAddress:  operatorIpAddress,
	}
}
