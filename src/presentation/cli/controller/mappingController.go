package cliController

import (
	"strings"

	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/useCase"
	"github.com/speedianet/control/src/domain/valueObject"
	"github.com/speedianet/control/src/infra"
	"github.com/speedianet/control/src/infra/db"
	cliHelper "github.com/speedianet/control/src/presentation/cli/helper"
	cliMiddleware "github.com/speedianet/control/src/presentation/cli/middleware"
	"github.com/spf13/cobra"
)

func GetMappingsController() *cobra.Command {
	var dbSvc *db.DatabaseService

	cmd := &cobra.Command{
		Use:   "get",
		Short: "GetMappings",
		PreRun: func(cmd *cobra.Command, args []string) {
			dbSvc = cliMiddleware.DatabaseInit()
		},
		Run: func(cmd *cobra.Command, args []string) {
			mappingQueryRepo := infra.NewMappingQueryRepo(dbSvc)
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
	var dbSvc *db.DatabaseService

	var accId uint64
	var hostnameStr string
	var hostPortUint uint64
	var hostProtocolStr string
	var targetsSlice []string

	cmd := &cobra.Command{
		Use:   "add",
		Short: "AddMapping",
		PreRun: func(cmd *cobra.Command, args []string) {
			dbSvc = cliMiddleware.DatabaseInit()
		},
		Run: func(cmd *cobra.Command, args []string) {
			accId := valueObject.NewAccountIdPanic(accId)
			var hostnamePtr *valueObject.Fqdn
			if hostnameStr != "" {
				hostname := valueObject.NewFqdnPanic(hostnameStr)
				hostnamePtr = &hostname
			}

			hostPort := valueObject.NewNetworkPortPanic(hostPortUint)

			if hostPort == 443 {
				hostProtocolStr = "https"
			}
			hostProtocol := valueObject.NewNetworkProtocolPanic(hostProtocolStr)

			mappingTargets := parseMappingTargets(
				targetsSlice,
				hostPort,
				hostProtocol,
			)

			addMappingDto := dto.NewAddMapping(
				accId,
				hostnamePtr,
				hostPort,
				hostProtocol,
				mappingTargets,
			)

			mappingQueryRepo := infra.NewMappingQueryRepo(dbSvc)
			mappingCmdRepo := infra.NewMappingCmdRepo(dbSvc)

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
	cmd.Flags().Uint64VarP(&accId, "acc-id", "a", 0, "AccountId")
	cmd.MarkFlagRequired("acc-id")
	cmd.Flags().StringVarP(&hostnameStr, "hostname", "n", "", "Hostname")
	cmd.Flags().Uint64VarP(&hostPortUint, "port", "p", 0, "Host Port")
	cmd.MarkFlagRequired("port")
	cmd.Flags().StringVarP(&hostProtocolStr, "protocol", "l", "tcp", "Host Protocol")
	cmd.Flags().StringSliceVarP(
		&targetsSlice,
		"target",
		"t",
		[]string{},
		"ContainerId (required), Port and Protocol (format: containerId:containerPort/protocol)",
	)
	cmd.MarkFlagRequired("target")
	return cmd
}

func DeleteMappingController() *cobra.Command {
	var dbSvc *db.DatabaseService

	var mappingIdUint uint64

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "DeleteMapping",
		PreRun: func(cmd *cobra.Command, args []string) {
			dbSvc = cliMiddleware.DatabaseInit()
		},
		Run: func(cmd *cobra.Command, args []string) {
			mappingId := valueObject.NewMappingIdPanic(mappingIdUint)

			mappingQueryRepo := infra.NewMappingQueryRepo(dbSvc)
			mappingCmdRepo := infra.NewMappingCmdRepo(dbSvc)

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
