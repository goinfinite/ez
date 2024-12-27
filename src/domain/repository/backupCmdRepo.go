package repository

import (
	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/valueObject"
)

type BackupCmdRepo interface {
	CreateDestination(dto.CreateBackupDestination) (valueObject.BackupDestinationId, error)
}
