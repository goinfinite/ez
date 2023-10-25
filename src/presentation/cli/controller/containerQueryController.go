package cliController

import (
	"github.com/goinfinite/fleet/src/domain/useCase"
	"github.com/goinfinite/fleet/src/infra"
	"github.com/goinfinite/fleet/src/infra/db"
	cliHelper "github.com/goinfinite/fleet/src/presentation/cli/helper"
	cliMiddleware "github.com/goinfinite/fleet/src/presentation/cli/middleware"
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
