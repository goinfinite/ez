package repository

import (
	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/entity"
)

type BackupQueryRepo interface {
	ReadDestination(
		requestDto dto.ReadBackupDestinationsRequest,
		withSecrets bool,
	) (dto.ReadBackupDestinationsResponse, error)
	ReadFirstDestination(
		requestDto dto.ReadBackupDestinationsRequest,
		withSecrets bool,
	) (entity.IBackupDestination, error)

	ReadJob(dto.ReadBackupJobsRequest) (dto.ReadBackupJobsResponse, error)
	ReadFirstJob(dto.ReadBackupJobsRequest) (entity.BackupJob, error)

	ReadTask(dto.ReadBackupTasksRequest) (dto.ReadBackupTasksResponse, error)
	ReadFirstTask(dto.ReadBackupTasksRequest) (entity.BackupTask, error)

	ReadTaskArchive(
		dto.ReadBackupTaskArchivesRequest,
	) (dto.ReadBackupTaskArchivesResponse, error)
}
