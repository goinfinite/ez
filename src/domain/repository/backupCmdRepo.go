package repository

import (
	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/valueObject"
)

type BackupCmdRepo interface {
	CreateDestination(dto.CreateBackupDestination) (valueObject.BackupDestinationId, error)
	CreateJob(dto.CreateBackupJob) (valueObject.BackupJobId, error)
	UpdateDestination(dto.UpdateBackupDestination) error
}
