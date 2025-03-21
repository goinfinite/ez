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

func (uc *RunBackupJob) updateBackupDestinationStats(
	destinationId valueObject.BackupDestinationId,
	taskSizeBytes *valueObject.Byte,
) error {
	iDestinationEntity, err := uc.backupQueryRepo.ReadFirstDestination(
		dto.ReadBackupDestinationsRequest{DestinationId: &destinationId}, false,
	)
	if err != nil {
		return errors.New("ReadBackupDestinationInfraError: " + err.Error())
	}

	newDestinationTasksCount := uint16(0)
	if iDestinationEntity.ReadTasksCount() > 0 {
		newDestinationTasksCount = iDestinationEntity.ReadTasksCount() + 1
	}

	newDestinationTotalSpaceUsageBytes := iDestinationEntity.ReadTotalSpaceUsageBytes()
	if taskSizeBytes != nil {
		newDestinationTotalSpaceUsageBytes += *taskSizeBytes
	}
	if newDestinationTotalSpaceUsageBytes < 0 {
		newDestinationTotalSpaceUsageBytes = valueObject.Byte(0)
	}

	err = uc.backupCmdRepo.UpdateDestination(dto.UpdateBackupDestination{
		DestinationId:        destinationId,
		TasksCount:           &newDestinationTasksCount,
		TotalSpaceUsageBytes: &newDestinationTotalSpaceUsageBytes,
	})
	if err != nil {
		return errors.New("UpdateBackupDestinationInfraError: " + err.Error())
	}

	return nil
}

func (uc *RunBackupJob) updateBackupJobStats(
	jobId valueObject.BackupJobId,
	taskIds []valueObject.BackupTaskId,
) error {
	jobEntity, err := uc.backupQueryRepo.ReadFirstJob(
		dto.ReadBackupJobsRequest{JobId: &jobId},
	)
	if err != nil {
		return errors.New("ReadBackupJobInfraError")
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

		err = uc.updateBackupDestinationStats(
			taskEntity.DestinationId, taskEntity.SizeBytes,
		)
		if err != nil {
			slog.Error(
				"UpdateBackupDestinationStatsError",
				slog.String("error", err.Error()),
				slog.String("taskId", taskId.String()),
				slog.String("destinationId", taskEntity.DestinationId.String()),
			)
		}

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
		JobId:                jobId,
		TasksCount:           &newTasksCount,
		TotalSpaceUsageBytes: &totalUsageBytes,
		LastRunAt:            &lastRunAt,
		LastRunStatus:        &lastRunStatus,
		NextRunAt:            nextRunAtPtr,
	})
	if err != nil {
		return errors.New("UpdateBackupJobStatsInfraError: " + err.Error())
	}

	return nil
}

func (uc *RunBackupJob) Execute(runDto dto.RunBackupJob) error {
	_, err := uc.backupQueryRepo.ReadFirstJob(
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

	err = uc.updateBackupJobStats(runDto.JobId, taskIds)
	if err != nil {
		slog.Error(
			"UpdateBackupJobStatsError",
			slog.String("error", err.Error()),
			slog.String("jobId", runDto.JobId.String()),
		)
		return nil
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
