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

func (uc *DeleteBackupTask) updateBackupDestinationStats(
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
		newDestinationTasksCount = iDestinationEntity.ReadTasksCount() - 1
	}

	newDestinationTotalSpaceUsageBytes := iDestinationEntity.ReadTotalSpaceUsageBytes()
	if taskSizeBytes != nil {
		newDestinationTotalSpaceUsageBytes -= *taskSizeBytes
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

func (uc *DeleteBackupTask) updateBackupJobStats(
	jobId valueObject.BackupJobId,
	taskSizeBytes *valueObject.Byte,
) error {
	jobEntity, err := uc.backupQueryRepo.ReadFirstJob(
		dto.ReadBackupJobsRequest{JobId: &jobId},
	)
	if err != nil {
		return errors.New("ReadBackupJobInfraError: " + err.Error())
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
		return errors.New("UpdateBackupJobInfraError: " + err.Error())
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

	err = uc.updateBackupJobStats(taskEntity.JobId, taskEntity.SizeBytes)
	if err != nil {
		slog.Error(
			"UpdateBackupJobStatsError",
			slog.String("error", err.Error()),
			slog.String("jobId", taskEntity.JobId.String()),
		)
	}

	err = uc.updateBackupDestinationStats(
		taskEntity.DestinationId, taskEntity.SizeBytes,
	)
	if err != nil {
		slog.Error(
			"UpdateBackupDestinationStatsError",
			slog.String("error", err.Error()),
			slog.String("destinationId", taskEntity.DestinationId.String()),
		)
	}

	return nil
}
