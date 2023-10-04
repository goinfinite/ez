package cliController

import (
	"strconv"
	"strings"

	"github.com/speedianet/sfm/src/domain/dto"
	"github.com/speedianet/sfm/src/domain/useCase"
	"github.com/speedianet/sfm/src/domain/valueObject"
	"github.com/speedianet/sfm/src/infra"
	cliHelper "github.com/speedianet/sfm/src/presentation/cli/helper"
	"github.com/spf13/cobra"
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
	var scalingThresholdStr string

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
			if scalingThresholdStr != "" {
				scalingThreshold, err := strconv.ParseUint(scalingThresholdStr, 10, 64)
				if err != nil {
					panic("InvalidScalingThreshold")
				}
				scalingThresholdPtr = &scalingThreshold
			}

			dto := dto.NewAddResourceProfile(
				name,
				baseSpecs,
				maxSpecsPtr,
				scalingPolicyPtr,
				scalingThresholdPtr,
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
		"scaling-policy",
		"s",
		"",
		"ScalingPolicy",
	)
	cmd.Flags().StringVarP(
		&scalingThresholdStr,
		"scaling-threshold",
		"t",
		"",
		"ScalingThreshold",
	)
	return cmd
}
