package infra

import (
	"errors"
	"log"
	"log/slog"

	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/valueObject"
	"github.com/speedianet/control/src/infra/db"
	dbModel "github.com/speedianet/control/src/infra/db/model"
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
		return accountEntities, errors.New("DatabaseQueryAccountsError")
	}

	for _, accountModel := range accountModels {
		accountEntity, err := accountModel.ToEntity()
		if err != nil {
			slog.Debug("ModelToEntityError",
				slog.Any("error", err.Error()),
				slog.Uint64("accountId", uint64(accountModel.ID)),
			)
			log.Printf("AccountModelToEntityError: %v", err.Error())
			continue
		}

		accountEntities = append(accountEntities, accountEntity)
	}

	return accountEntities, nil
}

func (repo *AccountQueryRepo) ReadByUsername(
	username valueObject.Username,
) (entity.Account, error) {
	accountEntities, err := repo.Read()
	if err != nil {
		return entity.Account{}, errors.New("AccountQueryError")
	}

	for _, accountEntity := range accountEntities {
		if accountEntity.Username.String() == username.String() {
			return accountEntity, nil
		}
	}

	return entity.Account{}, errors.New("AccountNotFound")
}

func (repo *AccountQueryRepo) ReadById(
	accountId valueObject.AccountId,
) (entity.Account, error) {
	accountEntities, err := repo.Read()
	if err != nil {
		return entity.Account{}, errors.New("AccountQueryError")
	}

	for _, accountEntity := range accountEntities {
		if accountEntity.Id.String() == accountId.String() {
			return accountEntity, nil
		}
	}

	return entity.Account{}, errors.New("AccountNotFound")
}
