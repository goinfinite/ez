package useCase

import (
	"errors"
	"log/slog"
	"time"

	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/repository"
	"github.com/goinfinite/ez/src/domain/valueObject"
)

type BackupHousekeeper struct {
	backupQueryRepo       repository.BackupQueryRepo
	backupCmdRepo         repository.BackupCmdRepo
	activityRecordCmdRepo repository.ActivityRecordCmdRepo
	deleteBackupTask      *DeleteBackupTask
}

func NewBackupHousekeeper(
	backupQueryRepo repository.BackupQueryRepo,
	backupCmdRepo repository.BackupCmdRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
) *BackupHousekeeper {
	return &BackupHousekeeper{
		backupQueryRepo:       backupQueryRepo,
		backupCmdRepo:         backupCmdRepo,
		activityRecordCmdRepo: activityRecordCmdRepo,
		deleteBackupTask: NewDeleteBackupTask(
			backupQueryRepo, backupCmdRepo, activityRecordCmdRepo,
		),
	}
}

func (uc *BackupHousekeeper) deleteAncientBackupTasks(
	jobId valueObject.BackupJobId,
	lastValidDay valueObject.UnixTime,
) error {
	readAncientTaskResponse, err := uc.backupQueryRepo.ReadTask(dto.ReadBackupTasksRequest{
		JobId:           &jobId,
		CreatedBeforeAt: &lastValidDay,
	})
	if err != nil {
		return errors.New("ReadAncientTaskInfraError: " + err.Error())
	}

	for _, taskEntity := range readAncientTaskResponse.Tasks {
		err = uc.deleteBackupTask.Execute(dto.DeleteBackupTask{
			TaskId:             taskEntity.TaskId,
			ShouldDiscardFiles: true,
			OperatorAccountId:  valueObject.SystemAccountId,
			OperatorIpAddress:  valueObject.SystemIpAddress,
		})
		if err != nil {
			continue
		}

		slog.Debug(
			"AncientBackupTaskDeleted",
			slog.String("taskId", taskEntity.TaskId.String()),
			slog.String("taskCreatedAt", taskEntity.CreatedAt.ReadAsRfcDate()),
		)
	}

	return nil
}

func (uc *BackupHousekeeper) CleanJobTasks(jobId valueObject.BackupJobId) error {
	jobEntity, err := uc.backupQueryRepo.ReadFirstJob(dto.ReadBackupJobsRequest{
		JobId: &jobId,
	})
	if err != nil {
		return errors.New("BackupJobNotFound")
	}

	lastValidDay := valueObject.NewUnixTimeWithGoTime(
		time.Now().Add(-time.Duration(*jobEntity.MaxTaskRetentionDays) * time.Hour * 24),
	)

	if jobEntity.MaxTaskRetentionDays != nil {
		err = uc.deleteAncientBackupTasks(jobEntity.JobId, lastValidDay)
		if err != nil {
			return errors.New("DeleteAncientTaskError: " + err.Error())
		}
	}

	if jobEntity.MaxTaskRetentionCount == nil {
		return nil
	}

	jobEntity, err = uc.backupQueryRepo.ReadFirstJob(
		dto.ReadBackupJobsRequest{JobId: &jobId},
	)
	if err != nil {
		return errors.New("ReloadBackupJobEntityError: " + err.Error())
	}

	if jobEntity.TasksCount > *jobEntity.MaxTaskRetentionCount {
		firstTaskEntity, err := uc.backupQueryRepo.ReadFirstTask(
			dto.ReadBackupTasksRequest{JobId: &jobEntity.JobId},
		)
		if err != nil {
			return errors.New("BackupTaskNotFound: " + err.Error())
		}

		err = uc.deleteBackupTask.Execute(dto.DeleteBackupTask{
			TaskId:             firstTaskEntity.TaskId,
			ShouldDiscardFiles: true,
			OperatorAccountId:  valueObject.SystemAccountId,
			OperatorIpAddress:  valueObject.SystemIpAddress,
		})
		if err != nil {
			return errors.New("DeleteOldestTaskInfraError: " + err.Error())
		}

		slog.Debug(
			"OldestBackupTaskDeleted",
			slog.String("taskId", firstTaskEntity.TaskId.String()),
			slog.String("taskCreatedAt", firstTaskEntity.CreatedAt.ReadAsRfcDate()),
		)
	}

	return nil
}
