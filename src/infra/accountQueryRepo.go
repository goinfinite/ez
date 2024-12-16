package infra

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/ez/src/domain/entity"
	"github.com/goinfinite/ez/src/domain/valueObject"
	"github.com/goinfinite/ez/src/infra/db"
	dbModel "github.com/goinfinite/ez/src/infra/db/model"
)

type AccountQueryRepo struct {
	persistentDbSvc *db.PersistentDatabaseService
}

func NewAccountQueryRepo(persistentDbSvc *db.PersistentDatabaseService) *AccountQueryRepo {
	return &AccountQueryRepo{persistentDbSvc: persistentDbSvc}
}

func (repo *AccountQueryRepo) Read() ([]entity.Account, error) {
	accountEntities := []entity.Account{}

	var accountModels []dbModel.Account
	err := repo.persistentDbSvc.Handler.
		Model(&dbModel.Account{}).
		Preload("Quota").
		Preload("QuotaUsage").Find(&accountModels).Error
	if err != nil {
		return accountEntities, errors.New("QueryAccountsError: " + err.Error())
	}

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

	return accountEntities, nil
}

func (repo *AccountQueryRepo) ReadByUsername(
	username valueObject.UnixUsername,
) (accountEntity entity.Account, err error) {
	accountEntities, err := repo.Read()
	if err != nil {
		return accountEntity, errors.New("ReadAccountsError: " + err.Error())
	}

	usernameStr := username.String()
	for _, accountEntity := range accountEntities {
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
	accountEntities, err := repo.Read()
	if err != nil {
		return accountEntity, errors.New("ReadAccountsError: " + err.Error())
	}

	for _, accountEntity := range accountEntities {
		if accountEntity.Id.String() != accountId.String() {
			continue
		}

		return accountEntity, nil
	}

	return accountEntity, errors.New("AccountNotFound")
}
