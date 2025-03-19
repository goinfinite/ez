package useCase

import (
	"errors"
	"time"

	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/repository"
	"github.com/goinfinite/ez/src/domain/valueObject"
)

func deleteAncientBackupTasks(
	backupCmdRepo repository.BackupCmdRepo,
	backupQueryRepo repository.BackupQueryRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	jobId valueObject.BackupJobId,
	lastValidDay valueObject.UnixTime,
) error {
	readAncientTaskResponse, err := backupQueryRepo.ReadTask(dto.ReadBackupTasksRequest{
		JobId:           &jobId,
		CreatedBeforeAt: &lastValidDay,
	})
	if err != nil {
		return errors.New("ReadAncientTaskInfraError: " + err.Error())
	}

	for _, taskEntity := range readAncientTaskResponse.Tasks {
		err = DeleteBackupTask(
			backupQueryRepo, backupCmdRepo, activityRecordCmdRepo,
			dto.DeleteBackupTask{
				TaskId:             taskEntity.TaskId,
				ShouldDiscardFiles: true,
				OperatorAccountId:  valueObject.SystemAccountId,
				OperatorIpAddress:  valueObject.SystemIpAddress,
			},
		)
		if err != nil {
			continue
		}
	}

	return nil
}

func BackupJobHousekeeper(
	backupQueryRepo repository.BackupQueryRepo,
	backupCmdRepo repository.BackupCmdRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	jobId valueObject.BackupJobId,
) error {
	jobEntity, err := backupQueryRepo.ReadFirstJob(dto.ReadBackupJobsRequest{
		JobId: &jobId,
	})
	if err != nil {
		return errors.New("BackupJobNotFound")
	}

	lastValidDay := valueObject.NewUnixTimeWithGoTime(
		time.Now().Add(-time.Duration(*jobEntity.MaxTaskRetentionDays) * time.Hour * 24),
	)

	if jobEntity.MaxTaskRetentionDays != nil {
		err = deleteAncientBackupTasks(
			backupCmdRepo, backupQueryRepo, activityRecordCmdRepo, jobEntity.JobId, lastValidDay,
		)
		if err != nil {
			return errors.New("DeleteAncientTaskError: " + err.Error())
		}
	}

	if jobEntity.MaxTaskRetentionCount == nil {
		return nil
	}

	jobEntity, err = backupQueryRepo.ReadFirstJob(dto.ReadBackupJobsRequest{
		JobId: &jobId,
	})
	if err != nil {
		return errors.New("ReloadBackupJobEntityError: " + err.Error())
	}

	if jobEntity.TasksCount > *jobEntity.MaxTaskRetentionCount {
		firstTaskEntity, err := backupQueryRepo.ReadFirstTask(
			dto.ReadBackupTasksRequest{JobId: &jobEntity.JobId},
		)
		if err != nil {
			return errors.New("BackupTaskNotFound: " + err.Error())
		}

		err = DeleteBackupTask(
			backupQueryRepo, backupCmdRepo, activityRecordCmdRepo,
			dto.DeleteBackupTask{
				TaskId:             firstTaskEntity.TaskId,
				ShouldDiscardFiles: true,
				OperatorAccountId:  valueObject.SystemAccountId,
				OperatorIpAddress:  valueObject.SystemIpAddress,
			},
		)
		if err != nil {
			return errors.New("DeleteOldestTaskInfraError: " + err.Error())
		}
	}

	return nil
}
