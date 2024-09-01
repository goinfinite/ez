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

func (controller *ContainerImageController) CreateSnapshot() *cobra.Command {
	var accountIdUint uint64
	var containerIdStr string

	cmd := &cobra.Command{
		Use:   "create-snapshot",
		Short: "CreateContainerSnapshotImage",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{
				"accountId":   accountIdUint,
				"containerId": containerIdStr,
			}
			cliHelper.ServiceResponseWrapper(
				controller.containerImageService.CreateSnapshot(requestBody, false),
			)
		},
	}

	cmd.Flags().Uint64VarP(&accountIdUint, "account-id", "a", 0, "AccountId")
	cmd.MarkFlagRequired("account-id")
	cmd.Flags().StringVarP(&containerIdStr, "container-id", "c", "", "ContainerId")
	cmd.MarkFlagRequired("container-id")
	return cmd
}

func (controller *ContainerImageController) Export() *cobra.Command {
	var accountIdUint uint64
	var imageIdStr string

	cmd := &cobra.Command{
		Use:   "export",
		Short: "ExportContainerImage",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{
				"accountId": accountIdUint,
				"imageId":   imageIdStr,
			}
			cliHelper.ServiceResponseWrapper(
				controller.containerImageService.Export(requestBody, false),
			)
		},
	}

	cmd.Flags().Uint64VarP(&accountIdUint, "account-id", "a", 0, "AccountId")
	cmd.MarkFlagRequired("account-id")
	cmd.Flags().StringVarP(&imageIdStr, "image-id", "i", "", "ImageId")
	cmd.MarkFlagRequired("image-id")
	return cmd
}

func (controller *ContainerImageController) ReadArchiveFiles() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-archive-files",
		Short: "ReadContainerImageArchiveFiles",
		Run: func(cmd *cobra.Command, args []string) {
			cliHelper.ServiceResponseWrapper(
				controller.containerImageService.ReadArchiveFiles(),
			)
		},
	}
	return cmd
}

func (controller *ContainerImageController) Delete() *cobra.Command {
	var accountIdUint uint64
	var imageIdStr string

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "DeleteContainerImage",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{
				"accountId": accountIdUint,
				"imageId":   imageIdStr,
			}

			cliHelper.ServiceResponseWrapper(
				controller.containerImageService.Delete(requestBody),
			)
		},
	}

	cmd.Flags().Uint64VarP(&accountIdUint, "account-id", "a", 0, "AccountId")
	cmd.MarkFlagRequired("account-id")
	cmd.Flags().StringVarP(&imageIdStr, "image-id", "i", "", "ImageId")
	cmd.MarkFlagRequired("image-id")
	return cmd
}
