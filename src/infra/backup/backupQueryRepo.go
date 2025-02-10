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
	requestDto dto.ReadBackupDestinationsRequest,
	withSecrets bool,
) (responseDto dto.ReadBackupDestinationsResponse, err error) {
	backupDestinationModel := dbModel.BackupDestination{}
	if requestDto.DestinationId != nil {
		backupDestinationModel.ID = requestDto.DestinationId.Uint64()
	}
	if requestDto.DestinationName != nil {
		backupDestinationModel.Name = requestDto.DestinationName.String()
	}
	if requestDto.DestinationType != nil {
		backupDestinationModel.Type = requestDto.DestinationType.String()
	}
	if requestDto.ObjectStorageProvider != nil {
		objectStorageProviderStr := requestDto.ObjectStorageProvider.String()
		backupDestinationModel.ObjectStorageProvider = &objectStorageProviderStr
	}
	if requestDto.RemoteHostType != nil {
		remoteHostTypeStr := requestDto.RemoteHostType.String()
		backupDestinationModel.RemoteHostType = &remoteHostTypeStr
	}
	if requestDto.RemoteHostname != nil {
		remoteHostnameStr := requestDto.RemoteHostname.String()
		backupDestinationModel.RemoteHostname = &remoteHostnameStr
	}

	dbQuery := repo.persistentDbSvc.Handler.Model(backupDestinationModel).Where(&backupDestinationModel)
	if requestDto.CreatedBeforeAt != nil {
		dbQuery = dbQuery.Where("created_at < ?", requestDto.CreatedBeforeAt.GetAsGoTime())
	}
	if requestDto.CreatedAfterAt != nil {
		dbQuery = dbQuery.Where("created_at > ?", requestDto.CreatedAfterAt.GetAsGoTime())
	}

	paginatedDbQuery, responsePagination, err := dbHelper.PaginationQueryBuilder(
		dbQuery, requestDto.Pagination,
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
		backupDestinationEntity, err := backupDestinationModel.ToEntity(withSecrets)
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

func (repo *BackupQueryRepo) ReadFirstDestination(
	requestDto dto.ReadBackupDestinationsRequest,
	withSecrets bool,
) (destinationEntity entity.IBackupDestination, err error) {
	requestDto.Pagination = dto.Pagination{
		PageNumber:   0,
		ItemsPerPage: 1,
	}

	responseDto, err := repo.ReadDestination(requestDto, withSecrets)
	if err != nil {
		return destinationEntity, err
	}

	if len(responseDto.Destinations) == 0 {
		return destinationEntity, errors.New("BackupDestinationNotFound")
	}

	return responseDto.Destinations[0], nil
}

func (repo *BackupQueryRepo) ReadJob(
	requestDto dto.ReadBackupJobsRequest,
) (responseDto dto.ReadBackupJobsResponse, err error) {
	backupJobModel := dbModel.BackupJob{}
	if requestDto.JobId != nil {
		backupJobModel.ID = requestDto.JobId.Uint64()
	}
	if requestDto.JobStatus != nil {
		backupJobModel.JobStatus = *requestDto.JobStatus
	}
	if requestDto.AccountId != nil {
		backupJobModel.AccountID = requestDto.AccountId.Uint64()
	}
	if requestDto.DestinationId != nil {
		backupJobModel.DestinationIds = []uint64{requestDto.DestinationId.Uint64()}
	}
	if requestDto.RetentionStrategy != nil {
		retentionStrategyStr := requestDto.RetentionStrategy.String()
		backupJobModel.RetentionStrategy = retentionStrategyStr
	}
	if requestDto.ArchiveCompressionFormat != nil {
		archiveCompressionFormatStr := requestDto.ArchiveCompressionFormat.String()
		backupJobModel.ArchiveCompressionFormat = archiveCompressionFormatStr
	}
	if requestDto.LastRunStatus != nil {
		lastRunStatusStr := requestDto.LastRunStatus.String()
		backupJobModel.LastRunStatus = &lastRunStatusStr
	}

	dbQuery := repo.persistentDbSvc.Handler.Model(backupJobModel).Where(&backupJobModel)
	if requestDto.LastRunBeforeAt != nil {
		dbQuery = dbQuery.Where("last_run_at < ?", requestDto.LastRunBeforeAt.GetAsGoTime())
	}
	if requestDto.LastRunAfterAt != nil {
		dbQuery = dbQuery.Where("last_run_at > ?", requestDto.LastRunAfterAt.GetAsGoTime())
	}
	if requestDto.NextRunBeforeAt != nil {
		dbQuery = dbQuery.Where("next_run_at < ?", requestDto.NextRunBeforeAt.GetAsGoTime())
	}
	if requestDto.NextRunAfterAt != nil {
		dbQuery = dbQuery.Where("next_run_at > ?", requestDto.NextRunAfterAt.GetAsGoTime())
	}
	if requestDto.CreatedBeforeAt != nil {
		dbQuery = dbQuery.Where("created_at < ?", requestDto.CreatedBeforeAt.GetAsGoTime())
	}
	if requestDto.CreatedAfterAt != nil {
		dbQuery = dbQuery.Where("created_at > ?", requestDto.CreatedAfterAt.GetAsGoTime())
	}

	paginatedDbQuery, responsePagination, err := dbHelper.PaginationQueryBuilder(
		dbQuery, requestDto.Pagination,
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

func (repo *BackupQueryRepo) ReadFirstJob(
	requestDto dto.ReadBackupJobsRequest,
) (jobEntity entity.BackupJob, err error) {
	requestDto.Pagination = dto.Pagination{
		PageNumber:   0,
		ItemsPerPage: 1,
	}

	responseDto, err := repo.ReadJob(requestDto)
	if err != nil {
		return jobEntity, err
	}

	if len(responseDto.Jobs) == 0 {
		return jobEntity, errors.New("BackupJobNotFound")
	}

	return responseDto.Jobs[0], nil
}

func (repo *BackupQueryRepo) ReadTask(
	requestDto dto.ReadBackupTasksRequest,
) (responseDto dto.ReadBackupTasksResponse, err error) {
	backupTaskModel := dbModel.BackupTask{}
	if requestDto.TaskId != nil {
		backupTaskModel.ID = requestDto.TaskId.Uint64()
	}
	if requestDto.AccountId != nil {
		backupTaskModel.AccountID = requestDto.AccountId.Uint64()
	}
	if requestDto.JobId != nil {
		backupTaskModel.JobID = requestDto.JobId.Uint64()
	}
	if requestDto.DestinationId != nil {
		backupTaskModel.DestinationID = requestDto.DestinationId.Uint64()
	}
	if requestDto.TaskStatus != nil {
		taskStatusStr := requestDto.TaskStatus.String()
		backupTaskModel.TaskStatus = taskStatusStr
	}
	if requestDto.RetentionStrategy != nil {
		retentionStrategyStr := requestDto.RetentionStrategy.String()
		backupTaskModel.RetentionStrategy = retentionStrategyStr
	}

	dbQuery := repo.persistentDbSvc.Handler.Model(backupTaskModel).Where(&backupTaskModel)
	if requestDto.ContainerId != nil {
		containerIdLikeStr := "%" + requestDto.ContainerId.String() + "%"
		dbQuery = dbQuery.Where(
			"successful_container_ids LIKE ? OR failed_container_ids LIKE ?",
			containerIdLikeStr, containerIdLikeStr,
		)
	}
	if requestDto.StartedBeforeAt != nil {
		dbQuery = dbQuery.Where("started_at < ?", requestDto.StartedBeforeAt.GetAsGoTime())
	}
	if requestDto.StartedAfterAt != nil {
		dbQuery = dbQuery.Where("started_at > ?", requestDto.StartedAfterAt.GetAsGoTime())
	}
	if requestDto.FinishedBeforeAt != nil {
		dbQuery = dbQuery.Where("finished_at < ?", requestDto.FinishedBeforeAt.GetAsGoTime())
	}
	if requestDto.FinishedAfterAt != nil {
		dbQuery = dbQuery.Where("finished_at > ?", requestDto.FinishedAfterAt.GetAsGoTime())
	}
	if requestDto.CreatedBeforeAt != nil {
		dbQuery = dbQuery.Where("created_at < ?", requestDto.CreatedBeforeAt.GetAsGoTime())
	}
	if requestDto.CreatedAfterAt != nil {
		dbQuery = dbQuery.Where("created_at > ?", requestDto.CreatedAfterAt.GetAsGoTime())
	}

	paginatedDbQuery, responsePagination, err := dbHelper.PaginationQueryBuilder(
		dbQuery, requestDto.Pagination,
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

func (repo *BackupQueryRepo) ReadFirstTask(
	requestDto dto.ReadBackupTasksRequest,
) (taskEntity entity.BackupTask, err error) {
	requestDto.Pagination = dto.Pagination{
		PageNumber:   0,
		ItemsPerPage: 1,
	}

	responseDto, err := repo.ReadTask(requestDto)
	if err != nil {
		return taskEntity, err
	}

	if len(responseDto.Tasks) == 0 {
		return taskEntity, errors.New("BackupTaskNotFound")
	}

	return responseDto.Tasks[0], nil
}

func (repo *BackupQueryRepo) ReadTaskArchive(
	requestDto dto.ReadBackupTaskArchivesRequest,
) (responseDto dto.ReadBackupTaskArchivesResponse, err error) {
	return responseDto, nil
}
