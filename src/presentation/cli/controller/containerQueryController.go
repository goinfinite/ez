package cliController

import (
	"github.com/speedianet/sfm/src/domain/useCase"
	"github.com/speedianet/sfm/src/infra"
	"github.com/speedianet/sfm/src/infra/db"
	cliHelper "github.com/speedianet/sfm/src/presentation/cli/helper"
	cliMiddleware "github.com/speedianet/sfm/src/presentation/cli/middleware"
	"github.com/spf13/cobra"
)

func GetContainersController() *cobra.Command {
	var dbSvc *db.DatabaseService

	cmd := &cobra.Command{
		Use:   "get",
		Short: "GetContainers",
		PreRun: func(cmd *cobra.Command, args []string) {
			dbSvc = cliMiddleware.DatabaseInit()
		},
		Run: func(cmd *cobra.Command, args []string) {
			containerQueryRepo := infra.NewContainerQueryRepo(dbSvc)
			containersList, err := useCase.GetContainers(containerQueryRepo)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, containersList)
		},
	}

	return cmd
}
