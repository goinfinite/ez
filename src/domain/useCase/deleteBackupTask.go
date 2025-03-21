package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/repository"
	"github.com/goinfinite/ez/src/domain/valueObject"
)

type DeleteBackupTask struct {
	backupQueryRepo       repository.BackupQueryRepo
	backupCmdRepo         repository.BackupCmdRepo
	activityRecordCmdRepo repository.ActivityRecordCmdRepo
}

func NewDeleteBackupTask(
	backupQueryRepo repository.BackupQueryRepo,
	backupCmdRepo repository.BackupCmdRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
) *DeleteBackupTask {
	return &DeleteBackupTask{
		backupQueryRepo:       backupQueryRepo,
		backupCmdRepo:         backupCmdRepo,
		activityRecordCmdRepo: activityRecordCmdRepo,
	}
}

func (uc *DeleteBackupTask) updateBackupJobStats(
	jobId valueObject.BackupJobId,
	taskSizeBytes *valueObject.Byte,
) error {
	jobEntity, err := uc.backupQueryRepo.ReadFirstJob(
		dto.ReadBackupJobsRequest{JobId: &jobId},
	)
	if err != nil {
		slog.Error(
			"BackupJobNotFound",
			slog.String("error", err.Error()),
			slog.String("jobId", jobId.String()),
		)
		return errors.New("BackupJobNotFound")
	}

	newJobTasksCount := uint16(0)
	if jobEntity.TasksCount > 0 {
		newJobTasksCount = jobEntity.TasksCount - 1
	}

	newJobTotalSpaceUsageBytes := jobEntity.TotalSpaceUsageBytes
	if taskSizeBytes != nil {
		newJobTotalSpaceUsageBytes -= *taskSizeBytes
	}
	if newJobTotalSpaceUsageBytes < 0 {
		newJobTotalSpaceUsageBytes = valueObject.Byte(0)
	}

	err = uc.backupCmdRepo.UpdateJob(dto.UpdateBackupJob{
		JobId:                jobId,
		TasksCount:           &newJobTasksCount,
		TotalSpaceUsageBytes: &newJobTotalSpaceUsageBytes,
	})
	if err != nil {
		slog.Error("UpdateBackupJobStatsInfraError", slog.String("error", err.Error()))
		return errors.New("UpdateBackupJobStatsInfraError")
	}

	return nil
}

func (uc *DeleteBackupTask) Execute(deleteDto dto.DeleteBackupTask) error {
	taskEntity, err := uc.backupQueryRepo.ReadFirstTask(
		dto.ReadBackupTasksRequest{TaskId: &deleteDto.TaskId},
	)
	if err != nil {
		return errors.New("BackupTaskNotFound")
	}

	err = uc.backupCmdRepo.DeleteTask(deleteDto)
	if err != nil {
		slog.Error(
			"DeleteBackupTaskInfraError",
			slog.String("error", err.Error()),
			slog.String("taskId", taskEntity.TaskId.String()),
		)
		return errors.New("DeleteBackupTaskError")
	}

	NewCreateSecurityActivityRecord(uc.activityRecordCmdRepo).
		DeleteBackupTask(deleteDto, taskEntity.AccountId)

	return uc.updateBackupJobStats(taskEntity.JobId, taskEntity.SizeBytes)
}
