package repository

import (
	"github.com/goinfinite/ez/src/domain/dto"
)

type BackupQueryRepo interface {
	ReadDestination(dto.ReadBackupDestinationsRequest) (dto.ReadBackupDestinationsResponse, error)
	ReadJob(dto.ReadBackupJobsRequest) (dto.ReadBackupJobsResponse, error)
	ReadTask(dto.ReadBackupTasksRequest) (dto.ReadBackupTasksResponse, error)
}
