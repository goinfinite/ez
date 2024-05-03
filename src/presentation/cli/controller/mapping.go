package cliController

import (
	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/useCase"
	"github.com/speedianet/control/src/domain/valueObject"
	"github.com/speedianet/control/src/infra"
	"github.com/speedianet/control/src/infra/db"
	cliHelper "github.com/speedianet/control/src/presentation/cli/helper"
	"github.com/speedianet/control/src/presentation/service"
	"github.com/spf13/cobra"
)

type MappingController struct {
	mappingService *service.MappingService
}

func NewMappingController(
	persistentDbSvc *db.PersistentDatabaseService,
) *MappingController {
	return &MappingController{
		mappingService: service.NewMappingService(persistentDbSvc),
	}
}

func (controller *MappingController) Read() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "ReadMappings",
		Run: func(cmd *cobra.Command, args []string) {
			mappingQueryRepo := infra.NewMappingQueryRepo(controller.persistentDbSvc)
			mappingsList, err := useCase.GetMappings(mappingQueryRepo)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, mappingsList)
		},
	}

	return cmd
}

func (controller *MappingController) Create() *cobra.Command {
	var accIdUint uint64
	var hostnameStr string
	var publicPortUint uint64
	var hostProtocolStr string
	var targetsSlice []string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "CreateMapping",
		Run: func(cmd *cobra.Command, args []string) {
			accId := valueObject.NewAccountIdPanic(accIdUint)
			var hostnamePtr *valueObject.Fqdn
			if hostnameStr != "" {
				hostname := valueObject.NewFqdnPanic(hostnameStr)
				hostnamePtr = &hostname
			}

			publicPort := valueObject.NewNetworkPortPanic(publicPortUint)

			protocolStr := valueObject.GuessNetworkProtocolByPort(publicPort).String()
			if hostProtocolStr != "" {
				protocolStr = hostProtocolStr
			}
			hostProtocol := valueObject.NewNetworkProtocolPanic(protocolStr)

			targets := []valueObject.ContainerId{}
			for _, targetStr := range targetsSlice {
				containerId, err := valueObject.NewContainerId(targetStr)
				if err != nil {
					cliHelper.ResponseWrapper(false, err.Error())
				}
				targets = append(targets, containerId)
			}
			createMappingDto := dto.NewCreateMapping(
				accId,
				hostnamePtr,
				publicPort,
				hostProtocol,
				targets,
			)

			mappingQueryRepo := infra.NewMappingQueryRepo(controller.persistentDbSvc)
			mappingCmdRepo := infra.NewMappingCmdRepo(controller.persistentDbSvc)
			containerQueryRepo := infra.NewContainerQueryRepo(controller.persistentDbSvc)

			err := useCase.CreateMapping(
				mappingQueryRepo,
				mappingCmdRepo,
				containerQueryRepo,
				createMappingDto,
			)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, "MappingCreated")
		},
	}
	cmd.Flags().Uint64VarP(&accIdUint, "acc-id", "a", 0, "AccountId")
	cmd.MarkFlagRequired("acc-id")
	cmd.Flags().StringVarP(&hostnameStr, "hostname", "n", "", "Hostname")
	cmd.Flags().Uint64VarP(&publicPortUint, "port", "p", 0, "Public Port")
	cmd.MarkFlagRequired("port")
	cmd.Flags().StringVarP(&hostProtocolStr, "protocol", "l", "", "Host Protocol")
	cmd.Flags().StringSliceVarP(
		&targetsSlice,
		"target",
		"t",
		[]string{},
		"ContainerIds",
	)
	cmd.MarkFlagRequired("target")
	return cmd
}

func (controller *MappingController) Delete() *cobra.Command {
	var mappingIdUint uint64

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "DeleteMapping",
		Run: func(cmd *cobra.Command, args []string) {
			mappingId := valueObject.NewMappingIdPanic(mappingIdUint)

			mappingQueryRepo := infra.NewMappingQueryRepo(controller.persistentDbSvc)
			mappingCmdRepo := infra.NewMappingCmdRepo(controller.persistentDbSvc)

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

func (controller *MappingController) CreateTarget() *cobra.Command {
	var mappingIdUint uint64
	var targetStr string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "CreateMappingTarget",
		Run: func(cmd *cobra.Command, args []string) {
			mappingId := valueObject.NewMappingIdPanic(mappingIdUint)
			target, err := valueObject.NewContainerId(targetStr)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			createTargetDto := dto.NewCreateMappingTarget(
				mappingId,
				target,
			)

			mappingQueryRepo := infra.NewMappingQueryRepo(controller.persistentDbSvc)
			mappingCmdRepo := infra.NewMappingCmdRepo(controller.persistentDbSvc)
			containerQueryRepo := infra.NewContainerQueryRepo(controller.persistentDbSvc)

			err = useCase.CreateMappingTarget(
				mappingQueryRepo,
				mappingCmdRepo,
				containerQueryRepo,
				createTargetDto,
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
		"ContainerId",
	)
	cmd.MarkFlagRequired("target")
	return cmd
}

func (controller *MappingController) DeleteTarget() *cobra.Command {
	var targetIdUint uint64

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "DeleteMappingTarget",
		Run: func(cmd *cobra.Command, args []string) {
			targetId := valueObject.NewMappingTargetIdPanic(targetIdUint)

			mappingQueryRepo := infra.NewMappingQueryRepo(controller.persistentDbSvc)
			mappingCmdRepo := infra.NewMappingCmdRepo(controller.persistentDbSvc)

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
