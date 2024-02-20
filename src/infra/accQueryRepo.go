package infra

import (
	"errors"
	"log"

	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/valueObject"
	"github.com/speedianet/control/src/infra/db"
	dbModel "github.com/speedianet/control/src/infra/db/model"
)

type AccQueryRepo struct {
	persistDbSvc *db.PersistentDatabaseService
}

func NewAccQueryRepo(persistDbSvc *db.PersistentDatabaseService) *AccQueryRepo {
	return &AccQueryRepo{persistDbSvc: persistDbSvc}
}

func (repo AccQueryRepo) Get() ([]entity.Account, error) {
	var accEntities []entity.Account

	var accModels []dbModel.Account

	err := repo.persistDbSvc.Orm.Model(&dbModel.Account{}).
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

func (repo AccQueryRepo) GetByUsername(
	username valueObject.Username,
) (entity.Account, error) {
	accEntities, err := repo.Get()
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

func (repo AccQueryRepo) GetById(
	accountId valueObject.AccountId,
) (entity.Account, error) {
	accEntities, err := repo.Get()
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
