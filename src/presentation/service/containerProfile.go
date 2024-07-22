package service

import (
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
	if _, exists := input["maxSpecs"]; exists {
		maxSpecs, assertOk := input["maxSpecs"].(valueObject.ContainerSpecs)
		if !assertOk {
			return NewServiceOutput(UserError, "InvalidMaxSpecsStructure")
		}
		maxSpecsPtr = &maxSpecs
	}

	var scalingPolicyPtr *valueObject.ScalingPolicy
	if _, exists := input["scalingPolicy"]; exists {
		scalingPolicy, err := valueObject.NewScalingPolicy(input["scalingPolicy"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		scalingPolicyPtr = &scalingPolicy
	}

	var scalingThresholdPtr *uint
	if _, exists := input["scalingThreshold"]; exists {
		scalingThreshold, err := voHelper.InterfaceToUint(input["scalingThreshold"])
		if err != nil {
			return NewServiceOutput(UserError, "InvalidScalingThreshold")
		}
		scalingThresholdPtr = &scalingThreshold
	}

	var scalingMaxDurationSecsPtr *uint
	if _, exists := input["scalingMaxDurationSecs"]; exists {
		scalingMaxDurationSecs, err := voHelper.InterfaceToUint(
			input["scalingMaxDurationSecs"],
		)
		if err != nil {
			return NewServiceOutput(UserError, "InvalidScalingMaxDurationSecs")
		}
		scalingMaxDurationSecsPtr = &scalingMaxDurationSecs
	}

	var scalingIntervalSecsPtr *uint
	if _, exists := input["scalingIntervalSecs"]; exists {
		scalingIntervalSecs, err := voHelper.InterfaceToUint(input["scalingIntervalSecs"])
		if err != nil {
			return NewServiceOutput(UserError, "InvalidScalingIntervalSecs")
		}
		scalingIntervalSecsPtr = &scalingIntervalSecs
	}

	var hostMinCapacityPercentPtr *valueObject.HostMinCapacity
	if _, exists := input["hostMinCapacityPercent"]; exists {
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
	requiredParams := []string{"id"}

	err := serviceHelper.RequiredParamsInspector(input, requiredParams)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	profileId, err := valueObject.NewContainerProfileId(input["id"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	var profileNamePtr *valueObject.ContainerProfileName
	if _, exists := input["name"]; exists {
		profileName, err := valueObject.NewContainerProfileName(input["name"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		profileNamePtr = &profileName
	}

	var baseSpecsPtr *valueObject.ContainerSpecs
	if _, exists := input["baseSpecs"]; exists {
		baseSpecs, assertOk := input["baseSpecs"].(valueObject.ContainerSpecs)
		if !assertOk {
			return NewServiceOutput(UserError, "InvalidBaseSpecsStructure")
		}
		baseSpecsPtr = &baseSpecs
	}

	var maxSpecsPtr *valueObject.ContainerSpecs
	if _, exists := input["maxSpecs"]; exists {
		maxSpecs, assertOk := input["maxSpecs"].(valueObject.ContainerSpecs)
		if !assertOk {
			return NewServiceOutput(UserError, "InvalidMaxSpecsStructure")
		}
		maxSpecsPtr = &maxSpecs
	}

	var scalingPolicyPtr *valueObject.ScalingPolicy
	if _, exists := input["scalingPolicy"]; exists {
		scalingPolicy, err := valueObject.NewScalingPolicy(input["scalingPolicy"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		scalingPolicyPtr = &scalingPolicy
	}

	var scalingThresholdPtr *uint
	if _, exists := input["scalingThreshold"]; exists {
		scalingThreshold, err := voHelper.InterfaceToUint(input["scalingThreshold"])
		if err != nil {
			return NewServiceOutput(UserError, "InvalidScalingThreshold")
		}
		scalingThresholdPtr = &scalingThreshold
	}

	var scalingMaxDurationSecsPtr *uint
	if _, exists := input["scalingMaxDurationSecs"]; exists {
		scalingMaxDurationSecs, err := voHelper.InterfaceToUint(
			input["scalingMaxDurationSecs"],
		)
		if err != nil {
			return NewServiceOutput(UserError, "InvalidScalingMaxDurationSecs")
		}
		scalingMaxDurationSecsPtr = &scalingMaxDurationSecs
	}

	var scalingIntervalSecsPtr *uint
	if _, exists := input["scalingIntervalSecs"]; exists {
		scalingIntervalSecs, err := voHelper.InterfaceToUint(input["scalingIntervalSecs"])
		if err != nil {
			return NewServiceOutput(UserError, "InvalidScalingIntervalSecs")
		}
		scalingIntervalSecsPtr = &scalingIntervalSecs
	}

	var hostMinCapacityPercentPtr *valueObject.HostMinCapacity
	if _, exists := input["hostMinCapacityPercent"]; exists {
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
	requiredParams := []string{"id"}

	err := serviceHelper.RequiredParamsInspector(input, requiredParams)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	profileId, err := valueObject.NewContainerProfileId(input["id"])
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
