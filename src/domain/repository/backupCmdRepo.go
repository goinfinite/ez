package repository

import (
	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/valueObject"
)

type BackupCmdRepo interface {
	CreateDestination(dto.CreateBackupDestinationRequest) (dto.CreateBackupDestinationResponse, error)
	UpdateDestination(dto.UpdateBackupDestination) error
	DeleteDestination(dto.DeleteBackupDestination) error

	CreateJob(dto.CreateBackupJob) (valueObject.BackupJobId, error)
	UpdateJob(dto.UpdateBackupJob) error
	DeleteJob(dto.DeleteBackupJob) error
	RunJob(dto.RunBackupJob) error

	UpdateTask(dto.UpdateBackupTask) error
	DeleteTask(dto.DeleteBackupTask) error
	RestoreTask(dto.RestoreBackupTaskRequest) (dto.RestoreBackupTaskResponse, error)

	CreateTaskArchive(dto.CreateBackupTaskArchive) (valueObject.BackupTaskArchiveId, error)
	DeleteTaskArchive(dto.DeleteBackupTaskArchive) error
}
