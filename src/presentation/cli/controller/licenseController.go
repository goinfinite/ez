package cliController

import (
	"github.com/speedianet/control/src/domain/useCase"
	"github.com/speedianet/control/src/infra"
	"github.com/speedianet/control/src/infra/db"
	cliHelper "github.com/speedianet/control/src/presentation/cli/helper"
	"github.com/spf13/cobra"
)

type LicenseController struct {
	dbSvc *db.DatabaseService
}

func NewLicenseController(dbSvc *db.DatabaseService) LicenseController {
	return LicenseController{dbSvc: dbSvc}
}

func (controller LicenseController) GetLicenseInfo() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "GetLicenseInfo",
		Run: func(cmd *cobra.Command, args []string) {
			licenseQueryRepo := infra.NewLicenseQueryRepo(controller.dbSvc)
			licenseStatus, err := useCase.GetLicenseInfo(licenseQueryRepo)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, licenseStatus)
		},
	}

	return cmd
}
