package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/repository"
)

func DeleteBackupTask(
	backupQueryRepo repository.BackupQueryRepo,
	backupCmdRepo repository.BackupCmdRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	deleteDto dto.DeleteBackupTask,
) error {
	taskEntity, err := backupQueryRepo.ReadFirstTask(
		dto.ReadBackupTasksRequest{TaskId: &deleteDto.TaskId},
	)
	if err != nil {
		return errors.New("BackupTaskNotFound")
	}

	err = backupCmdRepo.DeleteTask(deleteDto)
	if err != nil {
		slog.Error(
			"DeleteBackupTaskInfraError",
			slog.String("error", err.Error()),
			slog.String("taskId", taskEntity.TaskId.String()),
		)
		return errors.New("DeleteBackupTaskError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).
		DeleteBackupTask(deleteDto, taskEntity.AccountId)

	jobEntity, err := backupQueryRepo.ReadFirstJob(
		dto.ReadBackupJobsRequest{JobId: &taskEntity.JobId},
	)
	if err != nil {
		slog.Error(
			"BackupJobNotFound",
			slog.String("error", err.Error()),
			slog.String("jobId", taskEntity.JobId.String()),
		)
		return nil
	}

	newJobTasksCount := jobEntity.TasksCount - 1
	newJobTotalSpaceUsageBytes := jobEntity.TotalSpaceUsageBytes
	if taskEntity.SizeBytes != nil {
		newJobTotalSpaceUsageBytes -= *taskEntity.SizeBytes
	}

	err = backupCmdRepo.UpdateJob(dto.UpdateBackupJob{
		JobId:                taskEntity.JobId,
		TasksCount:           &newJobTasksCount,
		TotalSpaceUsageBytes: &newJobTotalSpaceUsageBytes,
	})
	if err != nil {
		slog.Error("UpdateBackupJobStatsInfraError", slog.String("error", err.Error()))
		return errors.New("UpdateBackupJobStatsInfraError")
	}

	return nil
}
