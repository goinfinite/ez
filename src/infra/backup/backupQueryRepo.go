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
	backupDestinationEntities := []entity.IBackupDestination{}

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
	backupJobEntities := []entity.BackupJob{}

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
	if readDto.BackupType != nil {
		backupTypeStr := readDto.BackupType.String()
		backupJobModel.BackupType = backupTypeStr
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
