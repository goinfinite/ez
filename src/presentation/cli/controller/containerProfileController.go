package cliController

import (
	"strconv"
	"strings"

	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/useCase"
	"github.com/speedianet/control/src/domain/valueObject"
	"github.com/speedianet/control/src/infra"
	"github.com/speedianet/control/src/infra/db"
	cliHelper "github.com/speedianet/control/src/presentation/cli/helper"
	"github.com/spf13/cobra"
)

type ContainerProfileController struct {
	persistentDbSvc *db.PersistentDatabaseService
}

func NewContainerProfileController(persistentDbSvc *db.PersistentDatabaseService) ContainerProfileController {
	return ContainerProfileController{persistentDbSvc: persistentDbSvc}
}

func (controller ContainerProfileController) GetContainerProfiles() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "GetContainerProfiles",
		Run: func(cmd *cobra.Command, args []string) {
			containerProfileQueryRepo := infra.NewContainerProfileQueryRepo(
				controller.persistentDbSvc,
			)
			containerProfilesList, err := useCase.GetContainerProfiles(
				containerProfileQueryRepo,
			)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, containerProfilesList)
		},
	}

	return cmd
}

func (controller ContainerProfileController) parseContainerSpecs(
	specStr string,
) valueObject.ContainerSpecs {
	specParts := strings.Split(specStr, ":")
	if len(specParts) != 2 {
		panic("InvalidContainerSpecs")
	}

	cpuCores, err := valueObject.NewCpuCoresCount(specParts[0])
	if err != nil {
		panic("InvalidCpuCoresLimit")
	}

	memory, err := strconv.ParseInt(specParts[1], 10, 64)
	if err != nil {
		panic("InvalidMemoryLimit")
	}

	return valueObject.NewContainerSpecs(
		cpuCores,
		valueObject.Byte(memory),
	)
}

func (controller ContainerProfileController) AddContainerProfile() *cobra.Command {
	var nameStr string
	var baseSpecsStr string
	var maxSpecsStr string
	var scalingPolicyStr string
	var scalingThreshold uint64
	var scalingMaxDurationSecs uint64
	var scalingIntervalSecs uint64
	var hostMinCapacityPercent float64

	cmd := &cobra.Command{
		Use:   "add",
		Short: "AddNewContainerProfile",
		Run: func(cmd *cobra.Command, args []string) {
			name := valueObject.NewContainerProfileNamePanic(nameStr)

			baseSpecs := controller.parseContainerSpecs(baseSpecsStr)

			var maxSpecsPtr *valueObject.ContainerSpecs
			if maxSpecsStr != "" {
				maxSpecs := controller.parseContainerSpecs(maxSpecsStr)
				maxSpecsPtr = &maxSpecs
			}

			var scalingPolicyPtr *valueObject.ScalingPolicy
			if scalingPolicyStr != "" {
				scalingPolicy := valueObject.NewScalingPolicyPanic(scalingPolicyStr)
				scalingPolicyPtr = &scalingPolicy
			}

			var scalingThresholdPtr *uint64
			if scalingThreshold != 0 {
				scalingThresholdPtr = &scalingThreshold
			}

			var scalingMaxDurationSecsPtr *uint64
			if scalingMaxDurationSecs != 0 {
				scalingMaxDurationSecsPtr = &scalingMaxDurationSecs
			}

			var scalingIntervalSecsPtr *uint64
			if scalingIntervalSecs != 0 {
				scalingIntervalSecsPtr = &scalingIntervalSecs
			}

			var hostMinCapacityPercentPtr *valueObject.HostMinCapacity
			if hostMinCapacityPercent != 0 {
				hostMinCapacityPercent := valueObject.NewHostMinCapacityPanic(hostMinCapacityPercent)
				hostMinCapacityPercentPtr = &hostMinCapacityPercent
			}

			dto := dto.NewAddContainerProfile(
				name,
				baseSpecs,
				maxSpecsPtr,
				scalingPolicyPtr,
				scalingThresholdPtr,
				scalingMaxDurationSecsPtr,
				scalingIntervalSecsPtr,
				hostMinCapacityPercentPtr,
			)

			containerProfileCmdRepo := infra.NewContainerProfileCmdRepo(controller.persistentDbSvc)

			err := useCase.AddContainerProfile(
				containerProfileCmdRepo,
				dto,
			)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, "ContainerProfileAdded")
		},
	}

	cmd.Flags().StringVarP(&nameStr, "name", "n", "", "Name")
	cmd.MarkFlagRequired("name")
	cmd.Flags().StringVarP(
		&baseSpecsStr,
		"base-specs",
		"b",
		"",
		"BaseSpecs (cpuCoresFloat:memoryBytesUint)",
	)
	cmd.MarkFlagRequired("base-specs")
	cmd.Flags().StringVarP(
		&maxSpecsStr,
		"max-specs",
		"m",
		"",
		"MaxSpecs (cpuCoresFloat:memoryBytesUint)",
	)
	cmd.Flags().StringVarP(
		&scalingPolicyStr,
		"policy",
		"p",
		"",
		"ScalingPolicy",
	)
	cmd.Flags().Uint64VarP(
		&scalingThreshold,
		"threshold",
		"t",
		0,
		"ScalingThreshold",
	)
	cmd.Flags().Uint64VarP(
		&scalingMaxDurationSecs,
		"max-duration",
		"d",
		0,
		"ScalingMaxDurationSecs",
	)
	cmd.Flags().Uint64VarP(
		&scalingIntervalSecs,
		"interval",
		"v",
		0,
		"ScalingIntervalSecs",
	)
	cmd.Flags().Float64VarP(
		&hostMinCapacityPercent,
		"min-capacity",
		"c",
		0,
		"HostMinCapacityPercent",
	)
	return cmd
}

func (controller ContainerProfileController) UpdateContainerProfile() *cobra.Command {
	var profileIdUint uint64
	var nameStr string
	var baseSpecsStr string
	var maxSpecsStr string
	var scalingPolicyStr string
	var scalingThreshold uint64
	var scalingMaxDurationSecs uint64
	var scalingIntervalSecs uint64
	var hostMinCapacityPercent float64

	cmd := &cobra.Command{
		Use:   "update",
		Short: "UpdateContainerProfile",
		Run: func(cmd *cobra.Command, args []string) {
			profileId := valueObject.NewContainerProfileIdPanic(profileIdUint)

			var namePtr *valueObject.ContainerProfileName
			if nameStr != "" {
				name := valueObject.NewContainerProfileNamePanic(nameStr)
				namePtr = &name
			}

			var baseSpecsPtr *valueObject.ContainerSpecs
			if baseSpecsStr != "" {
				baseSpecs := controller.parseContainerSpecs(baseSpecsStr)
				baseSpecsPtr = &baseSpecs
			}

			var maxSpecsPtr *valueObject.ContainerSpecs
			if maxSpecsStr != "" {
				maxSpecs := controller.parseContainerSpecs(maxSpecsStr)
				maxSpecsPtr = &maxSpecs
			}

			var scalingPolicyPtr *valueObject.ScalingPolicy
			if scalingPolicyStr != "" {
				scalingPolicy := valueObject.NewScalingPolicyPanic(scalingPolicyStr)
				scalingPolicyPtr = &scalingPolicy
			}

			var scalingThresholdPtr *uint64
			if scalingThreshold != 0 {
				scalingThresholdPtr = &scalingThreshold
			}

			var scalingMaxDurationSecsPtr *uint64
			if scalingMaxDurationSecs != 0 {
				scalingMaxDurationSecsPtr = &scalingMaxDurationSecs
			}

			var scalingIntervalSecsPtr *uint64
			if scalingIntervalSecs != 0 {
				scalingIntervalSecsPtr = &scalingIntervalSecs
			}

			var hostMinCapacityPercentPtr *valueObject.HostMinCapacity
			if hostMinCapacityPercent != 0 {
				hostMinCapacityPercent := valueObject.NewHostMinCapacityPanic(hostMinCapacityPercent)
				hostMinCapacityPercentPtr = &hostMinCapacityPercent
			}

			dto := dto.NewUpdateContainerProfile(
				profileId,
				namePtr,
				baseSpecsPtr,
				maxSpecsPtr,
				scalingPolicyPtr,
				scalingThresholdPtr,
				scalingMaxDurationSecsPtr,
				scalingIntervalSecsPtr,
				hostMinCapacityPercentPtr,
			)

			containerProfileQueryRepo := infra.NewContainerProfileQueryRepo(controller.persistentDbSvc)
			containerProfileCmdRepo := infra.NewContainerProfileCmdRepo(controller.persistentDbSvc)
			containerQueryRepo := infra.NewContainerQueryRepo(controller.persistentDbSvc)
			containerCmdRepo := infra.NewContainerCmdRepo(controller.persistentDbSvc)

			err := useCase.UpdateContainerProfile(
				containerProfileQueryRepo,
				containerProfileCmdRepo,
				containerQueryRepo,
				containerCmdRepo,
				dto,
			)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, "ContainerProfileUpdated")
		},
	}

	cmd.Flags().Uint64VarP(&profileIdUint, "id", "i", 0, "ContainerProfileId")
	cmd.MarkFlagRequired("id")
	cmd.Flags().StringVarP(&nameStr, "name", "n", "", "Name")
	cmd.Flags().StringVarP(
		&baseSpecsStr,
		"base-specs",
		"b",
		"",
		"BaseSpecs (cpuCoresFloat:memoryBytesUint)",
	)
	cmd.Flags().StringVarP(
		&maxSpecsStr,
		"max-specs",
		"m",
		"",
		"MaxSpecs (cpuCoresFloat:memoryBytesUint)",
	)
	cmd.Flags().StringVarP(
		&scalingPolicyStr,
		"policy",
		"p",
		"",
		"ScalingPolicy",
	)
	cmd.Flags().Uint64VarP(
		&scalingThreshold,
		"threshold",
		"t",
		0,
		"ScalingThreshold",
	)
	cmd.Flags().Uint64VarP(
		&scalingMaxDurationSecs,
		"max-duration",
		"d",
		0,
		"ScalingMaxDurationSecs",
	)
	cmd.Flags().Uint64VarP(
		&scalingIntervalSecs,
		"interval",
		"v",
		0,
		"ScalingIntervalSecs",
	)
	cmd.Flags().Float64VarP(
		&hostMinCapacityPercent,
		"min-capacity",
		"c",
		0,
		"HostMinCapacityPercent",
	)
	return cmd
}

func (controller ContainerProfileController) DeleteContainerProfile() *cobra.Command {
	var profileIdUint uint64

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "DeleteContainerProfile",
		Run: func(cmd *cobra.Command, args []string) {
			profileId := valueObject.NewContainerProfileIdPanic(profileIdUint)

			containerProfileQueryRepo := infra.NewContainerProfileQueryRepo(controller.persistentDbSvc)
			containerProfileCmdRepo := infra.NewContainerProfileCmdRepo(controller.persistentDbSvc)
			containerQueryRepo := infra.NewContainerQueryRepo(controller.persistentDbSvc)
			containerCmdRepo := infra.NewContainerCmdRepo(controller.persistentDbSvc)

			err := useCase.DeleteContainerProfile(
				containerProfileQueryRepo,
				containerProfileCmdRepo,
				containerQueryRepo,
				containerCmdRepo,
				profileId,
			)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, "ContainerProfileDeleted")
		},
	}

	cmd.Flags().Uint64VarP(&profileIdUint, "id", "i", 0, "ContainerProfileId")
	cmd.MarkFlagRequired("id")
	return cmd
}
