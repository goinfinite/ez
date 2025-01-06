package repository

import (
	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/valueObject"
)

type BackupCmdRepo interface {
	CreateDestination(dto.CreateBackupDestination) (valueObject.BackupDestinationId, error)
	UpdateDestination(dto.UpdateBackupDestination) error
	DeleteDestination(dto.DeleteBackupDestination) error
	CreateJob(dto.CreateBackupJob) (valueObject.BackupJobId, error)
	UpdateJob(dto.UpdateBackupJob) error
}
