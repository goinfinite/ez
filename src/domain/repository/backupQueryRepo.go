package repository

import (
	"github.com/goinfinite/ez/src/domain/dto"
)

type BackupQueryRepo interface {
	ReadDestination(
		requestDto dto.ReadBackupDestinationsRequest,
		withSecrets bool,
	) (dto.ReadBackupDestinationsResponse, error)
	ReadJob(dto.ReadBackupJobsRequest) (dto.ReadBackupJobsResponse, error)
	ReadTask(dto.ReadBackupTasksRequest) (dto.ReadBackupTasksResponse, error)
}
