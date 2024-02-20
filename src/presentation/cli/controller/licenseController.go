package cliController

import (
	"github.com/speedianet/control/src/domain/useCase"
	"github.com/speedianet/control/src/infra"
	"github.com/speedianet/control/src/infra/db"
	cliHelper "github.com/speedianet/control/src/presentation/cli/helper"
	"github.com/spf13/cobra"
)

type LicenseController struct {
	persistentDbSvc *db.PersistentDatabaseService
}

func NewLicenseController(persistentDbSvc *db.PersistentDatabaseService) LicenseController {
	return LicenseController{persistentDbSvc: persistentDbSvc}
}

func (controller LicenseController) GetLicenseInfo() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "GetLicenseInfo",
		Run: func(cmd *cobra.Command, args []string) {
			licenseQueryRepo := infra.NewLicenseQueryRepo(controller.persistentDbSvc)
			licenseStatus, err := useCase.GetLicenseInfo(licenseQueryRepo)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, licenseStatus)
		},
	}

	return cmd
}
