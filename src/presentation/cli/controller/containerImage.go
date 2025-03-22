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

func (controller *ContainerImageController) ReadArchives() *cobra.Command {
	var imageIdStr string
	var accountIdUint uint64
	var paginationPageNumberUint32 uint32
	var paginationItemsPerPageUint16 uint16
	var paginationSortByStr, paginationSortDirectionStr, paginationLastSeenIdStr string

	cmd := &cobra.Command{
		Use:   "get",
		Short: "ReadContainerImageArchives",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{}
			if accountIdUint != 0 {
				requestBody["accountId"] = accountIdUint
			}
			if imageIdStr != "" {
				requestBody["imageId"] = imageIdStr
			}
			requestBody = cliHelper.PaginationParser(
				requestBody, paginationPageNumberUint32, paginationItemsPerPageUint16,
				paginationSortByStr, paginationSortDirectionStr, paginationLastSeenIdStr,
			)

			cliHelper.ServiceResponseWrapper(
				controller.containerImageService.ReadArchives(requestBody, nil),
			)
		},
	}

	cmd.Flags().Uint64VarP(&accountIdUint, "account-id", "a", 0, "AccountId")
	cmd.Flags().StringVarP(&imageIdStr, "image-id", "i", "", "ImageId")
	cmd.Flags().Uint32VarP(
		&paginationPageNumberUint32, "page-number", "o", 0, "PageNumber (Pagination)",
	)
	cmd.Flags().Uint16VarP(
		&paginationItemsPerPageUint16, "items-per-page", "I", 0, "ItemsPerPage (Pagination)",
	)
	cmd.Flags().StringVarP(
		&paginationSortByStr, "sort-by", "y", "", "SortBy (Pagination)",
	)
	cmd.Flags().StringVarP(
		&paginationSortDirectionStr, "sort-direction", "x", "", "SortDirection (Pagination)",
	)
	cmd.Flags().StringVarP(
		&paginationLastSeenIdStr, "last-seen-id", "L", "", "LastSeenId (Pagination)",
	)
	return cmd
}

func (controller *ContainerImageController) CreateArchive() *cobra.Command {
	var accountIdUint uint64
	var imageIdStr, compressionFormatStr string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "CreateContainerImageArchive",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{
				"accountId": accountIdUint,
				"imageId":   imageIdStr,
			}

			if compressionFormatStr != "" {
				requestBody["compressionFormat"] = compressionFormatStr
			}

			cliHelper.ServiceResponseWrapper(
				controller.containerImageService.CreateArchive(requestBody, false),
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

func (controller *ContainerImageController) DeleteArchive() *cobra.Command {
	var accountIdUint uint64
	var imageIdStr string

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "DeleteContainerImageArchive",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{
				"accountId": accountIdUint,
				"imageId":   imageIdStr,
			}

			cliHelper.ServiceResponseWrapper(
				controller.containerImageService.DeleteArchive(requestBody),
			)
		},
	}

	cmd.Flags().Uint64VarP(&accountIdUint, "account-id", "a", 0, "AccountId")
	cmd.MarkFlagRequired("account-id")
	cmd.Flags().StringVarP(&imageIdStr, "image-id", "i", "", "ImageId")
	cmd.MarkFlagRequired("image-id")
	return cmd
}
