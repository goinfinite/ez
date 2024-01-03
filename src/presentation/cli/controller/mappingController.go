package cliController

import (
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

func AddMappingController() *cobra.Command {
	var dbSvc *db.DatabaseService

	var accIdUint uint64
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
			accId := valueObject.NewAccountIdPanic(accIdUint)
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

			mappingTargets := []dto.AddMappingTargetWithoutMappingId{}
			for _, targetStr := range targetsSlice {
				target, err := dto.NewAddMappingTargetWithoutMappingIdFromString(targetStr)
				if err != nil {
					cliHelper.ResponseWrapper(false, err.Error())
				}
				mappingTargets = append(mappingTargets, target)
			}

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
	cmd.Flags().Uint64VarP(&accIdUint, "acc-id", "a", 0, "AccountId")
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

func AddMappingTargetController() *cobra.Command {
	var dbSvc *db.DatabaseService

	var mappingIdUint uint64
	var targetStr string

	cmd := &cobra.Command{
		Use:   "add",
		Short: "AddMappingTarget",
		PreRun: func(cmd *cobra.Command, args []string) {
			dbSvc = cliMiddleware.DatabaseInit()
		},
		Run: func(cmd *cobra.Command, args []string) {
			mappingId := valueObject.NewMappingIdPanic(mappingIdUint)
			target, err := dto.NewAddMappingTargetWithoutMappingIdFromString(targetStr)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			addTargetDto := dto.NewAddMappingTarget(
				mappingId,
				target.ContainerId,
				target.Port,
				target.Protocol,
			)

			mappingQueryRepo := infra.NewMappingQueryRepo(dbSvc)
			mappingCmdRepo := infra.NewMappingCmdRepo(dbSvc)

			err = useCase.AddMappingTarget(
				mappingQueryRepo,
				mappingCmdRepo,
				addTargetDto,
			)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, "MappingTargetAdded")
		},
	}

	cmd.Flags().Uint64VarP(&mappingIdUint, "mapping-id", "m", 0, "MappingId")
	cmd.MarkFlagRequired("mapping-id")
	cmd.Flags().StringVarP(
		&targetStr,
		"target",
		"t",
		"",
		"ContainerId (required), Port and Protocol (format: containerId:containerPort/protocol)",
	)
	cmd.MarkFlagRequired("target")
	return cmd
}

func DeleteMappingTargetController() *cobra.Command {
	var dbSvc *db.DatabaseService

	var targetIdUint uint64

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "DeleteMappingTarget",
		PreRun: func(cmd *cobra.Command, args []string) {
			dbSvc = cliMiddleware.DatabaseInit()
		},
		Run: func(cmd *cobra.Command, args []string) {
			targetId := valueObject.NewMappingTargetIdPanic(targetIdUint)

			mappingQueryRepo := infra.NewMappingQueryRepo(dbSvc)
			mappingCmdRepo := infra.NewMappingCmdRepo(dbSvc)

			err := useCase.DeleteMappingTarget(
				mappingQueryRepo,
				mappingCmdRepo,
				targetId,
			)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, "MappingTargetDeleted")
		},
	}

	cmd.Flags().Uint64VarP(&targetIdUint, "id", "i", 0, "MappingTargetId")
	cmd.MarkFlagRequired("id")
	return cmd
}
