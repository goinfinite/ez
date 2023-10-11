package cliController

import (
	"strconv"
	"strings"

	"github.com/speedianet/sfm/src/domain/dto"
	"github.com/speedianet/sfm/src/domain/useCase"
	"github.com/speedianet/sfm/src/domain/valueObject"
	"github.com/speedianet/sfm/src/infra"
	cliHelper "github.com/speedianet/sfm/src/presentation/cli/helper"
	cliMiddleware "github.com/speedianet/sfm/src/presentation/cli/middleware"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

func GetResourceProfilesController() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "GetResourceProfiles",
		Run: func(cmd *cobra.Command, args []string) {
			resourceProfileQueryRepo := infra.ResourceProfileQueryRepo{}
			resourceProfilesList, err := useCase.GetResourceProfiles(
				resourceProfileQueryRepo,
			)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, resourceProfilesList)
		},
	}

	return cmd
}

func parseContainerSpecs(specStr string) valueObject.ContainerSpecs {
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

func AddResourceProfileController() *cobra.Command {
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
		Short: "AddNewResourceProfile",
		Run: func(cmd *cobra.Command, args []string) {
			name := valueObject.NewResourceProfileNamePanic(nameStr)

			baseSpecs := parseContainerSpecs(baseSpecsStr)

			var maxSpecsPtr *valueObject.ContainerSpecs
			if maxSpecsStr != "" {
				maxSpecs := parseContainerSpecs(maxSpecsStr)
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

			dto := dto.NewAddResourceProfile(
				name,
				baseSpecs,
				maxSpecsPtr,
				scalingPolicyPtr,
				scalingThresholdPtr,
				scalingMaxDurationSecsPtr,
				scalingIntervalSecsPtr,
				hostMinCapacityPercentPtr,
			)

			resourceProfileCmdRepo := infra.ResourceProfileCmdRepo{}

			err := useCase.AddResourceProfile(
				resourceProfileCmdRepo,
				dto,
			)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, "ResourceProfileAdded")
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
		"i",
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

func UpdateResourceProfileController() *cobra.Command {
	var dbSvc *gorm.DB

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
		Short: "UpdateResourceProfile",
		PreRun: func(cmd *cobra.Command, args []string) {
			dbSvc = cliMiddleware.DatabaseInit()
		},
		Run: func(cmd *cobra.Command, args []string) {
			resourceProfileId := valueObject.NewResourceProfileIdPanic(profileIdUint)

			var namePtr *valueObject.ResourceProfileName
			if nameStr != "" {
				name := valueObject.NewResourceProfileNamePanic(nameStr)
				namePtr = &name
			}

			var baseSpecsPtr *valueObject.ContainerSpecs
			if baseSpecsStr != "" {
				baseSpecs := parseContainerSpecs(baseSpecsStr)
				baseSpecsPtr = &baseSpecs
			}

			var maxSpecsPtr *valueObject.ContainerSpecs
			if maxSpecsStr != "" {
				maxSpecs := parseContainerSpecs(maxSpecsStr)
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

			dto := dto.NewUpdateResourceProfile(
				resourceProfileId,
				namePtr,
				baseSpecsPtr,
				maxSpecsPtr,
				scalingPolicyPtr,
				scalingThresholdPtr,
				scalingMaxDurationSecsPtr,
				scalingIntervalSecsPtr,
				hostMinCapacityPercentPtr,
			)

			resourceProfileQueryRepo := infra.ResourceProfileQueryRepo{}
			resourceProfileCmdRepo := infra.ResourceProfileCmdRepo{}
			containerQueryRepo := infra.NewContainerQueryRepo(dbSvc)
			containerCmdRepo := infra.ContainerCmdRepo{}

			err := useCase.UpdateResourceProfile(
				resourceProfileQueryRepo,
				resourceProfileCmdRepo,
				containerQueryRepo,
				containerCmdRepo,
				dto,
			)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, "ResourceProfileUpdated")
		},
	}

	cmd.Flags().Uint64VarP(&profileIdUint, "id", "i", 0, "ResourceProfileId")
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
		"i",
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

func DeleteResourceProfileController() *cobra.Command {
	var dbSvc *gorm.DB

	var profileIdUint uint64

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "DeleteResourceProfile",
		PreRun: func(cmd *cobra.Command, args []string) {
			dbSvc = cliMiddleware.DatabaseInit()
		},
		Run: func(cmd *cobra.Command, args []string) {
			resourceProfileId := valueObject.NewResourceProfileIdPanic(profileIdUint)

			resourceProfileQueryRepo := infra.ResourceProfileQueryRepo{}
			resourceProfileCmdRepo := infra.ResourceProfileCmdRepo{}
			containerQueryRepo := infra.NewContainerQueryRepo(dbSvc)
			containerCmdRepo := infra.ContainerCmdRepo{}

			err := useCase.DeleteResourceProfile(
				resourceProfileQueryRepo,
				resourceProfileCmdRepo,
				containerQueryRepo,
				containerCmdRepo,
				resourceProfileId,
			)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, "ResourceProfileDeleted")
		},
	}

	cmd.Flags().Uint64VarP(&profileIdUint, "id", "i", 0, "ResourceProfileId")
	cmd.MarkFlagRequired("id")
	return cmd
}
