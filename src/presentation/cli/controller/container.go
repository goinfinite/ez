package cliController

import (
	"log/slog"

	"github.com/goinfinite/ez/src/domain/valueObject"
	"github.com/goinfinite/ez/src/infra/db"
	infraHelper "github.com/goinfinite/ez/src/infra/helper"
	cliHelper "github.com/goinfinite/ez/src/presentation/cli/helper"
	"github.com/goinfinite/ez/src/presentation/service"
	"github.com/spf13/cobra"
)

type ContainerController struct {
	containerService *service.ContainerService
}

func NewContainerController(
	persistentDbSvc *db.PersistentDatabaseService,
	trailDbSvc *db.TrailDatabaseService,
) *ContainerController {
	return &ContainerController{
		containerService: service.NewContainerService(
			persistentDbSvc, trailDbSvc,
		),
	}
}

func (controller *ContainerController) Read() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "ReadContainers",
		Run: func(cmd *cobra.Command, args []string) {
			cliHelper.ServiceResponseWrapper(controller.containerService.Read())
		},
	}

	return cmd
}

func (controller *ContainerController) ReadWithMetrics() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-with-metrics",
		Short: "ReadContainersWithMetrics",
		Run: func(cmd *cobra.Command, args []string) {
			cliHelper.ServiceResponseWrapper(
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
	for _, rawPortBinding := range portBindingsSlice {
		portBinding, err := valueObject.NewPortBindingFromString(rawPortBinding)
		if err != nil {
			slog.Debug(
				"ParsePortBindingsError",
				slog.String("rawPortBinding", rawPortBinding),
				slog.Any("error", err),
			)
			continue
		}

		portBindings = append(portBindings, portBinding...)
	}

	return portBindings
}

func (controller *ContainerController) parseContainerEnvs(
	envsSlice []string,
) []valueObject.ContainerEnv {
	envs := []valueObject.ContainerEnv{}
	for _, rawEnv := range envsSlice {
		env, err := valueObject.NewContainerEnv(rawEnv)
		if err != nil {
			slog.Debug(
				"ParseContainerEnvsError",
				slog.String("rawEnv", rawEnv), slog.Any("error", err),
			)
			continue
		}
		envs = append(envs, env)
	}

	return envs
}

func (controller *ContainerController) Create() *cobra.Command {
	var accountIdUint uint64
	var hostnameStr, containerImageAddressStr string
	var portBindingsSlice []string
	var restartPolicyStr, entrypointStr string
	var profileId uint64
	var envsSlice []string
	var launchScriptFilePathStr, autoCreateMappingsBoolStr string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "CreateNewContainer",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{
				"accountId":          accountIdUint,
				"hostname":           hostnameStr,
				"autoCreateMappings": autoCreateMappingsBoolStr,
			}

			if containerImageAddressStr != "" {
				requestBody["imageAddress"] = containerImageAddressStr
			}

			if len(portBindingsSlice) > 0 {
				portBindings := controller.parsePortBindings(portBindingsSlice)
				requestBody["portBindings"] = portBindings
			}

			if restartPolicyStr != "" {
				requestBody["restartPolicy"] = restartPolicyStr
			}

			if entrypointStr != "" {
				requestBody["entrypoint"] = entrypointStr
			}

			if profileId != 0 {
				requestBody["profileId"] = profileId
			}

			if len(envsSlice) > 0 {
				envs := controller.parseContainerEnvs(envsSlice)
				requestBody["envs"] = envs
			}

			if launchScriptFilePathStr != "" {
				scriptFilePath, err := valueObject.NewUnixFilePath(launchScriptFilePathStr)
				if err != nil {
					cliHelper.ResponseWrapper(false, err.Error())
				}

				scriptFileContent, err := infraHelper.GetFileContent(scriptFilePath.String())
				if err != nil {
					cliHelper.ResponseWrapper(false, err.Error())
				}

				launchScript, err := valueObject.NewLaunchScript(scriptFileContent)
				if err != nil {
					cliHelper.ResponseWrapper(false, err.Error())
				}

				err = infraHelper.RemoveFile(scriptFilePath.String())
				if err != nil {
					cliHelper.ResponseWrapper(false, err.Error())
				}

				requestBody["launchScript"] = launchScript
			}

			cliHelper.ServiceResponseWrapper(
				controller.containerService.Create(requestBody, false),
			)
		},
	}

	cmd.Flags().Uint64VarP(&accountIdUint, "account-id", "a", 0, "AccountId")
	cmd.MarkFlagRequired("account-id")
	cmd.Flags().StringVarP(&hostnameStr, "hostname", "n", "", "Hostname")
	cmd.MarkFlagRequired("hostname")
	cmd.Flags().StringVarP(&containerImageAddressStr, "image-address", "i", "", "ImageAddress")
	cmd.Flags().StringSliceVarP(
		&portBindingsSlice, "port-bindings", "b", []string{},
		"PortBindings (serviceName[:publicPort][:containerPort][/protocol][:privatePort])",
	)
	cmd.Flags().StringVarP(&restartPolicyStr, "restart-policy", "r", "", "RestartPolicy")
	cmd.Flags().StringVarP(&entrypointStr, "entrypoint", "e", "", "Entrypoint")
	cmd.Flags().Uint64VarP(&profileId, "profile-id", "p", 0, "ContainerProfileId")
	cmd.Flags().StringSliceVarP(&envsSlice, "envs", "v", []string{}, "Envs (key=value)")
	cmd.Flags().StringVarP(
		&launchScriptFilePathStr, "launch-script-path", "l", "", "Launch script file path",
	)
	cmd.Flags().StringVarP(
		&autoCreateMappingsBoolStr, "auto-create-mappings", "m", "true", "AutoCreateMappings",
	)
	return cmd
}

func (controller *ContainerController) Update() *cobra.Command {
	var accountIdUint uint64
	var containerIdStr string
	var containerStatusStr string
	var profileId uint64

	cmd := &cobra.Command{
		Use:   "update",
		Short: "UpdateContainer",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{
				"accountId":   accountIdUint,
				"containerId": containerIdStr,
			}

			if containerStatusStr != "" {
				requestBody["status"] = containerStatusStr
			}

			if profileId != 0 {
				requestBody["profileId"] = profileId
			}

			cliHelper.ServiceResponseWrapper(
				controller.containerService.Update(requestBody),
			)
		},
	}

	cmd.Flags().Uint64VarP(&accountIdUint, "account-id", "a", 0, "AccountId")
	cmd.MarkFlagRequired("account-id")
	cmd.Flags().StringVarP(&containerIdStr, "container-id", "c", "", "ContainerId")
	cmd.MarkFlagRequired("container-id")
	cmd.Flags().StringVarP(&containerStatusStr, "status", "s", "", "Status (true/false)")
	cmd.Flags().Uint64VarP(&profileId, "profile-id", "p", 0, "ContainerProfileId")
	return cmd
}

func (controller *ContainerController) Delete() *cobra.Command {
	var accountIdUint uint64
	var containerIdStr string

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "DeleteContainer",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{
				"accountId":   accountIdUint,
				"containerId": containerIdStr,
			}

			cliHelper.ServiceResponseWrapper(
				controller.containerService.Delete(requestBody),
			)
		},
	}

	cmd.Flags().Uint64VarP(&accountIdUint, "account-id", "a", 0, "AccountId")
	cmd.MarkFlagRequired("account-id")
	cmd.Flags().StringVarP(&containerIdStr, "container-id", "c", "", "ContainerId")
	cmd.MarkFlagRequired("container-id")
	return cmd
}
