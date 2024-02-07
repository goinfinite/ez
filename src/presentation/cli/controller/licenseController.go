package cliController

import (
	"github.com/speedianet/control/src/domain/useCase"
	"github.com/speedianet/control/src/infra"
	"github.com/speedianet/control/src/infra/db"
	cliHelper "github.com/speedianet/control/src/presentation/cli/helper"
	cliMiddleware "github.com/speedianet/control/src/presentation/cli/middleware"
	"github.com/spf13/cobra"
)

func GetLicenseStatusController() *cobra.Command {
	var dbSvc *db.DatabaseService

	cmd := &cobra.Command{
		Use:   "get",
		Short: "GetLicenseStatus",
		PreRun: func(cmd *cobra.Command, args []string) {
			dbSvc = cliMiddleware.DatabaseInit()
		},
		Run: func(cmd *cobra.Command, args []string) {
			licenseQueryRepo := infra.NewLicenseQueryRepo(dbSvc)
			licenseStatus, err := useCase.GetLicenseStatus(licenseQueryRepo)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, licenseStatus)
		},
	}

	return cmd
}
