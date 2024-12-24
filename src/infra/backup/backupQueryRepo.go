package backupInfra

import (
	"errors"
	"log/slog"
	"math"

	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/entity"
	"github.com/goinfinite/ez/src/infra/db"
	dbModel "github.com/goinfinite/ez/src/infra/db/model"
	"github.com/iancoleman/strcase"
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

	var itemsTotal int64
	err = dbQuery.Count(&itemsTotal).Error
	if err != nil {
		return responseDto, errors.New("CountItemsTotalError: " + err.Error())
	}

	dbQuery = dbQuery.Limit(int(readDto.Pagination.ItemsPerPage))
	if readDto.Pagination.LastSeenId == nil {
		offset := int(readDto.Pagination.PageNumber) * int(readDto.Pagination.ItemsPerPage)
		dbQuery = dbQuery.Offset(offset)
	} else {
		dbQuery = dbQuery.Where("id > ?", readDto.Pagination.LastSeenId.String())
	}
	if readDto.Pagination.SortBy != nil {
		orderStatement := readDto.Pagination.SortBy.String()
		orderStatement = strcase.ToSnake(orderStatement)
		if orderStatement == "id" {
			orderStatement = "ID"
		}

		if readDto.Pagination.SortDirection != nil {
			orderStatement += " " + readDto.Pagination.SortDirection.String()
		}

		dbQuery = dbQuery.Order(orderStatement)
	}

	backupDestinationModels := []dbModel.BackupDestination{}
	err = dbQuery.Find(&backupDestinationModels).Error
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

	itemsTotalUint := uint64(itemsTotal)
	pagesTotal := uint32(
		math.Ceil(float64(itemsTotal) / float64(readDto.Pagination.ItemsPerPage)),
	)
	responsePagination := dto.Pagination{
		PageNumber:    readDto.Pagination.PageNumber,
		ItemsPerPage:  readDto.Pagination.ItemsPerPage,
		SortBy:        readDto.Pagination.SortBy,
		SortDirection: readDto.Pagination.SortDirection,
		PagesTotal:    &pagesTotal,
		ItemsTotal:    &itemsTotalUint,
	}

	return dto.ReadBackupDestinationsResponse{
		Pagination:   responsePagination,
		Destinations: backupDestinationEntities,
	}, nil
}
