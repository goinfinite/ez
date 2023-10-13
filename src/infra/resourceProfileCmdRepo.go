package infra

import (
	"github.com/speedianet/sfm/src/domain/dto"
	"github.com/speedianet/sfm/src/domain/valueObject"
	dbModel "github.com/speedianet/sfm/src/infra/db/model"
	"gorm.io/gorm"
)

type ResourceProfileCmdRepo struct {
	dbSvc *gorm.DB
}

func NewResourceProfileCmdRepo(dbSvc *gorm.DB) *ResourceProfileCmdRepo {
	return &ResourceProfileCmdRepo{dbSvc: dbSvc}
}

func (repo ResourceProfileCmdRepo) Add(
	addDto dto.AddResourceProfile,
) error {
	resourceProfileModel, err := dbModel.ResourceProfile{}.FromAddDtoToModel(addDto)
	if err != nil {
		return err
	}

	err = repo.dbSvc.Create(&resourceProfileModel).Error
	if err != nil {
		return err
	}

	return nil
}

func (repo ResourceProfileCmdRepo) Update(
	updateDto dto.UpdateResourceProfile,
) error {
	return nil
}

func (repo ResourceProfileCmdRepo) Delete(
	profileId valueObject.ResourceProfileId,
) error {
	return nil
}
