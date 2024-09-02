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
) *ContainerProfileController {
	return &ContainerProfileController{
		containerProfileService: service.NewContainerProfileService(persistentDbSvc),
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
	var nameStr string
	var baseSpecsStr string
	var maxSpecsStr string
	var scalingPolicyStr string
	var scalingThreshold uint
	var scalingMaxDurationSecs uint
	var scalingIntervalSecs uint
	var hostMinCapacityPercent uint8

	cmd := &cobra.Command{
		Use:   "create",
		Short: "CreateNewContainerProfile",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{
				"name": nameStr,
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
		&scalingPolicyStr, "policy", "p", "", "ScalingPolicy (cpu|memory|connection)",
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
	var profileIdUint uint64
	var nameStr string
	var baseSpecsStr string
	var maxSpecsStr string
	var scalingPolicyStr string
	var scalingThreshold uint
	var scalingMaxDurationSecs uint
	var scalingIntervalSecs uint
	var hostMinCapacityPercent uint8

	cmd := &cobra.Command{
		Use:   "update",
		Short: "UpdateContainerProfile",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{
				"id": profileIdUint,
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
		&scalingPolicyStr, "policy", "p", "", "ScalingPolicy (cpu|memory|connection)",
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
	var profileIdUint uint64

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "DeleteContainerProfile",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{
				"profileId": profileIdUint,
			}

			cliHelper.ServiceResponseWrapper(
				controller.containerProfileService.Delete(requestBody),
			)
		},
	}

	cmd.Flags().Uint64VarP(&profileIdUint, "id", "i", 0, "ContainerProfileId")
	cmd.MarkFlagRequired("id")
	return cmd
}
