package cliController

import (
	"strconv"

	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/useCase"
	"github.com/speedianet/control/src/domain/valueObject"
	"github.com/speedianet/control/src/infra"
	"github.com/speedianet/control/src/infra/db"
	cliHelper "github.com/speedianet/control/src/presentation/cli/helper"
	"github.com/speedianet/control/src/presentation/service"
	"github.com/spf13/cobra"
)

type ContainerController struct {
	persistentDbSvc  *db.PersistentDatabaseService
	containerService *service.ContainerService
}

func NewContainerController(
	persistentDbSvc *db.PersistentDatabaseService,
) *ContainerController {
	return &ContainerController{
		persistentDbSvc:  persistentDbSvc,
		containerService: service.NewContainerService(persistentDbSvc),
	}
}

func (controller *ContainerController) Read() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "ReadContainers",
		Run: func(cmd *cobra.Command, args []string) {
			cliHelper.NewResponseWrapper(controller.containerService.Read())
		},
	}

	return cmd
}

func (controller *ContainerController) ReadWithMetrics() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-with-metrics",
		Short: "ReadContainersWithMetrics",
		Run: func(cmd *cobra.Command, args []string) {
			cliHelper.NewResponseWrapper(
				controller.containerService.ReadWithMetrics(),
			)
		},
	}

	return cmd
}

func (controller *ContainerController) parsePortBindings(
	portBindingsSlice []string,
) []valueObject.PortBinding {
	portBindings := []valueObject.PortBinding{}
	for _, portBindingStr := range portBindingsSlice {
		portBinding, err := valueObject.NewPortBindingFromString(portBindingStr)
		if err != nil {
			panic(err.Error() + ": " + portBindingStr)
		}

		portBindings = append(portBindings, portBinding...)
	}

	return portBindings
}

func (controller *ContainerController) parseContainerEnvs(
	envsSlice []string,
) []valueObject.ContainerEnv {
	envs := []valueObject.ContainerEnv{}
	for _, envStr := range envsSlice {
		env := valueObject.NewContainerEnvPanic(envStr)
		envs = append(envs, env)
	}

	return envs
}

func (controller *ContainerController) Create() *cobra.Command {
	var accountIdUint uint64
	var hostnameStr string
	var containerImageAddressStr string
	var portBindingsSlice []string
	var restartPolicyStr string
	var entrypointStr string
	var profileId uint64
	var envsSlice []string
	var autoCreateMappings bool

	cmd := &cobra.Command{
		Use:   "create",
		Short: "CreateNewContainer",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{
				"accountId":          accountIdUint,
				"hostname":           hostnameStr,
				"imageAddress":       containerImageAddressStr,
				"restartPolicy":      restartPolicyStr,
				"entrypoint":         entrypointStr,
				"profileId":          profileId,
				"autoCreateMappings": autoCreateMappings,
			}

			portBindings := []valueObject.PortBinding{}
			if len(portBindingsSlice) > 0 {
				portBindings = controller.parsePortBindings(portBindingsSlice)
			}

			envs := []valueObject.ContainerEnv{}
			if len(envsSlice) > 0 {
				envs = controller.parseContainerEnvs(envsSlice)
			}

			requestBody["portBindings"] = portBindings
			requestBody["envs"] = envs

			cliHelper.NewResponseWrapper(
				controller.containerService.Create(requestBody),
			)
		},
	}

	cmd.Flags().Uint64VarP(&accountIdUint, "account-id", "a", 0, "AccountId")
	cmd.MarkFlagRequired("account-id")
	cmd.Flags().StringVarP(&hostnameStr, "hostname", "n", "", "Hostname")
	cmd.MarkFlagRequired("hostname")
	cmd.Flags().StringVarP(&containerImageAddressStr, "image", "i", "", "ImageAddress")
	cmd.MarkFlagRequired("image")
	cmd.Flags().StringSliceVarP(
		&portBindingsSlice,
		"port-bindings",
		"b",
		[]string{},
		"PortBindings (serviceName[:publicPort][:containerPort][/protocol][:privatePort])",
	)
	cmd.Flags().StringVarP(&restartPolicyStr, "restart-policy", "r", "", "RestartPolicy")
	cmd.Flags().StringVarP(&entrypointStr, "entrypoint", "e", "", "Entrypoint")
	cmd.Flags().Uint64VarP(&profileId, "profile-id", "p", 0, "ContainerProfileId")
	cmd.Flags().StringSliceVarP(&envsSlice, "envs", "v", []string{}, "Envs (key=value)")
	cmd.Flags().BoolVarP(
		&autoCreateMappings,
		"auto-create-mappings",
		"m",
		true,
		"AutoCreateMappings",
	)
	return cmd
}

func (controller *ContainerController) Update() *cobra.Command {
	var accId uint64
	var containerIdStr string
	var containerStatusStr string
	var profileId uint64

	cmd := &cobra.Command{
		Use:   "update",
		Short: "UpdateContainer",
		Run: func(cmd *cobra.Command, args []string) {
			accId := valueObject.NewAccountIdPanic(accId)
			containerId := valueObject.NewContainerIdPanic(containerIdStr)

			var containerStatusPtr *bool
			if containerStatusStr != "" {
				containerStatus, err := strconv.ParseBool(containerStatusStr)
				if err != nil {
					panic("InvalidContainerStatus")
				}
				containerStatusPtr = &containerStatus
			}

			var profileIdPtr *valueObject.ContainerProfileId
			if profileId != 0 {
				profileId := valueObject.NewContainerProfileIdPanic(
					profileId,
				)
				profileIdPtr = &profileId
			}

			updateContainerDto := dto.NewUpdateContainer(
				accId,
				containerId,
				containerStatusPtr,
				profileIdPtr,
			)

			containerQueryRepo := infra.NewContainerQueryRepo(controller.persistentDbSvc)
			containerCmdRepo := infra.NewContainerCmdRepo(controller.persistentDbSvc)
			accQueryRepo := infra.NewAccQueryRepo(controller.persistentDbSvc)
			accCmdRepo := infra.NewAccCmdRepo(controller.persistentDbSvc)
			containerProfileQueryRepo := infra.NewContainerProfileQueryRepo(controller.persistentDbSvc)

			err := useCase.UpdateContainer(
				containerQueryRepo,
				containerCmdRepo,
				accQueryRepo,
				accCmdRepo,
				containerProfileQueryRepo,
				updateContainerDto,
			)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, "ContainerUpdated")
		},
	}

	cmd.Flags().Uint64VarP(&accId, "acc-id", "a", 0, "AccountId")
	cmd.MarkFlagRequired("acc-id")
	cmd.Flags().StringVarP(&containerIdStr, "container-id", "c", "", "ContainerId")
	cmd.MarkFlagRequired("container-id")
	cmd.Flags().StringVarP(&containerStatusStr, "status", "s", "", "Status (true/false)")
	cmd.Flags().Uint64VarP(&profileId, "profile-id", "p", 0, "ContainerProfileId")
	return cmd
}

func (controller *ContainerController) Delete() *cobra.Command {
	var accId uint64
	var containerIdStr string

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "DeleteContainer",
		Run: func(cmd *cobra.Command, args []string) {
			accId := valueObject.NewAccountIdPanic(accId)
			containerId := valueObject.NewContainerIdPanic(containerIdStr)

			containerQueryRepo := infra.NewContainerQueryRepo(controller.persistentDbSvc)
			containerCmdRepo := infra.NewContainerCmdRepo(controller.persistentDbSvc)
			accCmdRepo := infra.NewAccCmdRepo(controller.persistentDbSvc)
			mappingQueryRepo := infra.NewMappingQueryRepo(controller.persistentDbSvc)
			mappingCmdRepo := infra.NewMappingCmdRepo(controller.persistentDbSvc)

			err := useCase.DeleteContainer(
				containerQueryRepo,
				containerCmdRepo,
				accCmdRepo,
				mappingQueryRepo,
				mappingCmdRepo,
				accId,
				containerId,
			)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, "ContainerDeleted")
		},
	}

	cmd.Flags().Uint64VarP(&accId, "acc-id", "a", 0, "AccountId")
	cmd.MarkFlagRequired("acc-id")
	cmd.Flags().StringVarP(&containerIdStr, "container-id", "c", "", "ContainerId")
	cmd.MarkFlagRequired("container-id")
	return cmd
}
