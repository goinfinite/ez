package infra

import (
	"errors"
	"log"

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
	var accEntities []entity.Account

	var accModels []dbModel.Account

	err := repo.persistentDbSvc.Handler.Model(&dbModel.Account{}).
		Preload("Quota").
		Preload("QuotaUsage").Find(&accModels).Error
	if err != nil {
		return accEntities, errors.New("DatabaseQueryAccountsError")
	}

	for _, accModel := range accModels {
		accEntity, err := accModel.ToEntity()
		if err != nil {
			log.Printf("AccountModelToEntityError: %v", err.Error())
			continue
		}

		accEntities = append(accEntities, accEntity)
	}

	return accEntities, nil
}

func (repo *AccountQueryRepo) ReadByUsername(
	username valueObject.Username,
) (entity.Account, error) {
	accEntities, err := repo.Read()
	if err != nil {
		return entity.Account{}, errors.New("AccountQueryError")
	}

	for _, accEntity := range accEntities {
		if accEntity.Username.String() == username.String() {
			return accEntity, nil
		}
	}

	return entity.Account{}, errors.New("AccountNotFound")
}

func (repo *AccountQueryRepo) ReadById(
	accountId valueObject.AccountId,
) (entity.Account, error) {
	accEntities, err := repo.Read()
	if err != nil {
		return entity.Account{}, errors.New("AccountQueryError")
	}

	for _, accEntity := range accEntities {
		if accEntity.Id.String() == accountId.String() {
			return accEntity, nil
		}
	}

	return entity.Account{}, errors.New("AccountNotFound")
}
