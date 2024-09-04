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
) (profileId valueObject.ContainerProfileId, err error) {
	var maxSpecsPtr *string
	if createDto.MaxSpecs != nil {
		maxSpecs := createDto.MaxSpecs.String()
		maxSpecsPtr = &maxSpecs
	}

	var scalingPolicyPtr *string
	if createDto.ScalingPolicy != nil {
		scalingPolicy := createDto.ScalingPolicy.String()
		scalingPolicyPtr = &scalingPolicy
	}

	var hostMinCapacityPercentPtr *uint8
	if createDto.HostMinCapacityPercent != nil {
		hostMinCapacityPercentUint8 := createDto.HostMinCapacityPercent.Uint8()
		hostMinCapacityPercentPtr = &hostMinCapacityPercentUint8
	}

	createModel := dbModel.NewContainerProfile(
		createDto.AccountId.Uint64(), createDto.Name.String(),
		createDto.BaseSpecs.String(), maxSpecsPtr, scalingPolicyPtr,
		createDto.ScalingThreshold, createDto.ScalingMaxDurationSecs,
		createDto.ScalingIntervalSecs, hostMinCapacityPercentPtr,
	)

	err = repo.persistentDbSvc.Handler.Create(&createModel).Error
	if err != nil {
		return profileId, err
	}

	return valueObject.NewContainerProfileId(createModel.ID)
}

func (repo *ContainerProfileCmdRepo) Update(
	updateDto dto.UpdateContainerProfile,
) error {
	updateMap := map[string]interface{}{
		"account_id": updateDto.AccountId.Uint64(),
	}

	if updateDto.Name != nil {
		updateMap["name"] = updateDto.Name.String()
	}

	if updateDto.BaseSpecs != nil {
		updateMap["base_specs"] = updateDto.BaseSpecs.String()
	}

	if updateDto.MaxSpecs != nil {
		maxSpecsStr := updateDto.MaxSpecs.String()
		updateMap["max_specs"] = maxSpecsStr
		if maxSpecsStr == "0:0:0" {
			updateMap["max_specs"] = nil
		}
	}

	if updateDto.ScalingPolicy != nil {
		updateMap["scaling_policy"] = updateDto.ScalingPolicy.String()
	}

	if updateDto.ScalingThreshold != nil {
		thresholdUint := *updateDto.ScalingThreshold
		updateMap["scaling_threshold"] = thresholdUint
		if thresholdUint == 0 {
			updateMap["scaling_threshold"] = nil
		}
	}

	if updateDto.ScalingMaxDurationSecs != nil {
		durationUint := *updateDto.ScalingMaxDurationSecs
		updateMap["scaling_max_duration_secs"] = durationUint
		if durationUint == 0 {
			updateMap["scaling_max_duration_secs"] = nil
		}
	}

	if updateDto.ScalingIntervalSecs != nil {
		intervalUint := *updateDto.ScalingIntervalSecs
		updateMap["scaling_interval_secs"] = intervalUint
		if intervalUint == 0 {
			updateMap["scaling_interval_secs"] = nil
		}
	}

	if updateDto.HostMinCapacityPercent != nil {
		percentUint := updateDto.HostMinCapacityPercent.Uint8()
		updateMap["host_min_capacity_percent"] = percentUint
		if percentUint == 0 {
			updateMap["host_min_capacity_percent"] = nil
		}
	}

	return repo.persistentDbSvc.Handler.
		Model(&dbModel.ContainerProfile{}).
		Where("id = ?", updateDto.ProfileId.String()).
		Updates(updateMap).Error
}

func (repo *ContainerProfileCmdRepo) Delete(
	deleteDto dto.DeleteContainerProfile,
) error {
	return repo.persistentDbSvc.Handler.
		Model(&dbModel.ContainerProfile{}).
		Delete(dbModel.ContainerProfile{}, deleteDto.ProfileId.Uint64()).Error
}
