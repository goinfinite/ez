package cliController

import (
	"github.com/speedianet/control/src/domain/valueObject"
	"github.com/speedianet/control/src/infra/db"
	cliHelper "github.com/speedianet/control/src/presentation/cli/helper"
	"github.com/speedianet/control/src/presentation/service"
	"github.com/spf13/cobra"
)

type ContainerProfileController struct {
	containerProfileService *service.ContainerProfileService
}

func NewContainerProfileController(
	persistentDbSvc *db.PersistentDatabaseService,
	trailDbSvc *db.TrailDatabaseService,
) *ContainerProfileController {
	return &ContainerProfileController{
		containerProfileService: service.NewContainerProfileService(
			persistentDbSvc, trailDbSvc,
		),
	}
}

func (controller *ContainerProfileController) Read() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "ReadContainerProfiles",
		Run: func(cmd *cobra.Command, args []string) {
			cliHelper.ServiceResponseWrapper(controller.containerProfileService.Read())
		},
	}
	return cmd
}

func (controller *ContainerProfileController) Create() *cobra.Command {
	var accountIdUint uint64
	var nameStr, baseSpecsStr, maxSpecsStr, scalingPolicyStr string
	var scalingThreshold, scalingMaxDurationSecs, scalingIntervalSecs uint
	var hostMinCapacityPercent uint8

	cmd := &cobra.Command{
		Use:   "create",
		Short: "CreateNewContainerProfile",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{
				"accountId": accountIdUint,
				"name":      nameStr,
			}

			if baseSpecsStr != "" {
				baseSpecs, err := valueObject.NewContainerSpecsFromString(baseSpecsStr)
				if err != nil {
					cliHelper.ResponseWrapper(false, err.Error())
				}

				requestBody["baseSpecs"] = baseSpecs
			}

			if maxSpecsStr != "" {
				maxSpecs, err := valueObject.NewContainerSpecsFromString(maxSpecsStr)
				if err != nil {
					cliHelper.ResponseWrapper(false, err.Error())
				}

				requestBody["maxSpecs"] = maxSpecs
			}

			if scalingPolicyStr != "" {
				requestBody["scalingPolicy"] = scalingPolicyStr
			}

			if scalingThreshold != 0 {
				requestBody["scalingThreshold"] = scalingThreshold
			}

			if scalingMaxDurationSecs != 0 {
				requestBody["scalingMaxDurationSecs"] = scalingMaxDurationSecs
			}

			if scalingIntervalSecs != 0 {
				requestBody["scalingIntervalSecs"] = scalingIntervalSecs
			}

			if hostMinCapacityPercent != 0 {
				requestBody["hostMinCapacityPercent"] = hostMinCapacityPercent
			}

			cliHelper.ServiceResponseWrapper(
				controller.containerProfileService.Create(requestBody),
			)
		},
	}

	cmd.Flags().Uint64VarP(&accountIdUint, "account-id", "a", 0, "AccountId")
	cmd.MarkFlagRequired("account-id")
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
		&scalingPolicyStr, "policy", "s", "", "ScalingPolicy (cpu|memory|connection)",
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
	cmd.Flags().Uint8VarP(
		&hostMinCapacityPercent, "min-capacity", "c", 0, "HostMinCapacityPercent (0-100)",
	)
	return cmd
}

func (controller *ContainerProfileController) Update() *cobra.Command {
	var accountIdUint, profileIdUint uint64
	var nameStr, baseSpecsStr, maxSpecsStr, scalingPolicyStr string
	var scalingThreshold, scalingMaxDurationSecs, scalingIntervalSecs uint
	var hostMinCapacityPercent uint8

	cmd := &cobra.Command{
		Use:   "update",
		Short: "UpdateContainerProfile",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{
				"accountId": accountIdUint,
				"profileId": profileIdUint,
			}

			if nameStr != "" {
				requestBody["name"] = nameStr
			}

			if baseSpecsStr != "" {
				baseSpecs, err := valueObject.NewContainerSpecsFromString(baseSpecsStr)
				if err != nil {
					cliHelper.ResponseWrapper(false, err.Error())
				}

				requestBody["baseSpecs"] = baseSpecs
			}

			if maxSpecsStr != "" {
				maxSpecs, err := valueObject.NewContainerSpecsFromString(maxSpecsStr)
				if err != nil {
					cliHelper.ResponseWrapper(false, err.Error())
				}

				requestBody["maxSpecs"] = maxSpecs
			}

			if scalingPolicyStr != "" {
				requestBody["scalingPolicy"] = scalingPolicyStr
			}

			if scalingThreshold != 0 {
				requestBody["scalingThreshold"] = scalingThreshold
			}

			if scalingMaxDurationSecs != 0 {
				requestBody["scalingMaxDurationSecs"] = scalingMaxDurationSecs
			}

			if scalingIntervalSecs != 0 {
				requestBody["scalingIntervalSecs"] = scalingIntervalSecs
			}

			if hostMinCapacityPercent != 0 {
				requestBody["hostMinCapacityPercent"] = hostMinCapacityPercent
			}

			cliHelper.ServiceResponseWrapper(
				controller.containerProfileService.Update(requestBody),
			)
		},
	}

	cmd.Flags().Uint64VarP(&accountIdUint, "account-id", "a", 0, "AccountId")
	cmd.MarkFlagRequired("account-id")
	cmd.Flags().Uint64VarP(&profileIdUint, "profile-id", "p", 0, "ContainerProfileId")
	cmd.MarkFlagRequired("profile-id")
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
		&scalingPolicyStr, "policy", "s", "", "ScalingPolicy (cpu|memory|connection)",
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
	cmd.Flags().Uint8VarP(
		&hostMinCapacityPercent, "min-capacity", "c", 0, "HostMinCapacityPercent (0-100)",
	)
	return cmd
}

func (controller *ContainerProfileController) Delete() *cobra.Command {
	var accountIdUint, profileIdUint uint64

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "DeleteContainerProfile",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{
				"accountId": accountIdUint,
				"profileId": profileIdUint,
			}

			cliHelper.ServiceResponseWrapper(
				controller.containerProfileService.Delete(requestBody),
			)
		},
	}

	cmd.Flags().Uint64VarP(&accountIdUint, "account-id", "a", 0, "AccountId")
	cmd.MarkFlagRequired("account-id")
	cmd.Flags().Uint64VarP(&profileIdUint, "profile-id", "p", 0, "ContainerProfileId")
	cmd.MarkFlagRequired("profile-id")
	return cmd
}
