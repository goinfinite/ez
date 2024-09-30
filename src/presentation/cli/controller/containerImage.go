package cliController

import (
	"github.com/goinfinite/ez/src/infra/db"
	cliHelper "github.com/goinfinite/ez/src/presentation/cli/helper"
	"github.com/goinfinite/ez/src/presentation/service"
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
	var containerIdStr, shouldCreateArchiveBoolStr, archiveCompressionFormatStr string
	var shouldDiscardImageBoolStr string

	cmd := &cobra.Command{
		Use:   "create-snapshot",
		Short: "CreateContainerSnapshotImage",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{
				"containerId":         containerIdStr,
				"shouldCreateArchive": shouldCreateArchiveBoolStr,
				"shouldDiscardImage":  shouldDiscardImageBoolStr,
			}

			if archiveCompressionFormatStr != "" {
				requestBody["archiveCompressionFormat"] = archiveCompressionFormatStr
			}

			cliHelper.ServiceResponseWrapper(
				controller.containerImageService.CreateSnapshot(requestBody, false),
			)
		},
	}

	cmd.Flags().StringVarP(&containerIdStr, "container-id", "c", "", "ContainerId")
	cmd.MarkFlagRequired("container-id")
	cmd.Flags().StringVarP(
		&shouldCreateArchiveBoolStr, "should-create-archive", "r", "false", "ShouldCreateArchive",
	)
	cmd.Flags().StringVarP(
		&archiveCompressionFormatStr, "archive-compression-format", "f", "",
		"ArchiveCompressionFormat (tar|gzip|zip|xz|br)",
	)
	cmd.Flags().StringVarP(
		&shouldDiscardImageBoolStr, "should-discard-image", "d", "false", "ShouldDiscardImage",
	)
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

func (controller *ContainerImageController) ReadArchiveFiles() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "ReadContainerImageArchiveFiles",
		Run: func(cmd *cobra.Command, args []string) {
			cliHelper.ServiceResponseWrapper(
				controller.containerImageService.ReadArchiveFiles(nil),
			)
		},
	}
	return cmd
}

func (controller *ContainerImageController) CreateArchiveFile() *cobra.Command {
	var accountIdUint uint64
	var imageIdStr, compressionFormatStr string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "CreateContainerImageArchiveFile",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{
				"accountId": accountIdUint,
				"imageId":   imageIdStr,
			}

			if compressionFormatStr != "" {
				requestBody["compressionFormat"] = compressionFormatStr
			}

			cliHelper.ServiceResponseWrapper(
				controller.containerImageService.CreateArchiveFile(requestBody, false),
			)
		},
	}

	cmd.Flags().Uint64VarP(&accountIdUint, "account-id", "a", 0, "AccountId")
	cmd.MarkFlagRequired("account-id")
	cmd.Flags().StringVarP(&imageIdStr, "image-id", "i", "", "ImageId")
	cmd.MarkFlagRequired("image-id")
	cmd.Flags().StringVarP(
		&compressionFormatStr, "compression-format", "f", "",
		"CompressionFormat (tar|gzip|zip|xz|br)",
	)
	return cmd
}

func (controller *ContainerImageController) DeleteArchiveFile() *cobra.Command {
	var accountIdUint uint64
	var imageIdStr string

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "DeleteContainerImageArchiveFile",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{
				"accountId": accountIdUint,
				"imageId":   imageIdStr,
			}

			cliHelper.ServiceResponseWrapper(
				controller.containerImageService.DeleteArchiveFile(requestBody),
			)
		},
	}

	cmd.Flags().Uint64VarP(&accountIdUint, "account-id", "a", 0, "AccountId")
	cmd.MarkFlagRequired("account-id")
	cmd.Flags().StringVarP(&imageIdStr, "image-id", "i", "", "ImageId")
	cmd.MarkFlagRequired("image-id")
	return cmd
}
