package dbModel

import (
	"log/slog"
	"time"

	"github.com/goinfinite/ez/src/domain/entity"
	"github.com/goinfinite/ez/src/domain/valueObject"
)

type BackupTask struct {
	ID                     uint64   `gorm:"primarykey"`
	AccountID              uint64   `gorm:"not null"`
	JobID                  uint64   `gorm:"not null"`
	DestinationID          uint64   `gorm:"not null"`
	TaskStatus             string   `gorm:"not null"`
	RetentionStrategy      string   `gorm:"not null"`
	BackupSchedule         string   `gorm:"not null"`
	TimeoutSecs            uint64   `gorm:"not null"`
	SuccessfulContainerIds []string `gorm:"serializer:json"`
	FailedContainerIds     []string `gorm:"serializer:json"`
	ExecutionOutput        *string
	StartedAt              *time.Time
	FinishedAt             *time.Time
	ElapsedSecs            *uint64
	CreatedAt              time.Time `gorm:"not null"`
	UpdatedAt              time.Time `gorm:"not null"`
}

func (model BackupTask) TableName() string {
	return "backup_tasks"
}

func NewBackupTask(
	id, accountId, jobId, destinationId uint64,
	taskStatus, retentionStrategy, backupSchedule string,
	timeoutSecs uint64,
	successfulContainerIds, failedContainerIds []string,
	executionOutput *string,
	startedAt, finishedAt *time.Time,
	elapsedSecs *uint64,
) BackupTask {
	taskModel := BackupTask{
		AccountID:              accountId,
		JobID:                  jobId,
		DestinationID:          destinationId,
		TaskStatus:             taskStatus,
		RetentionStrategy:      retentionStrategy,
		BackupSchedule:         backupSchedule,
		TimeoutSecs:            timeoutSecs,
		SuccessfulContainerIds: successfulContainerIds,
		FailedContainerIds:     failedContainerIds,
		ExecutionOutput:        executionOutput,
		StartedAt:              startedAt,
		FinishedAt:             finishedAt,
		ElapsedSecs:            elapsedSecs,
	}

	if id != 0 {
		taskModel.ID = id
	}

	return taskModel
}

func (model BackupTask) ToEntity() (taskEntity entity.BackupTask, err error) {
	taskId, err := valueObject.NewBackupTaskId(model.ID)
	if err != nil {
		return taskEntity, err
	}

	accountId, err := valueObject.NewAccountId(model.AccountID)
	if err != nil {
		return taskEntity, err
	}

	jobId, err := valueObject.NewBackupJobId(model.JobID)
	if err != nil {
		return taskEntity, err
	}

	destinationId, err := valueObject.NewBackupDestinationId(model.DestinationID)
	if err != nil {
		return taskEntity, err
	}

	taskStatus, err := valueObject.NewBackupTaskStatus(model.TaskStatus)
	if err != nil {
		return taskEntity, err
	}

	retentionStrategy, err := valueObject.NewBackupRetentionStrategy(model.RetentionStrategy)
	if err != nil {
		return taskEntity, err
	}

	backupSchedule, err := valueObject.NewCronSchedule(model.BackupSchedule)
	if err != nil {
		return taskEntity, err
	}

	successfulContainerIds := []valueObject.ContainerId{}
	for _, rawContainerId := range model.SuccessfulContainerIds {
		containerId, err := valueObject.NewContainerId(rawContainerId)
		if err != nil {
			slog.Debug(err.Error(), slog.String("containerId", rawContainerId))
			continue
		}
		successfulContainerIds = append(successfulContainerIds, containerId)
	}

	failedContainerIds := []valueObject.ContainerId{}
	for _, rawContainerId := range model.FailedContainerIds {
		containerId, err := valueObject.NewContainerId(rawContainerId)
		if err != nil {
			slog.Debug(err.Error(), slog.String("containerId", rawContainerId))
			continue
		}
		failedContainerIds = append(failedContainerIds, containerId)
	}

	var executionOutputPtr *valueObject.BackupTaskExecutionOutput
	if model.ExecutionOutput != nil {
		executionOutput, err := valueObject.NewBackupTaskExecutionOutput(*model.ExecutionOutput)
		if err != nil {
			return taskEntity, err
		}
		executionOutputPtr = &executionOutput
	}

	var startedAtPtr *valueObject.UnixTime
	if model.StartedAt != nil {
		startedAt, err := valueObject.NewUnixTime(*model.StartedAt)
		if err != nil {
			return taskEntity, err
		}
		startedAtPtr = &startedAt
	}

	var finishedAtPtr *valueObject.UnixTime
	if model.FinishedAt != nil {
		finishedAt, err := valueObject.NewUnixTime(*model.FinishedAt)
		if err != nil {
			return taskEntity, err
		}
		finishedAtPtr = &finishedAt
	}

	var elapsedSecsPtr *uint64
	if model.ElapsedSecs != nil {
		elapsedSecsPtr = model.ElapsedSecs
	}

	return entity.NewBackupTask(
		taskId, accountId, jobId, destinationId, taskStatus, retentionStrategy,
		backupSchedule, model.TimeoutSecs, successfulContainerIds, failedContainerIds,
		executionOutputPtr, startedAtPtr, finishedAtPtr, elapsedSecsPtr,
		valueObject.NewUnixTimeWithGoTime(model.CreatedAt),
		valueObject.NewUnixTimeWithGoTime(model.UpdatedAt),
	), nil
}
