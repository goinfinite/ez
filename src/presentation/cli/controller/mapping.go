package cliController

import (
	"github.com/speedianet/control/src/domain/valueObject"
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
			cliHelper.ServiceResponseWrapper(controller.mappingService.Read())
		},
	}

	return cmd
}

func (controller *MappingController) Create() *cobra.Command {
	var accountIdUint uint64
	var hostnameStr string
	var publicPortUint uint64
	var networkProtocolStr string
	var containerIdStrSlice []string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "CreateMapping",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{
				"accountId":  accountIdUint,
				"publicPort": publicPortUint,
			}

			if hostnameStr != "" {
				requestBody["hostname"] = hostnameStr
			}

			if networkProtocolStr != "" {
				requestBody["protocol"] = networkProtocolStr
			}

			containerIds := []valueObject.ContainerId{}
			for _, containerIdStr := range containerIdStrSlice {
				containerId, err := valueObject.NewContainerId(containerIdStr)
				if err != nil {
					cliHelper.ResponseWrapper(false, err.Error())
				}
				containerIds = append(containerIds, containerId)
			}
			requestBody["containerIds"] = containerIds

			cliHelper.ServiceResponseWrapper(
				controller.mappingService.Create(requestBody),
			)
		},
	}
	cmd.Flags().Uint64VarP(&accountIdUint, "account-id", "a", 0, "AccountId")
	cmd.MarkFlagRequired("account-id")
	cmd.Flags().StringVarP(&hostnameStr, "hostname", "n", "", "Hostname")
	cmd.Flags().Uint64VarP(&publicPortUint, "port", "p", 0, "Public Port")
	cmd.MarkFlagRequired("port")
	cmd.Flags().StringVarP(&networkProtocolStr, "protocol", "l", "", "Host Protocol")
	cmd.Flags().StringSliceVarP(
		&containerIdStrSlice, "container-ids", "c", []string{}, "ContainerIds",
	)
	cmd.MarkFlagRequired("container-ids")
	return cmd
}

func (controller *MappingController) Delete() *cobra.Command {
	var mappingIdUint uint64

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "DeleteMapping",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{
				"mappingId": mappingIdUint,
			}

			cliHelper.ServiceResponseWrapper(
				controller.mappingService.Delete(requestBody),
			)
		},
	}

	cmd.Flags().Uint64VarP(&mappingIdUint, "id", "i", 0, "MappingId")
	cmd.MarkFlagRequired("id")
	return cmd
}

func (controller *MappingController) CreateTarget() *cobra.Command {
	var mappingIdUint uint64
	var containerIdStr string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "CreateMappingTarget",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{
				"mappingId":   mappingIdUint,
				"containerId": containerIdStr,
			}

			cliHelper.ServiceResponseWrapper(
				controller.mappingService.CreateTarget(requestBody),
			)
		},
	}

	cmd.Flags().Uint64VarP(&mappingIdUint, "mapping-id", "m", 0, "MappingId")
	cmd.MarkFlagRequired("mapping-id")
	cmd.Flags().StringVarP(&containerIdStr, "container-id", "c", "", "ContainerId")
	cmd.MarkFlagRequired("container-id")
	return cmd
}

func (controller *MappingController) DeleteTarget() *cobra.Command {
	var mappingIdUint uint64
	var targetIdUint uint64

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "DeleteMappingTarget",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{
				"mappingId": mappingIdUint,
				"targetId":  targetIdUint,
			}

			cliHelper.ServiceResponseWrapper(
				controller.mappingService.DeleteTarget(requestBody),
			)
		},
	}

	cmd.Flags().Uint64VarP(&mappingIdUint, "mapping-id", "m", 0, "MappingId")
	cmd.MarkFlagRequired("mapping-id")
	cmd.Flags().Uint64VarP(&targetIdUint, "target-id", "t", 0, "MappingTargetId")
	cmd.MarkFlagRequired("target-id")
	return cmd
}
