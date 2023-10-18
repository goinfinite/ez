package infra

import (
	"errors"

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
	updateMap := map[string]interface{}{}

	if updateDto.Name != nil {
		updateMap["name"] = updateDto.Name.String()
	}

	if updateDto.BaseSpecs != nil {
		updateMap["base_specs"] = updateDto.BaseSpecs.String()
	}

	if updateDto.MaxSpecs != nil {
		updateMap["max_specs"] = updateDto.MaxSpecs.String()
	}

	if updateDto.ScalingPolicy != nil {
		updateMap["scaling_policy"] = updateDto.ScalingPolicy.String()
	}

	if updateDto.ScalingThreshold != nil {
		updateMap["scaling_threshold"] = uint64(*updateDto.ScalingThreshold)
	}

	if updateDto.ScalingMaxDurationSecs != nil {
		updateMap["scaling_max_duration_secs"] = uint64(
			*updateDto.ScalingMaxDurationSecs,
		)
	}

	if updateDto.ScalingIntervalSecs != nil {
		updateMap["scaling_interval_secs"] = uint64(
			*updateDto.ScalingIntervalSecs,
		)
	}

	if updateDto.HostMinCapacityPercent != nil {
		updateMap["host_min_capacity_percent"] = updateDto.HostMinCapacityPercent.Get()
	}

	err := repo.dbSvc.Table(dbModel.ResourceProfile{}.TableName()).
		Where("id = ?", updateDto.Id.String()).
		Updates(updateMap).Error
	if err != nil {
		return err
	}

	return nil
}

func (repo ResourceProfileCmdRepo) Delete(
	profileId valueObject.ResourceProfileId,
) error {
	err := repo.dbSvc.Delete(dbModel.ResourceProfile{}, profileId.Get()).Error
	if err != nil {
		return errors.New("DeleteResourceProfileDbError")
	}

	return nil
}
