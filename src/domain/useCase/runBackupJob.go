package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/repository"
	"github.com/goinfinite/ez/src/domain/valueObject"
)

type RunBackupJob struct {
	backupQueryRepo       repository.BackupQueryRepo
	backupCmdRepo         repository.BackupCmdRepo
	cronQueryRepo         repository.CronQueryRepo
	activityRecordCmdRepo repository.ActivityRecordCmdRepo
}

func NewRunBackupJob(
	backupQueryRepo repository.BackupQueryRepo,
	backupCmdRepo repository.BackupCmdRepo,
	cronQueryRepo repository.CronQueryRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
) *RunBackupJob {
	return &RunBackupJob{
		backupQueryRepo:       backupQueryRepo,
		backupCmdRepo:         backupCmdRepo,
		cronQueryRepo:         cronQueryRepo,
		activityRecordCmdRepo: activityRecordCmdRepo,
	}
}

func (uc *RunBackupJob) Execute(runDto dto.RunBackupJob) error {
	jobEntity, err := uc.backupQueryRepo.ReadFirstJob(
		dto.ReadBackupJobsRequest{JobId: &runDto.JobId},
	)
	if err != nil {
		return errors.New("BackupJobNotFound")
	}

	backupHousekeeper := NewBackupHousekeeper(
		uc.backupQueryRepo, uc.backupCmdRepo, uc.activityRecordCmdRepo,
	)
	err = backupHousekeeper.CleanJobTasks(runDto.JobId)
	if err != nil {
		slog.Error(
			"PreRunBackupJobHousekeeperError",
			slog.String("error", err.Error()),
			slog.String("jobId", runDto.JobId.String()),
		)
	}

	taskIds, err := uc.backupCmdRepo.RunJob(runDto)
	if err != nil {
		slog.Error(
			"RunBackupJobInfraError",
			slog.String("error", err.Error()),
			slog.String("jobId", runDto.JobId.String()),
		)
		return errors.New("RunBackupJobInfraError")
	}

	NewCreateSecurityActivityRecord(uc.activityRecordCmdRepo).RunBackupJob(runDto)

	jobEntity, err = uc.backupQueryRepo.ReadFirstJob(
		dto.ReadBackupJobsRequest{JobId: &runDto.JobId},
	)
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
		taskEntity, err := uc.backupQueryRepo.ReadFirstTask(
			dto.ReadBackupTasksRequest{TaskId: &taskId},
		)
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
	var nextRunAtPtr *valueObject.UnixTime
	nextRunAt, err := uc.cronQueryRepo.ReadNextRun(jobEntity.BackupSchedule)
	if err == nil {
		nextRunAtPtr = &nextRunAt
	}

	err = uc.backupCmdRepo.UpdateJob(dto.UpdateBackupJob{
		JobId:                runDto.JobId,
		TasksCount:           &newTasksCount,
		TotalSpaceUsageBytes: &totalUsageBytes,
		LastRunAt:            &lastRunAt,
		LastRunStatus:        &lastRunStatus,
		NextRunAt:            nextRunAtPtr,
	})
	if err != nil {
		slog.Error(
			"UpdateBackupJobStatsInfraError",
			slog.String("error", err.Error()),
			slog.String("jobId", runDto.JobId.String()),
		)
	}

	err = backupHousekeeper.CleanJobTasks(runDto.JobId)
	if err != nil {
		slog.Error(
			"PostRunBackupJobHousekeeperError",
			slog.String("error", err.Error()),
			slog.String("jobId", runDto.JobId.String()),
		)
	}

	return nil
}
