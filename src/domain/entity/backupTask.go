package entity

import (
	"github.com/goinfinite/ez/src/domain/valueObject"
)

type BackupTask struct {
	TaskId          valueObject.BackupTaskId               `json:"taskId"`
	AccountId       valueObject.AccountId                  `json:"accountId"`
	JobId           valueObject.BackupJobId                `json:"jobId"`
	DestinationId   valueObject.BackupDestinationId        `json:"destinationId"`
	TaskStatus      valueObject.BackupTaskStatus           `json:"taskStatus"`
	BackupType      valueObject.BackupJobType              `json:"backupType"`
	BackupSchedule  valueObject.CronSchedule               `json:"backupSchedule"`
	TimeoutSecs     *uint64                                `json:"timeoutSecs"`
	ContainerIds    []valueObject.ContainerId              `json:"containerIds"`
	ExecutionOutput *valueObject.BackupTaskExecutionOutput `json:"executionOutput"`
	StartedAt       *valueObject.UnixTime                  `json:"startedAt"`
	FinishedAt      *valueObject.UnixTime                  `json:"finishedAt"`
	ElapsedSecs     *uint64                                `json:"elapsedSecs"`
	CreatedAt       valueObject.UnixTime                   `json:"createdAt"`
}

func NewBackupTask(
	taskId valueObject.BackupTaskId,
	accountId valueObject.AccountId,
	jobId valueObject.BackupJobId,
	destinationId valueObject.BackupDestinationId,
	taskStatus valueObject.BackupTaskStatus,
	backupType valueObject.BackupJobType,
	backupSchedule valueObject.CronSchedule,
	timeoutSecs *uint64,
	containerIds []valueObject.ContainerId,
	executionOutput *valueObject.BackupTaskExecutionOutput,
	startedAt *valueObject.UnixTime,
	finishedAt *valueObject.UnixTime,
	elapsedSecs *uint64,
	createdAt valueObject.UnixTime,
) BackupTask {
	return BackupTask{
		TaskId:          taskId,
		AccountId:       accountId,
		JobId:           jobId,
		DestinationId:   destinationId,
		TaskStatus:      taskStatus,
		BackupType:      backupType,
		BackupSchedule:  backupSchedule,
		TimeoutSecs:     timeoutSecs,
		ContainerIds:    containerIds,
		ExecutionOutput: executionOutput,
		StartedAt:       startedAt,
		FinishedAt:      finishedAt,
		ElapsedSecs:     elapsedSecs,
		CreatedAt:       createdAt,
	}
}
