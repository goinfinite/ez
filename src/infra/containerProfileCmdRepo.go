package infra

import (
	"errors"

	"github.com/speedianet/sfm/src/domain/dto"
	"github.com/speedianet/sfm/src/domain/valueObject"
	dbModel "github.com/speedianet/sfm/src/infra/db/model"
	"gorm.io/gorm"
)

type ContainerProfileCmdRepo struct {
	dbSvc *gorm.DB
}

func NewContainerProfileCmdRepo(dbSvc *gorm.DB) *ContainerProfileCmdRepo {
	return &ContainerProfileCmdRepo{dbSvc: dbSvc}
}

func (repo ContainerProfileCmdRepo) Add(
	addDto dto.AddContainerProfile,
) error {
	containerProfileModel, err := dbModel.ContainerProfile{}.FromAddDtoToModel(addDto)
	if err != nil {
		return err
	}

	err = repo.dbSvc.Create(&containerProfileModel).Error
	if err != nil {
		return err
	}

	return nil
}

func (repo ContainerProfileCmdRepo) Update(
	updateDto dto.UpdateContainerProfile,
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

	err := repo.dbSvc.Table(dbModel.ContainerProfile{}.TableName()).
		Where("id = ?", updateDto.Id.String()).
		Updates(updateMap).Error
	if err != nil {
		return err
	}

	return nil
}

func (repo ContainerProfileCmdRepo) Delete(
	profileId valueObject.ContainerProfileId,
) error {
	err := repo.dbSvc.Delete(dbModel.ContainerProfile{}, profileId.Get()).Error
	if err != nil {
		return errors.New("DeleteContainerProfileDbError")
	}

	return nil
}
