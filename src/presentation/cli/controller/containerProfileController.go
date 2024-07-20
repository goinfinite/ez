package cliController

import (
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

func NewContainerProfileController(
	persistentDbSvc *db.PersistentDatabaseService,
) *ContainerProfileController {
	return &ContainerProfileController{persistentDbSvc: persistentDbSvc}
}

func (controller *ContainerProfileController) ReadContainerProfiles() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "ReadContainerProfiles",
		Run: func(cmd *cobra.Command, args []string) {
			containerProfileQueryRepo := infra.NewContainerProfileQueryRepo(
				controller.persistentDbSvc,
			)
			containerProfilesList, err := useCase.ReadContainerProfiles(
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

func (controller *ContainerProfileController) CreateContainerProfile() *cobra.Command {
	var nameStr string
	var baseSpecsStr string
	var maxSpecsStr string
	var scalingPolicyStr string
	var scalingThreshold uint
	var scalingMaxDurationSecs uint
	var scalingIntervalSecs uint
	var hostMinCapacityPercent float64

	cmd := &cobra.Command{
		Use:   "create",
		Short: "CreateNewContainerProfile",
		Run: func(cmd *cobra.Command, args []string) {
			name := valueObject.NewContainerProfileNamePanic(nameStr)

			baseSpecs, err := valueObject.NewContainerSpecsFromString(baseSpecsStr)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			var maxSpecsPtr *valueObject.ContainerSpecs
			if maxSpecsStr != "" {
				maxSpecs, err := valueObject.NewContainerSpecsFromString(maxSpecsStr)
				if err != nil {
					cliHelper.ResponseWrapper(false, err.Error())
				}
				maxSpecsPtr = &maxSpecs
			}

			var scalingPolicyPtr *valueObject.ScalingPolicy
			if scalingPolicyStr != "" {
				scalingPolicy := valueObject.NewScalingPolicyPanic(scalingPolicyStr)
				scalingPolicyPtr = &scalingPolicy
			}

			var scalingThresholdPtr *uint
			if scalingThreshold != 0 {
				scalingThresholdPtr = &scalingThreshold
			}

			var scalingMaxDurationSecsPtr *uint
			if scalingMaxDurationSecs != 0 {
				scalingMaxDurationSecsPtr = &scalingMaxDurationSecs
			}

			var scalingIntervalSecsPtr *uint
			if scalingIntervalSecs != 0 {
				scalingIntervalSecsPtr = &scalingIntervalSecs
			}

			var hostMinCapacityPercentPtr *valueObject.HostMinCapacity
			if hostMinCapacityPercent != 0 {
				hostMinCapacityPercent, err := valueObject.NewHostMinCapacity(
					hostMinCapacityPercent,
				)
				if err != nil {
					cliHelper.ResponseWrapper(false, err.Error())
				}
				hostMinCapacityPercentPtr = &hostMinCapacityPercent
			}

			dto := dto.NewCreateContainerProfile(
				name, baseSpecs, maxSpecsPtr, scalingPolicyPtr, scalingThresholdPtr,
				scalingMaxDurationSecsPtr, scalingIntervalSecsPtr, hostMinCapacityPercentPtr,
			)

			containerProfileCmdRepo := infra.NewContainerProfileCmdRepo(
				controller.persistentDbSvc,
			)

			err = useCase.CreateContainerProfile(containerProfileCmdRepo, dto)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, "ContainerProfileAdded")
		},
	}

	cmd.Flags().StringVarP(&nameStr, "name", "n", "", "Name")
	cmd.MarkFlagRequired("name")
	cmd.Flags().StringVarP(
		&baseSpecsStr, "base-specs", "b", "",
		"BaseSpecs (millicores:memoryBytes:storagePerformanceUnits)",
	)
	cmd.MarkFlagRequired("base-specs")
	cmd.Flags().StringVarP(
		&maxSpecsStr, "max-specs", "m", "",
		"MaxSpecs (millicores:memoryBytes:storagePerformanceUnits)",
	)
	cmd.Flags().StringVarP(
		&scalingPolicyStr, "policy", "p", "", "ScalingPolicy",
	)
	cmd.Flags().UintVarP(
		&scalingThreshold, "threshold", "t", 0, "ScalingThreshold",
	)
	cmd.Flags().UintVarP(
		&scalingMaxDurationSecs, "max-duration", "d", 0, "ScalingMaxDurationSecs",
	)
	cmd.Flags().UintVarP(
		&scalingIntervalSecs, "interval", "v", 0, "ScalingIntervalSecs",
	)
	cmd.Flags().Float64VarP(
		&hostMinCapacityPercent, "min-capacity", "c", 0, "HostMinCapacityPercent",
	)
	return cmd
}

func (controller *ContainerProfileController) UpdateContainerProfile() *cobra.Command {
	var profileIdUint uint64
	var nameStr string
	var baseSpecsStr string
	var maxSpecsStr string
	var scalingPolicyStr string
	var scalingThreshold uint
	var scalingMaxDurationSecs uint
	var scalingIntervalSecs uint
	var hostMinCapacityPercent float64

	cmd := &cobra.Command{
		Use:   "update",
		Short: "UpdateContainerProfile",
		Run: func(cmd *cobra.Command, args []string) {
			profileId, err := valueObject.NewContainerProfileId(profileIdUint)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			var namePtr *valueObject.ContainerProfileName
			if nameStr != "" {
				name := valueObject.NewContainerProfileNamePanic(nameStr)
				namePtr = &name
			}

			var baseSpecsPtr *valueObject.ContainerSpecs
			if baseSpecsStr != "" {
				baseSpecs, err := valueObject.NewContainerSpecsFromString(baseSpecsStr)
				if err != nil {
					cliHelper.ResponseWrapper(false, err.Error())
				}
				baseSpecsPtr = &baseSpecs
			}

			var maxSpecsPtr *valueObject.ContainerSpecs
			if maxSpecsStr != "" {
				maxSpecs, err := valueObject.NewContainerSpecsFromString(maxSpecsStr)
				if err != nil {
					cliHelper.ResponseWrapper(false, err.Error())
				}
				maxSpecsPtr = &maxSpecs
			}

			var scalingPolicyPtr *valueObject.ScalingPolicy
			if scalingPolicyStr != "" {
				scalingPolicy := valueObject.NewScalingPolicyPanic(scalingPolicyStr)
				scalingPolicyPtr = &scalingPolicy
			}

			var scalingThresholdPtr *uint
			if scalingThreshold != 0 {
				scalingThresholdPtr = &scalingThreshold
			}

			var scalingMaxDurationSecsPtr *uint
			if scalingMaxDurationSecs != 0 {
				scalingMaxDurationSecsPtr = &scalingMaxDurationSecs
			}

			var scalingIntervalSecsPtr *uint
			if scalingIntervalSecs != 0 {
				scalingIntervalSecsPtr = &scalingIntervalSecs
			}

			var hostMinCapacityPercentPtr *valueObject.HostMinCapacity
			if hostMinCapacityPercent != 0 {
				hostMinCapacityPercent, err := valueObject.NewHostMinCapacity(
					hostMinCapacityPercent,
				)
				if err != nil {
					cliHelper.ResponseWrapper(false, err.Error())
				}
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

			containerProfileQueryRepo := infra.NewContainerProfileQueryRepo(
				controller.persistentDbSvc,
			)
			containerProfileCmdRepo := infra.NewContainerProfileCmdRepo(controller.persistentDbSvc)
			containerQueryRepo := infra.NewContainerQueryRepo(controller.persistentDbSvc)
			containerCmdRepo := infra.NewContainerCmdRepo(controller.persistentDbSvc)

			err = useCase.UpdateContainerProfile(
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
		&baseSpecsStr, "base-specs", "b", "",
		"BaseSpecs (millicores:memoryBytes:storagePerformanceUnits)",
	)
	cmd.Flags().StringVarP(
		&maxSpecsStr, "max-specs", "m", "",
		"MaxSpecs (millicores:memoryBytes:storagePerformanceUnits)",
	)
	cmd.Flags().StringVarP(
		&scalingPolicyStr, "policy", "p", "", "ScalingPolicy",
	)
	cmd.Flags().UintVarP(
		&scalingThreshold, "threshold", "t", 0, "ScalingThreshold",
	)
	cmd.Flags().UintVarP(
		&scalingMaxDurationSecs, "max-duration", "d", 0, "ScalingMaxDurationSecs",
	)
	cmd.Flags().UintVarP(
		&scalingIntervalSecs, "interval", "v", 0, "ScalingIntervalSecs",
	)
	cmd.Flags().Float64VarP(
		&hostMinCapacityPercent, "min-capacity", "c", 0, "HostMinCapacityPercent",
	)
	return cmd
}

func (controller *ContainerProfileController) DeleteContainerProfile() *cobra.Command {
	var profileIdUint uint64

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "DeleteContainerProfile",
		Run: func(cmd *cobra.Command, args []string) {
			profileId, err := valueObject.NewContainerProfileId(profileIdUint)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			containerProfileQueryRepo := infra.NewContainerProfileQueryRepo(controller.persistentDbSvc)
			containerProfileCmdRepo := infra.NewContainerProfileCmdRepo(controller.persistentDbSvc)
			containerQueryRepo := infra.NewContainerQueryRepo(controller.persistentDbSvc)
			containerCmdRepo := infra.NewContainerCmdRepo(controller.persistentDbSvc)

			err = useCase.DeleteContainerProfile(
				containerProfileQueryRepo, containerProfileCmdRepo, containerQueryRepo,
				containerCmdRepo, profileId,
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
