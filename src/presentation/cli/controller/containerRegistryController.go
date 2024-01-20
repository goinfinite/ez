package cliController

import (
	"github.com/speedianet/control/src/domain/useCase"
	"github.com/speedianet/control/src/domain/valueObject"
	"github.com/speedianet/control/src/infra"
	"github.com/speedianet/control/src/infra/db"
	cliHelper "github.com/speedianet/control/src/presentation/cli/helper"
	cliMiddleware "github.com/speedianet/control/src/presentation/cli/middleware"
	"github.com/spf13/cobra"
)

func GetRegistryImagesController() *cobra.Command {
	var dbSvc *db.DatabaseService
	var imageNameStr string

	cmd := &cobra.Command{
		Use:   "get",
		Short: "GetRegistryImages",
		PreRun: func(cmd *cobra.Command, args []string) {
			dbSvc = cliMiddleware.DatabaseInit()
		},
		Run: func(cmd *cobra.Command, args []string) {
			containerRegistryQueryRepo := infra.NewContainerRegistryQueryRepo(dbSvc)

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
