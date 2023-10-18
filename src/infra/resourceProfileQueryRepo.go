package infra

import (
	"errors"
	"log"

	"github.com/speedianet/sfm/src/domain/entity"
	"github.com/speedianet/sfm/src/domain/valueObject"
	dbModel "github.com/speedianet/sfm/src/infra/db/model"
	"gorm.io/gorm"
)

type ResourceProfileQueryRepo struct {
	dbSvc *gorm.DB
}

func NewResourceProfileQueryRepo(dbSvc *gorm.DB) *ResourceProfileQueryRepo {
	return &ResourceProfileQueryRepo{dbSvc: dbSvc}
}

func (repo ResourceProfileQueryRepo) Get() ([]entity.ResourceProfile, error) {
	var profileEntities []entity.ResourceProfile
	var profileModels []dbModel.ResourceProfile

	err := repo.dbSvc.Model(&dbModel.ResourceProfile{}).Find(&profileModels).Error
	if err != nil {
		return profileEntities, errors.New("DatabaseQueryResourceProfileError")
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

func (repo ResourceProfileQueryRepo) GetById(
	id valueObject.ResourceProfileId,
) (entity.ResourceProfile, error) {
	return entity.ResourceProfile{}, nil
}
