package service

import (
	"errors"

	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/useCase"
	"github.com/speedianet/control/src/domain/valueObject"
	voHelper "github.com/speedianet/control/src/domain/valueObject/helper"
	"github.com/speedianet/control/src/infra"
	"github.com/speedianet/control/src/infra/db"
	serviceHelper "github.com/speedianet/control/src/presentation/service/helper"
)

type ContainerProfileService struct {
	persistentDbSvc           *db.PersistentDatabaseService
	containerProfileQueryRepo *infra.ContainerProfileQueryRepo
}

func NewContainerProfileService(
	persistentDbSvc *db.PersistentDatabaseService,
) *ContainerProfileService {
	return &ContainerProfileService{
		persistentDbSvc:           persistentDbSvc,
		containerProfileQueryRepo: infra.NewContainerProfileQueryRepo(persistentDbSvc),
	}
}

func (service *ContainerProfileService) Read() ServiceOutput {
	profilesList, err := useCase.ReadContainerProfiles(service.containerProfileQueryRepo)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, profilesList)
}

func (service *ContainerProfileService) parseScalingMaxDuration(
	input map[string]interface{},
) (scalingMaxDurationSecsPtr *uint, err error) {
	if input["scalingMaxDurationSecs"] != nil {
		scalingMaxDurationSecs, err := voHelper.InterfaceToUint(
			input["scalingMaxDurationSecs"],
		)
		if err != nil {
			return nil, errors.New("InvalidScalingMaxDurationSecs")
		}
		scalingMaxDurationSecsPtr = &scalingMaxDurationSecs
	}
	if input["scalingMaxDurationMinutes"] != nil {
		scalingMaxDurationMinutes, err := voHelper.InterfaceToUint16(
			input["scalingMaxDurationMinutes"],
		)
		if err != nil {
			return nil, errors.New("InvalidScalingMaxDurationMinutes")
		}
		scalingMaxDurationSecs := uint(scalingMaxDurationMinutes) * 60
		scalingMaxDurationSecsPtr = &scalingMaxDurationSecs
	}
	if input["scalingMaxDurationHours"] != nil {
		scalingMaxDurationHours, err := voHelper.InterfaceToUint8(
			input["scalingMaxDurationHours"],
		)
		if err != nil {
			return nil, errors.New("InvalidScalingMaxDurationHours")
		}
		scalingMaxDurationSecs := uint(scalingMaxDurationHours) * 3600
		scalingMaxDurationSecsPtr = &scalingMaxDurationSecs
	}

	return scalingMaxDurationSecsPtr, nil
}

func (service *ContainerProfileService) parseScalingInterval(
	input map[string]interface{},
) (scalingIntervalSecsPtr *uint, err error) {
	if input["scalingIntervalSecs"] != nil {
		scalingIntervalSecs, err := voHelper.InterfaceToUint(input["scalingIntervalSecs"])
		if err != nil {
			return nil, errors.New("InvalidScalingIntervalSecs")
		}
		scalingIntervalSecsPtr = &scalingIntervalSecs
	}
	if input["scalingIntervalMinutes"] != nil {
		scalingIntervalMinutes, err := voHelper.InterfaceToUint16(
			input["scalingIntervalMinutes"],
		)
		if err != nil {
			return nil, errors.New("InvalidScalingIntervalMinutes")
		}
		scalingIntervalSecs := uint(scalingIntervalMinutes) * 60
		scalingIntervalSecsPtr = &scalingIntervalSecs
	}
	if input["scalingIntervalHours"] != nil {
		scalingIntervalHours, err := voHelper.InterfaceToUint8(
			input["scalingIntervalHours"],
		)
		if err != nil {
			return nil, errors.New("InvalidScalingIntervalHours")
		}
		scalingIntervalSecs := uint(scalingIntervalHours) * 3600
		scalingIntervalSecsPtr = &scalingIntervalSecs
	}

	return scalingIntervalSecsPtr, nil
}

func (service *ContainerProfileService) Create(
	input map[string]interface{},
) ServiceOutput {
	requiredParams := []string{"name", "baseSpecs"}

	err := serviceHelper.RequiredParamsInspector(input, requiredParams)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	profileName, err := valueObject.NewContainerProfileName(input["name"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	baseSpecs, assertOk := input["baseSpecs"].(valueObject.ContainerSpecs)
	if !assertOk {
		return NewServiceOutput(UserError, "InvalidBaseSpecsStructure")
	}

	var maxSpecsPtr *valueObject.ContainerSpecs
	if input["maxSpecs"] != nil {
		maxSpecs, assertOk := input["maxSpecs"].(valueObject.ContainerSpecs)
		if !assertOk {
			return NewServiceOutput(UserError, "InvalidMaxSpecsStructure")
		}
		maxSpecsPtr = &maxSpecs
	}

	var scalingPolicyPtr *valueObject.ScalingPolicy
	if input["scalingPolicy"] != nil {
		scalingPolicy, err := valueObject.NewScalingPolicy(input["scalingPolicy"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		scalingPolicyPtr = &scalingPolicy
	}

	var scalingThresholdPtr *uint
	if input["scalingThreshold"] != nil {
		scalingThreshold, err := voHelper.InterfaceToUint(input["scalingThreshold"])
		if err != nil {
			return NewServiceOutput(UserError, "InvalidScalingThreshold")
		}
		scalingThresholdPtr = &scalingThreshold
	}

	scalingMaxDurationSecsPtr, err := service.parseScalingMaxDuration(input)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	scalingIntervalSecsPtr, err := service.parseScalingInterval(input)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	var hostMinCapacityPercentPtr *valueObject.HostMinCapacity
	if input["hostMinCapacityPercent"] != nil {
		hostMinCapacityPercent, err := valueObject.NewHostMinCapacity(
			input["hostMinCapacityPercent"],
		)
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		hostMinCapacityPercentPtr = &hostMinCapacityPercent
	}

	createDto := dto.NewCreateContainerProfile(
		profileName, baseSpecs, maxSpecsPtr, scalingPolicyPtr, scalingThresholdPtr,
		scalingMaxDurationSecsPtr, scalingIntervalSecsPtr, hostMinCapacityPercentPtr,
	)

	containerProfileCmdRepo := infra.NewContainerProfileCmdRepo(service.persistentDbSvc)

	err = useCase.CreateContainerProfile(containerProfileCmdRepo, createDto)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Created, "ContainerProfileCreated")
}

func (service *ContainerProfileService) Update(
	input map[string]interface{},
) ServiceOutput {
	if input["id"] != nil {
		input["profileId"] = input["id"]
	}
	requiredParams := []string{"profileId"}

	err := serviceHelper.RequiredParamsInspector(input, requiredParams)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	profileId, err := valueObject.NewContainerProfileId(input["profileId"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	var profileNamePtr *valueObject.ContainerProfileName
	if input["name"] != nil {
		profileName, err := valueObject.NewContainerProfileName(input["name"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		profileNamePtr = &profileName
	}

	var baseSpecsPtr *valueObject.ContainerSpecs
	if input["baseSpecs"] != nil {
		baseSpecs, assertOk := input["baseSpecs"].(valueObject.ContainerSpecs)
		if !assertOk {
			return NewServiceOutput(UserError, "InvalidBaseSpecsStructure")
		}
		baseSpecsPtr = &baseSpecs
	}

	var maxSpecsPtr *valueObject.ContainerSpecs
	if input["maxSpecs"] != nil {
		maxSpecs, assertOk := input["maxSpecs"].(valueObject.ContainerSpecs)
		if !assertOk {
			return NewServiceOutput(UserError, "InvalidMaxSpecsStructure")
		}
		maxSpecsPtr = &maxSpecs
	}

	var scalingPolicyPtr *valueObject.ScalingPolicy
	if input["scalingPolicy"] != nil {
		scalingPolicy, err := valueObject.NewScalingPolicy(input["scalingPolicy"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		scalingPolicyPtr = &scalingPolicy
	}

	var scalingThresholdPtr *uint
	if input["scalingThreshold"] != nil {
		scalingThreshold, err := voHelper.InterfaceToUint(input["scalingThreshold"])
		if err != nil {
			return NewServiceOutput(UserError, "InvalidScalingThreshold")
		}
		scalingThresholdPtr = &scalingThreshold
	}

	scalingMaxDurationSecsPtr, err := service.parseScalingMaxDuration(input)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	scalingIntervalSecsPtr, err := service.parseScalingInterval(input)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	var hostMinCapacityPercentPtr *valueObject.HostMinCapacity
	if input["hostMinCapacityPercent"] != nil {
		hostMinCapacityPercent, err := valueObject.NewHostMinCapacity(
			input["hostMinCapacityPercent"],
		)
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		hostMinCapacityPercentPtr = &hostMinCapacityPercent
	}

	updateDto := dto.NewUpdateContainerProfile(
		profileId, profileNamePtr, baseSpecsPtr, maxSpecsPtr, scalingPolicyPtr,
		scalingThresholdPtr, scalingMaxDurationSecsPtr, scalingIntervalSecsPtr,
		hostMinCapacityPercentPtr,
	)

	containerProfileCmdRepo := infra.NewContainerProfileCmdRepo(service.persistentDbSvc)
	containerQueryRepo := infra.NewContainerQueryRepo(service.persistentDbSvc)
	containerCmdRepo := infra.NewContainerCmdRepo(service.persistentDbSvc)

	err = useCase.UpdateContainerProfile(
		service.containerProfileQueryRepo, containerProfileCmdRepo, containerQueryRepo,
		containerCmdRepo, updateDto,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, "ContainerProfileUpdated")
}

func (service *ContainerProfileService) Delete(
	input map[string]interface{},
) ServiceOutput {
	if input["id"] != nil {
		input["profileId"] = input["id"]
	}
	requiredParams := []string{"profileId"}

	err := serviceHelper.RequiredParamsInspector(input, requiredParams)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	profileId, err := valueObject.NewContainerProfileId(input["profileId"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	containerProfileCmdRepo := infra.NewContainerProfileCmdRepo(service.persistentDbSvc)
	containerQueryRepo := infra.NewContainerQueryRepo(service.persistentDbSvc)
	containerCmdRepo := infra.NewContainerCmdRepo(service.persistentDbSvc)

	err = useCase.DeleteContainerProfile(
		service.containerProfileQueryRepo, containerProfileCmdRepo, containerQueryRepo,
		containerCmdRepo, profileId,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, "ContainerProfileDeleted")
}
