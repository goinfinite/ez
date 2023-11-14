package infra

import (
	"errors"
	"log"

	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/valueObject"
	"github.com/speedianet/control/src/infra/db"
	dbModel "github.com/speedianet/control/src/infra/db/model"
)

type ContainerProfileQueryRepo struct {
	dbSvc *db.DatabaseService
}

func NewContainerProfileQueryRepo(dbSvc *db.DatabaseService) *ContainerProfileQueryRepo {
	return &ContainerProfileQueryRepo{dbSvc: dbSvc}
}

func (repo ContainerProfileQueryRepo) Get() ([]entity.ContainerProfile, error) {
	var profileEntities []entity.ContainerProfile
	var profileModels []dbModel.ContainerProfile

	err := repo.dbSvc.Orm.Model(&dbModel.ContainerProfile{}).Find(&profileModels).Error
	if err != nil {
		return profileEntities, errors.New("DbQueryContainerProfilesError")
	}

	for _, profileModel := range profileModels {
		profileEntity, err := profileModel.ToEntity()
		if err != nil {
			log.Printf("ProfileModelToEntityError: %v", err.Error())
			continue
		}

		profileEntities = append(profileEntities, profileEntity)
	}

	return profileEntities, nil
}

func (repo ContainerProfileQueryRepo) GetById(
	id valueObject.ContainerProfileId,
) (entity.ContainerProfile, error) {
	profileModel := dbModel.ContainerProfile{ID: uint(id.Get())}
	err := repo.dbSvc.Orm.Model(&profileModel).First(&profileModel).Error
	if err != nil {
		return entity.ContainerProfile{}, err
	}

	profileEntity, err := profileModel.ToEntity()
	if err != nil {
		return entity.ContainerProfile{}, errors.New("ProfileModelToEntityError")
	}

	return profileEntity, nil
}
