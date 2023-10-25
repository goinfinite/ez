package infra

import (
	"errors"
	"log"

	"github.com/goinfinite/fleet/src/domain/entity"
	"github.com/goinfinite/fleet/src/domain/valueObject"
	"github.com/goinfinite/fleet/src/infra/db"
	dbModel "github.com/goinfinite/fleet/src/infra/db/model"
)

type AccQueryRepo struct {
	dbSvc *db.DatabaseService
}

func NewAccQueryRepo(dbSvc *db.DatabaseService) *AccQueryRepo {
	return &AccQueryRepo{dbSvc: dbSvc}
}

func (repo AccQueryRepo) Get() ([]entity.Account, error) {
	var accEntities []entity.Account

	var accModels []dbModel.Account

	err := repo.dbSvc.Orm.Model(&dbModel.Account{}).
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
