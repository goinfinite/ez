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
	trailDbSvc *db.TrailDatabaseService,
) *ContainerImageController {
	return &ContainerImageController{
		containerImageService: service.NewContainerImageService(persistentDbSvc, trailDbSvc),
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

func (controller *ContainerImageController) Delete() *cobra.Command {
	var accountIdStr, imageIdStr string

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "DeleteContainerImage",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{
				"accountId": accountIdStr,
				"imageId":   imageIdStr,
			}

			cliHelper.ServiceResponseWrapper(
				controller.containerImageService.Delete(requestBody),
			)
		},
	}

	cmd.Flags().StringVarP(&accountIdStr, "account-id", "a", "", "AccountId")
	cmd.MarkFlagRequired("account-id")
	cmd.Flags().StringVarP(&imageIdStr, "image-id", "i", "", "ImageId")
	cmd.MarkFlagRequired("image-id")
	return cmd
}
