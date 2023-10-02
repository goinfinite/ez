package cliController

import (
	"strings"

	"github.com/speedianet/sfm/src/domain/dto"
	"github.com/speedianet/sfm/src/domain/useCase"
	"github.com/speedianet/sfm/src/domain/valueObject"
	"github.com/speedianet/sfm/src/infra"
	cliHelper "github.com/speedianet/sfm/src/presentation/cli/helper"
	"github.com/spf13/cobra"
)

func GetContainersController() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "GetContainers",
		Run: func(cmd *cobra.Command, args []string) {
			containerQueryRepo := infra.ContainerQueryRepo{}
			containersList, err := useCase.GetContainers(containerQueryRepo)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, containersList)
		},
	}

	return cmd
}

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
	var accId uint64
	var hostnameStr string
	var containerImgAddressStr string
	var portBindingsSlice []string
	var restartPolicyStr string
	var entrypointStr string
	var baseSpecStr string
	var maxSpecStr string
	var envsSlice []string

	cmd := &cobra.Command{
		Use:   "add",
		Short: "AddNewContainer",
		Run: func(cmd *cobra.Command, args []string) {
			accId := valueObject.NewAccountIdPanic(accId)
			hostname := valueObject.NewFqdnPanic(hostnameStr)
			imgAddr := valueObject.NewContainerImgAddressPanic(
				containerImgAddressStr,
			)

			var portBindingsPtr *[]valueObject.PortBinding
			if len(portBindingsSlice) > 0 {
				portBindings := parsePortBindings(portBindingsSlice)
				portBindingsPtr = &portBindings
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

			var baseSpecsPtr *valueObject.ContainerSpecs
			if baseSpecStr != "" {
				baseSpecs, err := valueObject.NewContainerSpecsFromString(baseSpecStr)
				if err != nil {
					cliHelper.ResponseWrapper(false, err.Error())
				}
				baseSpecsPtr = &baseSpecs
			}

			var maxSpecsPtr *valueObject.ContainerSpecs
			if maxSpecStr != "" {
				maxSpecs, err := valueObject.NewContainerSpecsFromString(maxSpecStr)
				if err != nil {
					cliHelper.ResponseWrapper(false, err.Error())
				}
				maxSpecsPtr = &maxSpecs
			}

			var envsPtr *[]valueObject.ContainerEnv
			if len(envsSlice) > 0 {
				envs := parseContainerEnvs(envsSlice)
				envsPtr = &envs
			}

			addContainerDto := dto.NewAddContainer(
				accId,
				hostname,
				imgAddr,
				portBindingsPtr,
				restartPolicyPtr,
				entrypointPtr,
				baseSpecsPtr,
				maxSpecsPtr,
				envsPtr,
			)

			containerCmdRepo := infra.ContainerCmdRepo{}
			accQueryRepo := infra.AccQueryRepo{}
			accCmdRepo := infra.AccCmdRepo{}

			err := useCase.AddContainer(
				containerCmdRepo,
				accQueryRepo,
				accCmdRepo,
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
		"p",
		[]string{},
		"PortBindings (hostPort:containerPort[/protocol])",
	)
	cmd.Flags().StringVarP(&restartPolicyStr, "restart-policy", "r", "", "RestartPolicy")
	cmd.Flags().StringVarP(&entrypointStr, "entrypoint", "e", "", "Entrypoint")
	cmd.Flags().StringVarP(
		&baseSpecStr,
		"base-specs",
		"b",
		"",
		"BaseSpecs (cpuCoresFloat:memoryBytesUint)",
	)
	cmd.Flags().StringVarP(
		&maxSpecStr,
		"max-specs",
		"m",
		"",
		"MaxSpecs (cpuCoresFloat:memoryBytesUint)",
	)
	cmd.Flags().StringSliceVarP(&envsSlice, "envs", "v", []string{}, "Envs (key=value)")
	return cmd
}

func DeleteContainerController() *cobra.Command {
	var accId uint64
	var containerIdStr string

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "DeleteContainer",
		Run: func(cmd *cobra.Command, args []string) {
			accId := valueObject.NewAccountIdPanic(accId)
			containerId := valueObject.NewContainerIdPanic(containerIdStr)

			containerQueryRepo := infra.ContainerQueryRepo{}
			containerCmdRepo := infra.ContainerCmdRepo{}
			accCmdRepo := infra.AccCmdRepo{}

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

func UpdateContainerController() *cobra.Command {
	var accId uint64
	var containerIdStr string
	var containerStatus bool
	var baseSpecStr string
	var maxSpecStr string

	cmd := &cobra.Command{
		Use:   "update",
		Short: "UpdateContainer",
		Run: func(cmd *cobra.Command, args []string) {
			accId := valueObject.NewAccountIdPanic(accId)
			containerId := valueObject.NewContainerIdPanic(containerIdStr)

			var baseSpecsPtr *valueObject.ContainerSpecs
			if baseSpecStr != "" {
				baseSpecs, err := valueObject.NewContainerSpecsFromString(baseSpecStr)
				if err != nil {
					cliHelper.ResponseWrapper(false, err.Error())
				}
				baseSpecsPtr = &baseSpecs
			}

			var maxSpecsPtr *valueObject.ContainerSpecs
			if maxSpecStr != "" {
				maxSpecs, err := valueObject.NewContainerSpecsFromString(maxSpecStr)
				if err != nil {
					cliHelper.ResponseWrapper(false, err.Error())
				}
				maxSpecsPtr = &maxSpecs
			}

			updateContainerDto := dto.NewUpdateContainer(
				accId,
				containerId,
				containerStatus,
				baseSpecsPtr,
				maxSpecsPtr,
			)

			containerQueryRepo := infra.ContainerQueryRepo{}
			containerCmdRepo := infra.ContainerCmdRepo{}
			accQueryRepo := infra.AccQueryRepo{}
			accCmdRepo := infra.AccCmdRepo{}

			err := useCase.UpdateContainer(
				containerQueryRepo,
				containerCmdRepo,
				accQueryRepo,
				accCmdRepo,
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
	cmd.Flags().BoolVarP(&containerStatus, "status", "s", false, "ContainerStatus")
	cmd.Flags().StringVarP(
		&baseSpecStr,
		"base-specs",
		"b",
		"",
		"BaseSpecs (cpuCoresFloat:memoryBytesUint)",
	)
	cmd.Flags().StringVarP(
		&maxSpecStr,
		"max-specs",
		"m",
		"",
		"MaxSpecs (cpuCoresFloat:memoryBytesUint)",
	)
	return cmd
}
