package backupInfra

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/entity"
	"github.com/goinfinite/ez/src/infra/db"
	dbHelper "github.com/goinfinite/ez/src/infra/db/helper"
	dbModel "github.com/goinfinite/ez/src/infra/db/model"
)

type BackupQueryRepo struct {
	persistentDbSvc *db.PersistentDatabaseService
}

func NewBackupQueryRepo(
	persistentDbSvc *db.PersistentDatabaseService,
) *BackupQueryRepo {
	return &BackupQueryRepo{persistentDbSvc: persistentDbSvc}
}

func (repo *BackupQueryRepo) ReadDestination(
	readDto dto.ReadBackupDestinationsRequest,
) (responseDto dto.ReadBackupDestinationsResponse, err error) {
	backupDestinationModel := dbModel.BackupDestination{}
	if readDto.DestinationId != nil {
		backupDestinationModel.ID = readDto.DestinationId.Uint64()
	}
	if readDto.DestinationName != nil {
		backupDestinationModel.Name = readDto.DestinationName.String()
	}
	if readDto.DestinationType != nil {
		backupDestinationModel.Type = readDto.DestinationType.String()
	}
	if readDto.ObjectStorageProvider != nil {
		objectStorageProviderStr := readDto.ObjectStorageProvider.String()
		backupDestinationModel.ObjectStorageProvider = &objectStorageProviderStr
	}
	if readDto.RemoteHostType != nil {
		remoteHostTypeStr := readDto.RemoteHostType.String()
		backupDestinationModel.RemoteHostType = &remoteHostTypeStr
	}
	if readDto.RemoteHostname != nil {
		remoteHostnameStr := readDto.RemoteHostname.String()
		backupDestinationModel.RemoteHostname = &remoteHostnameStr
	}

	dbQuery := repo.persistentDbSvc.Handler.Model(backupDestinationModel).Where(&backupDestinationModel)
	if readDto.CreatedBeforeAt != nil {
		dbQuery = dbQuery.Where("created_at < ?", readDto.CreatedBeforeAt.GetAsGoTime())
	}
	if readDto.CreatedAfterAt != nil {
		dbQuery = dbQuery.Where("created_at > ?", readDto.CreatedAfterAt.GetAsGoTime())
	}

	paginatedDbQuery, responsePagination, err := dbHelper.PaginationQueryBuilder(
		dbQuery, readDto.Pagination,
	)
	if err != nil {
		return responseDto, errors.New("PaginationQueryBuilderError: " + err.Error())
	}

	backupDestinationModels := []dbModel.BackupDestination{}
	err = paginatedDbQuery.Find(&backupDestinationModels).Error
	if err != nil {
		return responseDto, errors.New("FindBackupDestinationsError: " + err.Error())
	}

	backupDestinationEntities := []entity.IBackupDestination{}
	for _, backupDestinationModel := range backupDestinationModels {
		backupDestinationEntity, err := backupDestinationModel.ToEntity()
		if err != nil {
			slog.Debug(
				"ModelToEntityError",
				slog.Uint64("id", backupDestinationModel.ID),
				slog.Any("error", err),
			)
			continue
		}
		backupDestinationEntities = append(backupDestinationEntities, backupDestinationEntity)
	}

	return dto.ReadBackupDestinationsResponse{
		Pagination:   responsePagination,
		Destinations: backupDestinationEntities,
	}, nil
}

func (repo *BackupQueryRepo) ReadJob(
	readDto dto.ReadBackupJobsRequest,
) (responseDto dto.ReadBackupJobsResponse, err error) {
	backupJobModel := dbModel.BackupJob{}
	if readDto.JobId != nil {
		backupJobModel.ID = readDto.JobId.Uint64()
	}
	if readDto.JobStatus != nil {
		backupJobModel.JobStatus = *readDto.JobStatus
	}
	if readDto.AccountId != nil {
		backupJobModel.AccountID = readDto.AccountId.Uint64()
	}
	if readDto.DestinationId != nil {
		backupJobModel.DestinationIds = []uint64{readDto.DestinationId.Uint64()}
	}
	if readDto.RetentionStrategy != nil {
		retentionStrategyStr := readDto.RetentionStrategy.String()
		backupJobModel.RetentionStrategy = retentionStrategyStr
	}
	if readDto.ArchiveCompressionFormat != nil {
		archiveCompressionFormatStr := readDto.ArchiveCompressionFormat.String()
		backupJobModel.ArchiveCompressionFormat = archiveCompressionFormatStr
	}
	if readDto.LastRunStatus != nil {
		lastRunStatusStr := readDto.LastRunStatus.String()
		backupJobModel.LastRunStatus = &lastRunStatusStr
	}

	dbQuery := repo.persistentDbSvc.Handler.Model(backupJobModel).Where(&backupJobModel)
	if readDto.LastRunBeforeAt != nil {
		dbQuery = dbQuery.Where("last_run_at < ?", readDto.LastRunBeforeAt.GetAsGoTime())
	}
	if readDto.LastRunAfterAt != nil {
		dbQuery = dbQuery.Where("last_run_at > ?", readDto.LastRunAfterAt.GetAsGoTime())
	}
	if readDto.NextRunBeforeAt != nil {
		dbQuery = dbQuery.Where("next_run_at < ?", readDto.NextRunBeforeAt.GetAsGoTime())
	}
	if readDto.NextRunAfterAt != nil {
		dbQuery = dbQuery.Where("next_run_at > ?", readDto.NextRunAfterAt.GetAsGoTime())
	}
	if readDto.CreatedBeforeAt != nil {
		dbQuery = dbQuery.Where("created_at < ?", readDto.CreatedBeforeAt.GetAsGoTime())
	}
	if readDto.CreatedAfterAt != nil {
		dbQuery = dbQuery.Where("created_at > ?", readDto.CreatedAfterAt.GetAsGoTime())
	}

	paginatedDbQuery, responsePagination, err := dbHelper.PaginationQueryBuilder(
		dbQuery, readDto.Pagination,
	)
	if err != nil {
		return responseDto, errors.New("PaginationQueryBuilderError: " + err.Error())
	}

	backupJobModels := []dbModel.BackupJob{}
	err = paginatedDbQuery.Find(&backupJobModels).Error
	if err != nil {
		return responseDto, errors.New("FindBackupJobsError: " + err.Error())
	}

	backupJobEntities := []entity.BackupJob{}
	for _, backupJobModel := range backupJobModels {
		backupJobEntity, err := backupJobModel.ToEntity()
		if err != nil {
			slog.Debug(
				"ModelToEntityError",
				slog.Uint64("id", backupJobModel.ID),
				slog.Any("error", err),
			)
			continue
		}
		backupJobEntities = append(backupJobEntities, backupJobEntity)
	}

	return dto.ReadBackupJobsResponse{
		Pagination: responsePagination,
		Jobs:       backupJobEntities,
	}, nil
}

func (repo *BackupQueryRepo) ReadTask(
	readDto dto.ReadBackupTasksRequest,
) (responseDto dto.ReadBackupTasksResponse, err error) {
	backupTaskModel := dbModel.BackupTask{}
	if readDto.TaskId != nil {
		backupTaskModel.ID = readDto.TaskId.Uint64()
	}
	if readDto.AccountId != nil {
		backupTaskModel.AccountID = readDto.AccountId.Uint64()
	}
	if readDto.JobId != nil {
		backupTaskModel.JobID = readDto.JobId.Uint64()
	}
	if readDto.DestinationId != nil {
		backupTaskModel.DestinationID = readDto.DestinationId.Uint64()
	}
	if readDto.TaskStatus != nil {
		taskStatusStr := readDto.TaskStatus.String()
		backupTaskModel.TaskStatus = taskStatusStr
	}
	if readDto.RetentionStrategy != nil {
		retentionStrategyStr := readDto.RetentionStrategy.String()
		backupTaskModel.RetentionStrategy = retentionStrategyStr
	}
	if readDto.ContainerId != nil {
		backupTaskModel.ContainerIds = []string{readDto.ContainerId.String()}
	}

	dbQuery := repo.persistentDbSvc.Handler.Model(backupTaskModel).Where(&backupTaskModel)
	if readDto.StartedBeforeAt != nil {
		dbQuery = dbQuery.Where("started_at < ?", readDto.StartedBeforeAt.GetAsGoTime())
	}
	if readDto.StartedAfterAt != nil {
		dbQuery = dbQuery.Where("started_at > ?", readDto.StartedAfterAt.GetAsGoTime())
	}
	if readDto.FinishedBeforeAt != nil {
		dbQuery = dbQuery.Where("finished_at < ?", readDto.FinishedBeforeAt.GetAsGoTime())
	}
	if readDto.FinishedAfterAt != nil {
		dbQuery = dbQuery.Where("finished_at > ?", readDto.FinishedAfterAt.GetAsGoTime())
	}
	if readDto.CreatedBeforeAt != nil {
		dbQuery = dbQuery.Where("created_at < ?", readDto.CreatedBeforeAt.GetAsGoTime())
	}
	if readDto.CreatedAfterAt != nil {
		dbQuery = dbQuery.Where("created_at > ?", readDto.CreatedAfterAt.GetAsGoTime())
	}

	paginatedDbQuery, responsePagination, err := dbHelper.PaginationQueryBuilder(
		dbQuery, readDto.Pagination,
	)
	if err != nil {
		return responseDto, errors.New("PaginationQueryBuilderError: " + err.Error())
	}

	backupTaskModels := []dbModel.BackupTask{}
	err = paginatedDbQuery.Find(&backupTaskModels).Error
	if err != nil {
		return responseDto, errors.New("FindBackupTasksError: " + err.Error())
	}

	backupTaskEntities := []entity.BackupTask{}
	for _, backupTaskModel := range backupTaskModels {
		backupTaskEntity, err := backupTaskModel.ToEntity()
		if err != nil {
			slog.Debug(
				"ModelToEntityError",
				slog.Uint64("id", backupTaskModel.ID),
				slog.Any("error", err),
			)
			continue
		}
		backupTaskEntities = append(backupTaskEntities, backupTaskEntity)
	}

	return dto.ReadBackupTasksResponse{
		Pagination: responsePagination,
		Tasks:      backupTaskEntities,
	}, nil
}
