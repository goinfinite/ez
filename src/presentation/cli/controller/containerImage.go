package cliController

import (
	"github.com/speedianet/control/src/infra/db"
	cliHelper "github.com/speedianet/control/src/presentation/cli/helper"
	"github.com/speedianet/control/src/presentation/service"
	"github.com/spf13/cobra"
)

type ContainerImageController struct {
	containerImageService *service.ContainerImageService
}

func NewContainerImageController(
	persistentDbSvc *db.PersistentDatabaseService,
) *ContainerImageController {
	return &ContainerImageController{
		containerImageService: service.NewContainerImageService(persistentDbSvc),
	}
}

func (controller *ContainerImageController) Read() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "ReadContainerImages",
		Run: func(cmd *cobra.Command, args []string) {
			cliHelper.ServiceResponseWrapper(controller.containerImageService.Read())
		},
	}
	return cmd
}
