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
	var containerIdStr string
	var containerAccountIdUint uint64
	var containerHostnameStr, containerStatusStr string
	var containerImageIdStr, containerImageAddressStr, containerImageHashStr string
	var portBindingsSlice []string
	var restartPolicyStr string
	var profileIdUint uint64
	var envsSlice []string
	var createdBeforeAtInt64, createdAfterAtInt64 int64
	var startedBeforeAtInt64, startedAfterAtInt64 int64
	var stoppedBeforeAtInt64, stoppedAfterAtInt64 int64
	var withMetricsBoolStr string
	var paginationPageNumberUint32 uint32
	var paginationItemsPerPageUint16 uint16
	var paginationSortByStr, paginationSortDirectionStr string
	var paginationLastSeenIdStr string

	cmd := &cobra.Command{
		Use:   "get",
		Short: "ReadContainers",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{}

			if containerIdStr != "" {
				requestBody["containerId"] = containerIdStr
			}
			if containerAccountIdUint != 0 {
				requestBody["containerAccountId"] = containerAccountIdUint
			}
			if containerHostnameStr != "" {
				requestBody["containerHostname"] = containerHostnameStr
			}
			if containerStatusStr != "" {
				requestBody["containerStatus"] = containerStatusStr
			}
			if containerImageIdStr != "" {
				requestBody["containerImageId"] = containerImageIdStr
			}
			if containerImageAddressStr != "" {
				requestBody["containerImageAddress"] = containerImageAddressStr
			}
			if containerImageHashStr != "" {
				requestBody["containerImageHash"] = containerImageHashStr
			}
			if len(portBindingsSlice) > 0 {
				portBindings := controller.parsePortBindings(portBindingsSlice)
				requestBody["containerPortBindings"] = portBindings
			}
			if restartPolicyStr != "" {
				requestBody["containerRestartPolicy"] = restartPolicyStr
			}
			if profileIdUint != 0 {
				requestBody["containerProfileId"] = profileIdUint
			}
			if len(envsSlice) > 0 {
				envs := controller.parseContainerEnvs(envsSlice)
				requestBody["containerEnv"] = envs
			}
			if createdBeforeAtInt64 != 0 {
				requestBody["createdBeforeAt"] = createdBeforeAtInt64
			}
			if createdAfterAtInt64 != 0 {
				requestBody["createdAfterAt"] = createdAfterAtInt64
			}
			if startedBeforeAtInt64 != 0 {
				requestBody["startedBeforeAt"] = startedBeforeAtInt64
			}
			if startedAfterAtInt64 != 0 {
				requestBody["startedAfterAt"] = startedAfterAtInt64
			}
			if stoppedBeforeAtInt64 != 0 {
				requestBody["stoppedBeforeAt"] = stoppedBeforeAtInt64
			}
			if stoppedAfterAtInt64 != 0 {
				requestBody["stoppedAfterAt"] = stoppedAfterAtInt64
			}
			if withMetricsBoolStr != "" {
				requestBody["withMetrics"] = withMetricsBoolStr
			}
			if paginationPageNumberUint32 != 0 {
				requestBody["pageNumber"] = paginationPageNumberUint32
			}
			if paginationItemsPerPageUint16 != 0 {
				requestBody["itemsPerPage"] = paginationItemsPerPageUint16
			}
			if paginationSortByStr != "" {
				requestBody["sortBy"] = paginationSortByStr
			}
			if paginationSortDirectionStr != "" {
				requestBody["sortDirection"] = paginationSortDirectionStr
			}
			if paginationLastSeenIdStr != "" {
				requestBody["lastSeenId"] = paginationLastSeenIdStr
			}

			cliHelper.ServiceResponseWrapper(
				controller.containerService.Read(requestBody),
			)
		},
	}

	cmd.Flags().StringVarP(&containerIdStr, "container-id", "c", "", "ContainerId")
	cmd.Flags().Uint64VarP(
		&containerAccountIdUint, "account-id", "a", 0, "ContainerAccountId",
	)
	cmd.Flags().StringVarP(
		&containerHostnameStr, "hostname", "n", "", "ContainerHostname",
	)
	cmd.Flags().StringVarP(&containerStatusStr, "status", "s", "", "ContainerStatus")
	cmd.Flags().StringVarP(&containerImageIdStr, "image-id", "i", "", "ContainerImageId")
	cmd.Flags().StringVarP(
		&containerImageAddressStr, "image-address", "d", "", "ContainerImageAddress",
	)
	cmd.Flags().StringVarP(&containerImageHashStr, "image-hash", "h", "", "ContainerImageHash")
	cmd.Flags().StringSliceVarP(
		&portBindingsSlice, "port-bindings", "b", []string{},
		"ContainerPortBindings (serviceName[:publicPort][:containerPort][/protocol][:privatePort])",
	)
	cmd.Flags().StringVarP(&restartPolicyStr, "restart-policy", "r", "", "ContainerRestartPolicy")
	cmd.Flags().Uint64VarP(&profileIdUint, "profile-id", "p", 0, "ContainerProfileId")
	cmd.Flags().StringSliceVarP(&envsSlice, "envs", "v", []string{}, "ContainerEnvs (key=value)")
	cmd.Flags().Int64VarP(
		&createdBeforeAtInt64, "created-before-at", "e", 0, "CreatedBeforeAt",
	)
	cmd.Flags().Int64VarP(
		&createdAfterAtInt64, "created-after-at", "d", 0, "CreatedAfterAt",
	)
	cmd.Flags().Int64VarP(
		&startedBeforeAtInt64, "started-before-at", "b", 0, "StartedBeforeAt",
	)
	cmd.Flags().Int64VarP(
		&startedAfterAtInt64, "started-after-at", "a", 0, "StartedAfterAt",
	)
	cmd.Flags().Int64VarP(
		&stoppedBeforeAtInt64, "stopped-before-at", "t", 0, "StoppedBeforeAt",
	)
	cmd.Flags().Int64VarP(
		&stoppedAfterAtInt64, "stopped-after-at", "o", 0, "StoppedAfterAt",
	)
	cmd.Flags().StringVarP(&withMetricsBoolStr, "with-metrics", "w", "", "WithMetrics")
	cmd.Flags().Uint32VarP(
		&paginationPageNumberUint32, "page-number", "p", 0, "PageNumber (Pagination)",
	)
	cmd.Flags().Uint16VarP(
		&paginationItemsPerPageUint16, "items-per-page", "m", 0, "ItemsPerPage (Pagination)",
	)
	cmd.Flags().StringVarP(
		&paginationSortByStr, "sort-by", "y", "", "SortBy (Pagination)",
	)
	cmd.Flags().StringVarP(
		&paginationSortDirectionStr, "sort-direction", "r", "", "SortDirection (Pagination)",
	)
	cmd.Flags().StringVarP(
		&paginationLastSeenIdStr, "last-seen-id", "l", "", "LastSeenId (Pagination)",
	)

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
	var hostnameStr, containerImageAddressStr, containerImageIdStr string
	var portBindingsSlice []string
	var restartPolicyStr, entrypointStr string
	var profileId uint64
	var envsSlice []string
	var launchScriptFilePathStr, autoCreateMappingsBoolStr, useImageExposedPortsBoolStr string
	var existingContainerIdStr string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "CreateNewContainer",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{
				"accountId":            accountIdUint,
				"hostname":             hostnameStr,
				"autoCreateMappings":   autoCreateMappingsBoolStr,
				"useImageExposedPorts": useImageExposedPortsBoolStr,
			}

			if containerImageAddressStr != "" {
				requestBody["imageAddress"] = containerImageAddressStr
			}

			if containerImageIdStr != "" {
				requestBody["imageId"] = containerImageIdStr
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

				scriptFileContent, err := infraHelper.ReadFileContent(scriptFilePath.String())
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

			if existingContainerIdStr != "" {
				requestBody["existingContainerId"] = existingContainerIdStr
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
	cmd.Flags().StringVarP(&containerImageIdStr, "image-id", "d", "", "ImageId")
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
	cmd.Flags().StringVarP(
		&useImageExposedPortsBoolStr, "use-image-exposed-ports", "u", "false",
		"UseImageExposedPorts (valid when imageId is provided)",
	)
	cmd.Flags().StringVarP(
		&existingContainerIdStr, "existing-container-id", "x", "", "ExistingContainerId",
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
