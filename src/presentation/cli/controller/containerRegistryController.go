package cliController

import (
	"github.com/speedianet/control/src/domain/useCase"
	"github.com/speedianet/control/src/domain/valueObject"
	"github.com/speedianet/control/src/infra"
	"github.com/speedianet/control/src/infra/db"
	cliHelper "github.com/speedianet/control/src/presentation/cli/helper"
	"github.com/spf13/cobra"
)

type ContainerRegistryController struct {
	persistDbSvc *db.PersistentDatabaseService
}

func NewContainerRegistryController(persistDbSvc *db.PersistentDatabaseService) ContainerRegistryController {
	return ContainerRegistryController{persistDbSvc: persistDbSvc}
}

func (controller ContainerRegistryController) GetRegistryImages() *cobra.Command {
	var imageNameStr string

	cmd := &cobra.Command{
		Use:   "get",
		Short: "GetRegistryImages",
		Run: func(cmd *cobra.Command, args []string) {
			containerRegistryQueryRepo := infra.NewContainerRegistryQueryRepo(controller.persistDbSvc)

			var imageNamePtr *valueObject.RegistryImageName
			if imageNameStr != "" {
				imageName := valueObject.NewRegistryImageNamePanic(imageNameStr)
				imageNamePtr = &imageName
			}

			imagesList, err := useCase.GetRegistryImages(
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

func (controller ContainerRegistryController) GetRegistryTaggedImage() *cobra.Command {
	var imageAddressStr string

	cmd := &cobra.Command{
		Use:   "get-tagged",
		Short: "GetRegistryTaggedImage",
		Run: func(cmd *cobra.Command, args []string) {
			containerRegistryQueryRepo := infra.NewContainerRegistryQueryRepo(controller.persistDbSvc)

			imageAddress := valueObject.NewContainerImageAddressPanic(imageAddressStr)

			image, err := useCase.GetRegistryTaggedImage(
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
