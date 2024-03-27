package infra

import (
	"errors"

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

func (repo *ContainerProfileCmdRepo) Add(
	addDto dto.AddContainerProfile,
) error {
	containerProfileModel, err := dbModel.ContainerProfile{}.AddDtoToModel(addDto)
	if err != nil {
		return err
	}

	err = repo.persistentDbSvc.Handler.Create(&containerProfileModel).Error
	if err != nil {
		return err
	}

	return nil
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

	err := repo.persistentDbSvc.Handler.Table(dbModel.ContainerProfile{}.TableName()).
		Where("id = ?", updateDto.Id.String()).
		Updates(updateMap).Error
	if err != nil {
		return err
	}

	return nil
}

func (repo *ContainerProfileCmdRepo) Delete(
	profileId valueObject.ContainerProfileId,
) error {
	err := repo.persistentDbSvc.Handler.Delete(
		dbModel.ContainerProfile{},
		profileId.Get(),
	).Error
	if err != nil {
		return errors.New("DeleteContainerProfileDbError")
	}

	return nil
}
