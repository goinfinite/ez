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

func parsePortBindings(portBindingsSlice []string) []valueObject.PortBinding {
	portBindings := []valueObject.PortBinding{}
	for _, portBindingStr := range portBindingsSlice {
		protocol := valueObject.NewNetworkProtocolPanic("tcp")

		portBindingParts := strings.Split(portBindingStr, ":")
		hostPort, err := valueObject.NewNetworkPort(portBindingParts[0])
		if err != nil {
			panic("InvalidPortBindingHostPort")
		}
		containerPortStr := portBindingParts[1]

		containerPortParts := strings.Split(containerPortStr, "/")
		containerPort, err := valueObject.NewNetworkPort(containerPortParts[0])
		if err != nil {
			panic("InvalidPortBindingContainerPort")
		}
		if len(containerPortParts) == 1 {
			protocol = valueObject.NewNetworkProtocolPanic(containerPortParts[1])
		}

		portBinding := valueObject.NewPortBinding(protocol, hostPort, containerPort)
		portBindings = append(portBindings, portBinding)
	}

	return portBindings
}

func parseContainerEnvs(envsSlice []string) []valueObject.ContainerEnv {
	envs := []valueObject.ContainerEnv{}
	for _, envStr := range envsSlice {
		env := valueObject.NewContainerEnvPanic(envStr)
		envs = append(envs, env)
	}

	return envs
}

func AddContainerController() *cobra.Command {
	var dbSvc *db.DatabaseService

	var accId uint64
	var hostnameStr string
	var containerImgAddressStr string
	var portBindingsSlice []string
	var restartPolicyStr string
	var entrypointStr string
	var profileId uint64
	var envsSlice []string

	cmd := &cobra.Command{
		Use:   "add",
		Short: "AddNewContainer",
		PreRun: func(cmd *cobra.Command, args []string) {
			dbSvc = cliMiddleware.DatabaseInit()
		},
		Run: func(cmd *cobra.Command, args []string) {
			accId := valueObject.NewAccountIdPanic(accId)
			hostname := valueObject.NewFqdnPanic(hostnameStr)
			imgAddr := valueObject.NewContainerImgAddressPanic(
				containerImgAddressStr,
			)

			portBindings := []valueObject.PortBinding{}
			if len(portBindingsSlice) > 0 {
				portBindings = parsePortBindings(portBindingsSlice)
			}

			var restartPolicyPtr *valueObject.ContainerRestartPolicy
			if restartPolicyStr != "" {
				restartPolicy := valueObject.NewContainerRestartPolicyPanic(
					restartPolicyStr,
				)
				restartPolicyPtr = &restartPolicy
			}

			var entrypointPtr *valueObject.ContainerEntrypoint
			if entrypointStr != "" {
				entrypoint := valueObject.NewContainerEntrypointPanic(entrypointStr)
				entrypointPtr = &entrypoint
			}

			var profileIdPtr *valueObject.ContainerProfileId
			if profileId != 0 {
				profileId := valueObject.NewContainerProfileIdPanic(
					profileId,
				)
				profileIdPtr = &profileId
			}

			envs := []valueObject.ContainerEnv{}
			if len(envsSlice) > 0 {
				envs = parseContainerEnvs(envsSlice)
			}

			addContainerDto := dto.NewAddContainer(
				accId,
				hostname,
				imgAddr,
				portBindings,
				restartPolicyPtr,
				entrypointPtr,
				profileIdPtr,
				envs,
			)

			containerCmdRepo := infra.ContainerCmdRepo{}
			accQueryRepo := infra.NewAccQueryRepo(dbSvc)
			accCmdRepo := infra.NewAccCmdRepo(dbSvc)
			containerProfileQueryRepo := infra.NewContainerProfileQueryRepo(dbSvc)

			err := useCase.AddContainer(
				containerCmdRepo,
				accQueryRepo,
				accCmdRepo,
				containerProfileQueryRepo,
				addContainerDto,
			)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, "ContainerAdded")
		},
	}

	cmd.Flags().Uint64VarP(&accId, "acc-id", "a", 0, "AccountId")
	cmd.MarkFlagRequired("acc-id")
	cmd.Flags().StringVarP(&hostnameStr, "hostname", "n", "", "Hostname")
	cmd.MarkFlagRequired("hostname")
	cmd.Flags().StringVarP(&containerImgAddressStr, "image", "i", "", "ImageAddress")
	cmd.MarkFlagRequired("image")
	cmd.Flags().StringSliceVarP(
		&portBindingsSlice,
		"port-bindings",
		"b",
		[]string{},
		"PortBindings (hostPort:containerPort[/protocol])",
	)
	cmd.Flags().StringVarP(&restartPolicyStr, "restart-policy", "r", "", "RestartPolicy")
	cmd.Flags().StringVarP(&entrypointStr, "entrypoint", "e", "", "Entrypoint")
	cmd.Flags().Uint64VarP(&profileId, "profile-id", "p", 0, "ContainerProfileId")
	cmd.Flags().StringSliceVarP(&envsSlice, "envs", "v", []string{}, "Envs (key=value)")
	return cmd
}

func UpdateContainerController() *cobra.Command {
	var dbSvc *db.DatabaseService

	var accId uint64
	var containerIdStr string
	var containerStatusStr string
	var profileId uint64

	cmd := &cobra.Command{
		Use:   "update",
		Short: "UpdateContainer",
		PreRun: func(cmd *cobra.Command, args []string) {
			dbSvc = cliMiddleware.DatabaseInit()
		},
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

			containerQueryRepo := infra.NewContainerQueryRepo(dbSvc)
			containerCmdRepo := infra.ContainerCmdRepo{}
			accQueryRepo := infra.NewAccQueryRepo(dbSvc)
			accCmdRepo := infra.NewAccCmdRepo(dbSvc)
			containerProfileQueryRepo := infra.NewContainerProfileQueryRepo(dbSvc)

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
	cmd.Flags().StringVarP(&containerStatusStr, "status", "s", "", "ContainerStatus (true or false)")
	cmd.Flags().Uint64VarP(&profileId, "profile-id", "p", 0, "ContainerProfileId")
	return cmd
}

func DeleteContainerController() *cobra.Command {
	var dbSvc *db.DatabaseService

	var accId uint64
	var containerIdStr string

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "DeleteContainer",
		PreRun: func(cmd *cobra.Command, args []string) {
			dbSvc = cliMiddleware.DatabaseInit()
		},
		Run: func(cmd *cobra.Command, args []string) {
			accId := valueObject.NewAccountIdPanic(accId)
			containerId := valueObject.NewContainerIdPanic(containerIdStr)

			containerQueryRepo := infra.NewContainerQueryRepo(dbSvc)
			containerCmdRepo := infra.ContainerCmdRepo{}
			accCmdRepo := infra.NewAccCmdRepo(dbSvc)

			err := useCase.DeleteContainer(
				containerQueryRepo,
				containerCmdRepo,
				accCmdRepo,
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
