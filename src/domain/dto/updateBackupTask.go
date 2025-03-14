package dto

import (
	"github.com/goinfinite/ez/src/domain/valueObject"
)

type UpdateBackupTask struct {
	TaskId            valueObject.BackupTaskId      `json:"jobId"`
	TaskStatus        *valueObject.BackupTaskStatus `json:"taskStatus,omitempty"`
	OperatorAccountId valueObject.AccountId         `json:"-"`
	OperatorIpAddress valueObject.IpAddress         `json:"-"`
}

func NewUpdateBackupTask(
	taskId valueObject.BackupTaskId,
	taskStatus *valueObject.BackupTaskStatus,
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
) UpdateBackupTask {
	return UpdateBackupTask{
		TaskId:            taskId,
		TaskStatus:        taskStatus,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
	}
}
