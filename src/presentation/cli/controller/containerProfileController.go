package cliController

import (
	"strconv"
	"strings"

	"github.com/goinfinite/fleet/src/domain/dto"
	"github.com/goinfinite/fleet/src/domain/useCase"
	"github.com/goinfinite/fleet/src/domain/valueObject"
	"github.com/goinfinite/fleet/src/infra"
	"github.com/goinfinite/fleet/src/infra/db"
	cliHelper "github.com/goinfinite/fleet/src/presentation/cli/helper"
	cliMiddleware "github.com/goinfinite/fleet/src/presentation/cli/middleware"
	"github.com/spf13/cobra"
)

func GetContainerProfilesController() *cobra.Command {
	var dbSvc *db.DatabaseService

	cmd := &cobra.Command{
		Use:   "get",
		Short: "GetContainerProfiles",
		PreRun: func(cmd *cobra.Command, args []string) {
			dbSvc = cliMiddleware.DatabaseInit()
		},
		Run: func(cmd *cobra.Command, args []string) {
			containerProfileQueryRepo := infra.NewContainerProfileQueryRepo(dbSvc)
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

func AddContainerProfileController() *cobra.Command {
	var dbSvc *db.DatabaseService

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
		PreRun: func(cmd *cobra.Command, args []string) {
			dbSvc = cliMiddleware.DatabaseInit()
		},
		Run: func(cmd *cobra.Command, args []string) {
			name := valueObject.NewContainerProfileNamePanic(nameStr)

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

			containerProfileCmdRepo := infra.NewContainerProfileCmdRepo(dbSvc)

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

func UpdateContainerProfileController() *cobra.Command {
	var dbSvc *db.DatabaseService

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
		PreRun: func(cmd *cobra.Command, args []string) {
			dbSvc = cliMiddleware.DatabaseInit()
		},
		Run: func(cmd *cobra.Command, args []string) {
			profileId := valueObject.NewContainerProfileIdPanic(profileIdUint)

			var namePtr *valueObject.ContainerProfileName
			if nameStr != "" {
				name := valueObject.NewContainerProfileNamePanic(nameStr)
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

			containerProfileQueryRepo := infra.NewContainerProfileQueryRepo(dbSvc)
			containerProfileCmdRepo := infra.NewContainerProfileCmdRepo(dbSvc)
			containerQueryRepo := infra.NewContainerQueryRepo(dbSvc)
			containerCmdRepo := infra.NewContainerCmdRepo(dbSvc)

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

func DeleteContainerProfileController() *cobra.Command {
	var dbSvc *db.DatabaseService

	var profileIdUint uint64

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "DeleteContainerProfile",
		PreRun: func(cmd *cobra.Command, args []string) {
			dbSvc = cliMiddleware.DatabaseInit()
		},
		Run: func(cmd *cobra.Command, args []string) {
			profileId := valueObject.NewContainerProfileIdPanic(profileIdUint)

			containerProfileQueryRepo := infra.NewContainerProfileQueryRepo(dbSvc)
			containerProfileCmdRepo := infra.NewContainerProfileCmdRepo(dbSvc)
			containerQueryRepo := infra.NewContainerQueryRepo(dbSvc)
			containerCmdRepo := infra.NewContainerCmdRepo(dbSvc)

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
