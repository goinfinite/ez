package cliController

import (
	"github.com/goinfinite/ez/src/domain/useCase"
	"github.com/goinfinite/ez/src/domain/valueObject"
	"github.com/goinfinite/ez/src/infra"
	"github.com/goinfinite/ez/src/infra/db"
	cliHelper "github.com/goinfinite/ez/src/presentation/cli/helper"
	"github.com/spf13/cobra"
)

type ContainerRegistryController struct {
	persistentDbSvc *db.PersistentDatabaseService
}

func NewContainerRegistryController(
	persistentDbSvc *db.PersistentDatabaseService,
) *ContainerRegistryController {
	return &ContainerRegistryController{persistentDbSvc: persistentDbSvc}
}

func (controller *ContainerRegistryController) ReadRegistryImages() *cobra.Command {
	var imageNameStr string

	cmd := &cobra.Command{
		Use:   "get",
		Short: "ReadRegistryImages",
		Run: func(cmd *cobra.Command, args []string) {
			containerRegistryQueryRepo := infra.NewContainerRegistryQueryRepo(controller.persistentDbSvc)

			var imageNamePtr *valueObject.RegistryImageName
			if imageNameStr != "" {
				imageName := valueObject.NewRegistryImageNamePanic(imageNameStr)
				imageNamePtr = &imageName
			}

			imagesList, err := useCase.ReadRegistryImages(
				containerRegistryQueryRepo,
				imageNamePtr,
			)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, imagesList)
		},
	}

	cmd.Flags().StringVarP(&imageNameStr, "image-name", "n", "", "ImageName")
	return cmd
}

func (controller *ContainerRegistryController) ReadRegistryTaggedImage() *cobra.Command {
	var imageAddressStr string

	cmd := &cobra.Command{
		Use:   "get-tagged",
		Short: "ReadRegistryTaggedImage",
		Run: func(cmd *cobra.Command, args []string) {
			containerRegistryQueryRepo := infra.NewContainerRegistryQueryRepo(controller.persistentDbSvc)

			imageAddress := valueObject.NewContainerImageAddressPanic(imageAddressStr)

			image, err := useCase.ReadRegistryTaggedImage(
				containerRegistryQueryRepo,
				imageAddress,
			)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, image)
		},
	}

	cmd.Flags().StringVarP(&imageAddressStr, "image-address", "a", "", "ImageAddress")
	cmd.MarkFlagRequired("image-address")
	return cmd
}
