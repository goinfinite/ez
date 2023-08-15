package cliController

import (
	"errors"
	"strconv"
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

func parsePortBindings(portBindingsSlice []string) ([]valueObject.PortBinding, error) {
	portBindings := []valueObject.PortBinding{}
	for _, portBindingStr := range portBindingsSlice {
		protocol := valueObject.NewNetworkProtocolPanic("tcp")

		portBindingParts := strings.Split(portBindingStr, ":")
		hostPort, err := strconv.ParseUint(portBindingParts[0], 10, 16)
		if err != nil {
			return nil, err
		}
		containerPortStr := portBindingParts[1]

		containerPortParts := strings.Split(containerPortStr, "/")
		containerPort, err := strconv.ParseUint(containerPortParts[0], 10, 16)
		if err != nil {
			return nil, err
		}
		if len(containerPortParts) == 1 {
			protocol = valueObject.NewNetworkProtocolPanic(containerPortParts[1])
		}

		portBinding := valueObject.NewPortBinding(protocol, hostPort, containerPort)
		portBindings = append(portBindings, portBinding)
	}

	return portBindings, nil
}

func parseContainerSpecs(specStr string) (valueObject.ContainerSpecs, error) {
	specParts := strings.Split(specStr, ":")
	cpuCores, err := strconv.ParseFloat(specParts[0], 64)
	if err != nil {
		return valueObject.ContainerSpecs{}, errors.New("InvalidCpuCoresLimit")
	}

	memory, err := strconv.ParseInt(specParts[1], 10, 64)
	if err != nil {
		return valueObject.ContainerSpecs{}, errors.New("InvalidMemoryLimit")
	}

	storage, err := strconv.ParseInt(specParts[2], 10, 64)
	if err != nil {
		return valueObject.ContainerSpecs{}, errors.New("InvalidStorageLimit")
	}

	return valueObject.NewContainerSpecs(
		cpuCores,
		valueObject.Byte(memory),
		valueObject.Byte(storage),
	), nil
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
	var hostnameStr string
	var containerImageAddressStr string
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
			hostname := valueObject.NewFqdnPanic(hostnameStr)
			imgAddr := valueObject.NewContainerImgAddressPanic(containerImageAddressStr)

			portBindings, err := parsePortBindings(portBindingsSlice)
			if err != nil {
				cliHelper.ResponseWrapper(false, "InvalidPortBindingsFormat")
				return
			}

			restartPolicy := valueObject.NewContainerRestartPolicyPanic(restartPolicyStr)
			entrypoint := valueObject.NewContainerEntrypointPanic(entrypointStr)
			baseSpecs, err := parseContainerSpecs(baseSpecStr)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
				return
			}

			maxSpecs, err := parseContainerSpecs(maxSpecStr)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
				return
			}

			envs := parseContainerEnvs(envsSlice)

			addContainerDto := dto.NewAddContainer(containername, password)

			containerQueryRepo := infra.ContainerQueryRepo{}
			containerCmdRepo := infra.ContainerCmdRepo{}

			err := useCase.AddContainer(
				containerQueryRepo,
				containerCmdRepo,
				addContainerDto,
			)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, "ContainerAdded")
		},
	}

	cmd.Flags().StringVarP(&containernameStr, "containername", "u", "", "Containername")
	cmd.MarkFlagRequired("containername")
	cmd.Flags().StringVarP(&passwordStr, "password", "p", "", "Password")
	cmd.MarkFlagRequired("password")
	return cmd
}

func DeleteContainerController() *cobra.Command {
	var containerIdStr string

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "DeleteContainer",
		Run: func(cmd *cobra.Command, args []string) {
			containerId := valueObject.NewContainerIdFromStringPanic(containerIdStr)

			containerQueryRepo := infra.ContainerQueryRepo{}
			containerCmdRepo := infra.ContainerCmdRepo{}

			err := useCase.DeleteContainer(
				containerQueryRepo,
				containerCmdRepo,
				containerId,
			)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, "ContainerDeleted")
		},
	}

	cmd.Flags().StringVarP(&containerIdStr, "container-id", "u", "", "ContainerId")
	cmd.MarkFlagRequired("container-id")
	return cmd
}

func UpdateContainerController() *cobra.Command {
	var containerIdStr string
	var passwordStr string
	shouldUpdateApiKeyBool := false

	cmd := &cobra.Command{
		Use:   "update",
		Short: "UpdateContainer (pass or apiKey)",
		Run: func(cmd *cobra.Command, args []string) {
			containerId := valueObject.NewContainerIdFromStringPanic(containerIdStr)

			var passPtr *valueObject.Password
			if passwordStr != "" {
				password := valueObject.NewPasswordPanic(passwordStr)
				passPtr = &password
			}

			var shouldUpdateApiKeyPtr *bool
			if shouldUpdateApiKeyBool {
				shouldUpdateApiKeyPtr = &shouldUpdateApiKeyBool
			}

			updateContainerDto := dto.NewUpdateContainer(
				containerId,
				passPtr,
				shouldUpdateApiKeyPtr,
			)

			containerQueryRepo := infra.ContainerQueryRepo{}
			containerCmdRepo := infra.ContainerCmdRepo{}

			if updateContainerDto.Password != nil {
				useCase.UpdateContainerPassword(
					containerQueryRepo,
					containerCmdRepo,
					updateContainerDto,
				)
			}

			if shouldUpdateApiKeyBool {
				newKey, err := useCase.UpdateContainerApiKey(
					containerQueryRepo,
					containerCmdRepo,
					updateContainerDto,
				)
				if err != nil {
					cliHelper.ResponseWrapper(false, err.Error())
				}

				cliHelper.ResponseWrapper(true, newKey)
			}
		},
	}

	cmd.Flags().StringVarP(&containerIdStr, "container-id", "u", "", "ContainerId")
	cmd.MarkFlagRequired("container-id")
	cmd.Flags().StringVarP(&passwordStr, "password", "p", "", "Password")
	cmd.Flags().BoolVarP(
		&shouldUpdateApiKeyBool,
		"update-api-key",
		"k",
		false,
		"ShouldUpdateApiKey",
	)
	return cmd
}
