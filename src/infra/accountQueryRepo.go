package infra

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/entity"
	"github.com/goinfinite/ez/src/domain/valueObject"
	"github.com/goinfinite/ez/src/infra/db"
	dbHelper "github.com/goinfinite/ez/src/infra/db/helper"
	dbModel "github.com/goinfinite/ez/src/infra/db/model"
)

type AccountQueryRepo struct {
	persistentDbSvc *db.PersistentDatabaseService
}

func NewAccountQueryRepo(persistentDbSvc *db.PersistentDatabaseService) *AccountQueryRepo {
	return &AccountQueryRepo{persistentDbSvc: persistentDbSvc}
}

func (repo *AccountQueryRepo) Read(
	requestDto dto.ReadAccountsRequest,
) (responseDto dto.ReadAccountsResponse, err error) {
	model := dbModel.Account{}
	if requestDto.AccountId != nil {
		model.ID = requestDto.AccountId.Uint64()
	}
	if requestDto.AccountUsername != nil {
		model.Username = requestDto.AccountUsername.String()
	}

	dbQuery := repo.persistentDbSvc.Handler.Model(model).Where(&model)

	if requestDto.Pagination.SortBy != nil {
		sortByStr := requestDto.Pagination.SortBy.String()
		switch sortByStr {
		case "id", "accountId":
			sortByStr = "ID"
		case "username", "accountUsername":
			sortByStr = "Username"
		}

		sortBy, err := valueObject.NewPaginationSortBy(sortByStr)
		if err == nil {
			requestDto.Pagination.SortBy = &sortBy
		}
	}
	paginatedDbQuery, responsePagination, err := dbHelper.PaginationQueryBuilder(
		dbQuery, requestDto.Pagination,
	)
	if err != nil {
		return responseDto, errors.New("PaginationQueryBuilderError: " + err.Error())
	}

	var accountModels []dbModel.Account
	err = paginatedDbQuery.
		Preload("Quota").
		Preload("QuotaUsage").
		Find(&accountModels).Error
	if err != nil {
		return responseDto, errors.New("QueryAccountsError: " + err.Error())
	}

	accountEntities := []entity.Account{}
	for _, accountModel := range accountModels {
		accountEntity, err := accountModel.ToEntity()
		if err != nil {
			slog.Debug(
				"ModelToEntityError",
				slog.Any("error", err.Error()),
				slog.Uint64("accountId", accountModel.ID),
			)
			continue
		}

		accountEntities = append(accountEntities, accountEntity)
	}

	return dto.ReadAccountsResponse{
		Pagination: responsePagination,
		Accounts:   accountEntities,
	}, nil
}

func (repo *AccountQueryRepo) ReadByUsername(
	username valueObject.UnixUsername,
) (accountEntity entity.Account, err error) {
	requestDto := dto.ReadAccountsRequest{
		Pagination: dto.Pagination{
			PageNumber:   0,
			ItemsPerPage: 1,
		},
		AccountUsername: &username,
	}
	responseDto, err := repo.Read(requestDto)
	if err != nil {
		return accountEntity, errors.New("ReadAccountsError: " + err.Error())
	}

	usernameStr := username.String()
	for _, accountEntity := range responseDto.Accounts {
		if accountEntity.Username.String() != usernameStr {
			continue
		}

		return accountEntity, nil
	}

	return accountEntity, errors.New("AccountNotFound")
}

func (repo *AccountQueryRepo) ReadById(
	accountId valueObject.AccountId,
) (accountEntity entity.Account, err error) {
	requestDto := dto.ReadAccountsRequest{
		Pagination: dto.Pagination{
			PageNumber:   0,
			ItemsPerPage: 1,
		},
		AccountId: &accountId,
	}
	responseDto, err := repo.Read(requestDto)
	if err != nil {
		return accountEntity, errors.New("ReadAccountsError: " + err.Error())
	}

	for _, accountEntity := range responseDto.Accounts {
		if accountEntity.Id.String() != accountId.String() {
			continue
		}

		return accountEntity, nil
	}

	return accountEntity, errors.New("AccountNotFound")
}
