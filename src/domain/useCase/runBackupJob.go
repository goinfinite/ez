package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/repository"
	"github.com/goinfinite/ez/src/domain/valueObject"
)

func RunBackupJob(
	backupQueryRepo repository.BackupQueryRepo,
	backupCmdRepo repository.BackupCmdRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	runDto dto.RunBackupJob,
) error {
	jobEntity, err := backupQueryRepo.ReadFirstJob(dto.ReadBackupJobsRequest{
		JobId: &runDto.JobId,
	})
	if err != nil {
		return errors.New("BackupJobNotFound")
	}

	err = BackupJobHousekeeper(
		backupQueryRepo, backupCmdRepo, activityRecordCmdRepo, runDto.JobId,
	)
	if err != nil {
		slog.Error(
			"PreRunBackupJobHousekeeperError",
			slog.String("error", err.Error()),
			slog.String("jobId", runDto.JobId.String()),
		)
	}

	taskIds, err := backupCmdRepo.RunJob(runDto)
	if err != nil {
		slog.Error(
			"RunBackupJobInfraError",
			slog.String("error", err.Error()),
			slog.String("jobId", runDto.JobId.String()),
		)
		return errors.New("RunBackupJobInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).RunBackupJob(runDto)

	jobEntity, err = backupQueryRepo.ReadFirstJob(dto.ReadBackupJobsRequest{
		JobId: &runDto.JobId,
	})
	if err != nil {
		slog.Error(
			"ReloadBackupJobInfraError",
			slog.String("error", err.Error()),
			slog.String("jobId", runDto.JobId.String()),
		)
		return nil
	}

	totalUsageBytes := jobEntity.TotalSpaceUsageBytes
	lastRunStatus := valueObject.BackupTaskStatusCompleted
	for _, taskId := range taskIds {
		taskEntity, err := backupQueryRepo.ReadFirstTask(dto.ReadBackupTasksRequest{
			TaskId: &taskId,
		})
		if err != nil {
			slog.Error(
				"ReadBackupTaskInfraError",
				slog.String("error", err.Error()),
				slog.String("taskId", taskId.String()),
			)
			continue
		}

		lastRunStatus = taskEntity.TaskStatus

		if taskEntity.SizeBytes == nil {
			continue
		}
		totalUsageBytes += *taskEntity.SizeBytes
	}

	newTasksCount := jobEntity.TasksCount + uint16(len(taskIds))
	lastRunAt := valueObject.NewUnixTimeNow()

	err = backupCmdRepo.UpdateJob(dto.UpdateBackupJob{
		JobId:                runDto.JobId,
		TasksCount:           &newTasksCount,
		LastRunAt:            &lastRunAt,
		LastRunStatus:        &lastRunStatus,
		TotalSpaceUsageBytes: &totalUsageBytes,
	})
	if err != nil {
		slog.Error(
			"UpdateBackupJobStatsInfraError",
			slog.String("error", err.Error()),
			slog.String("jobId", runDto.JobId.String()),
		)
	}

	err = BackupJobHousekeeper(
		backupQueryRepo, backupCmdRepo, activityRecordCmdRepo, runDto.JobId,
	)
	if err != nil {
		slog.Error(
			"PostRunBackupJobHousekeeperError",
			slog.String("error", err.Error()),
			slog.String("jobId", runDto.JobId.String()),
		)
	}

	return nil
}
