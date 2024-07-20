package infra

import (
	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/valueObject"
	"github.com/speedianet/control/src/infra/db"
	dbModel "github.com/speedianet/control/src/infra/db/model"
)

type ContainerProfileCmdRepo struct {
	persistentDbSvc *db.PersistentDatabaseService
}

func NewContainerProfileCmdRepo(
	persistentDbSvc *db.PersistentDatabaseService,
) *ContainerProfileCmdRepo {
	return &ContainerProfileCmdRepo{persistentDbSvc: persistentDbSvc}
}

func (repo *ContainerProfileCmdRepo) Create(
	createDto dto.CreateContainerProfile,
) error {
	containerProfileModel, err := dbModel.ContainerProfile{}.AddDtoToModel(createDto)
	if err != nil {
		return err
	}

	return repo.persistentDbSvc.Handler.Create(&containerProfileModel).Error
}

func (repo *ContainerProfileCmdRepo) Update(
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
		updateMap["scaling_threshold"] = *updateDto.ScalingThreshold
	}

	if updateDto.ScalingMaxDurationSecs != nil {
		updateMap["scaling_max_duration_secs"] = *updateDto.ScalingMaxDurationSecs
	}

	if updateDto.ScalingIntervalSecs != nil {
		updateMap["scaling_interval_secs"] = *updateDto.ScalingIntervalSecs
	}

	if updateDto.HostMinCapacityPercent != nil {
		updateMap["host_min_capacity_percent"] = updateDto.HostMinCapacityPercent.Float64()
	}

	return repo.persistentDbSvc.Handler.
		Model(&dbModel.ContainerProfile{}).
		Where("id = ?", updateDto.Id.String()).
		Updates(updateMap).Error
}

func (repo *ContainerProfileCmdRepo) Delete(
	profileId valueObject.ContainerProfileId,
) error {
	return repo.persistentDbSvc.Handler.
		Model(&dbModel.ContainerProfile{}).
		Delete(dbModel.ContainerProfile{}, profileId.Uint64()).Error
}
