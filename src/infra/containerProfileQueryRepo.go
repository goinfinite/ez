package infra

import (
	"errors"
	"log"

	"github.com/speedianet/sfm/src/domain/entity"
	"github.com/speedianet/sfm/src/domain/valueObject"
	dbModel "github.com/speedianet/sfm/src/infra/db/model"
	"gorm.io/gorm"
)

type ContainerProfileQueryRepo struct {
	dbSvc *gorm.DB
}

func NewContainerProfileQueryRepo(dbSvc *gorm.DB) *ContainerProfileQueryRepo {
	return &ContainerProfileQueryRepo{dbSvc: dbSvc}
}

func (repo ContainerProfileQueryRepo) Get() ([]entity.ContainerProfile, error) {
	var profileEntities []entity.ContainerProfile
	var profileModels []dbModel.ContainerProfile

	err := repo.dbSvc.Model(&dbModel.ContainerProfile{}).Find(&profileModels).Error
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
	var profileModel dbModel.ContainerProfile
	err := repo.dbSvc.Model(&dbModel.ContainerProfile{}).First(&profileModel).Error
	if err != nil {
		return entity.ContainerProfile{}, errors.New("DbQueryContainerProfileError")
	}

	profileEntity, err := profileModel.ToEntity()
	if err != nil {
		return entity.ContainerProfile{}, errors.New("ProfileModelToEntityError")
	}

	return profileEntity, nil
}
