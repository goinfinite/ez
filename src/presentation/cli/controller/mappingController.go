package cliController

import (
	"strings"

	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/useCase"
	"github.com/speedianet/control/src/domain/valueObject"
	"github.com/speedianet/control/src/infra"
	cliHelper "github.com/speedianet/control/src/presentation/cli/helper"
	"github.com/spf13/cobra"
)

func GetMappingsController() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "GetMappings",
		Run: func(cmd *cobra.Command, args []string) {
			mappingQueryRepo := infra.MappingQueryRepo{}
			mappingsList, err := useCase.GetMappings(mappingQueryRepo)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, mappingsList)
		},
	}

	return cmd
}

func parsePortProtocol(portProtocol string) map[string]string {
	portProtocolParts := strings.Split(portProtocol, "/")
	hostPortStr := portProtocolParts[0]
	hostProtocolStr := "tcp"
	if len(portProtocolParts) == 2 {
		hostProtocolStr = portProtocolParts[1]
	}

	return map[string]string{
		"port":     hostPortStr,
		"protocol": hostProtocolStr,
	}
}

func parseMappingTargets(
	targetsStrSlice []string,
	hostPort valueObject.NetworkPort,
	hostProtocol valueObject.NetworkProtocol,
) []valueObject.MappingTarget {
	var mappingTargets []valueObject.MappingTarget
	for _, idPortProtocol := range targetsStrSlice {
		idPortProtocolParts := strings.Split(idPortProtocol, ":")
		containerId, err := valueObject.NewContainerId(idPortProtocolParts[0])
		if err != nil {
			continue
		}

		port := hostPort
		protocol := hostProtocol

		if len(idPortProtocolParts) > 1 {
			portProtocolMap := parsePortProtocol(idPortProtocolParts[1])
			port, err = valueObject.NewNetworkPort(portProtocolMap["port"])
			if err != nil {
				continue
			}

			protocol, err = valueObject.NewNetworkProtocol(portProtocolMap["protocol"])
			if err != nil {
				continue
			}
		}

		mappingTarget := valueObject.NewMappingTarget(
			containerId,
			port,
			protocol,
		)

		mappingTargets = append(
			mappingTargets,
			mappingTarget,
		)
	}

	return mappingTargets
}

func AddMappingController() *cobra.Command {
	var hostnameStr string
	var hostPortUint uint64
	var hostProtocolStr string
	var targetsSlice []string

	cmd := &cobra.Command{
		Use:   "add",
		Short: "AddMapping",
		Run: func(cmd *cobra.Command, args []string) {
			var hostnamePtr *valueObject.Fqdn
			if hostnameStr != "" {
				hostname := valueObject.NewFqdnPanic(hostnameStr)
				hostnamePtr = &hostname
			}

			hostPort := valueObject.NewNetworkPortPanic(hostPortUint)

			rawHostProtocol := "tcp"
			if hostPort == 443 {
				rawHostProtocol = "https"
			}

			if hostProtocolStr != "" {
				rawHostProtocol = hostProtocolStr
			}
			hostProtocol := valueObject.NewNetworkProtocolPanic(rawHostProtocol)

			mappingTargets := parseMappingTargets(
				targetsSlice,
				hostPort,
				hostProtocol,
			)

			addMappingDto := dto.NewAddMapping(
				hostnamePtr,
				hostPort,
				hostProtocol,
				mappingTargets,
			)

			mappingQueryRepo := infra.MappingQueryRepo{}
			mappingCmdRepo := infra.MappingCmdRepo{}

			err := useCase.AddMapping(
				mappingQueryRepo,
				mappingCmdRepo,
				addMappingDto,
			)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, "MappingAdded")
		},
	}
	cmd.Flags().StringVarP(&hostnameStr, "hostname", "n", "", "Hostname")
	cmd.Flags().Uint64VarP(&hostPortUint, "port", "p", 0, "Host Port")
	cmd.MarkFlagRequired("port")
	cmd.Flags().StringVarP(&hostProtocolStr, "protocol", "r", "tcp", "Host Protocol")
	cmd.Flags().StringSliceVarP(
		&targetsSlice,
		"target",
		"t",
		[]string{},
		"Container Id (required), Port and Protocol (format: containerId:containerPort/protocol)",
	)
	cmd.MarkFlagRequired("target")
	return cmd
}

func DeleteMappingController() *cobra.Command {
	var mappingIdUint uint64

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "DeleteMapping",
		Run: func(cmd *cobra.Command, args []string) {
			mappingId := valueObject.NewMappingIdPanic(mappingIdUint)

			mappingQueryRepo := infra.MappingQueryRepo{}
			mappingCmdRepo := infra.MappingCmdRepo{}

			err := useCase.DeleteMapping(
				mappingQueryRepo,
				mappingCmdRepo,
				mappingId,
			)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, "MappingDeleted")
		},
	}

	cmd.Flags().Uint64VarP(&mappingIdUint, "id", "i", 0, "MappingId")
	cmd.MarkFlagRequired("id")
	return cmd
}
