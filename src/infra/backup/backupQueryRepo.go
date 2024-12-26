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
