package sharedMiddleware

import (
	"log"

	"github.com/speedianet/control/src/domain/useCase"
	"github.com/speedianet/control/src/infra"
	"github.com/speedianet/control/src/infra/db"
)

func InvalidLicenseBlocker(dbSvc *db.DatabaseService) {
	licenseQueryRepo := infra.NewLicenseQueryRepo(dbSvc)
	containerQueryRepo := infra.NewContainerQueryRepo(dbSvc)
	containerCmdRepo := infra.NewContainerCmdRepo(dbSvc)

	err := useCase.InvalidLicenseBlocker(
		licenseQueryRepo,
		containerQueryRepo,
		containerCmdRepo,
	)
	if err != nil {
		log.Fatal(err)
	}
}
