package cliMiddleware

import (
	"log"

	"github.com/speedianet/control/src/domain/useCase"
	"github.com/speedianet/control/src/infra"
	"github.com/speedianet/control/src/infra/db"
)

func InvalidLicenseBlocker(
	persistentDbSvc *db.PersistentDatabaseService,
	transientDbSvc *db.TransientDatabaseService,
) {
	licenseQueryRepo := infra.NewLicenseQueryRepo(persistentDbSvc, transientDbSvc)
	containerQueryRepo := infra.NewContainerQueryRepo(persistentDbSvc)
	containerCmdRepo := infra.NewContainerCmdRepo(persistentDbSvc)

	err := useCase.InvalidLicenseBlocker(
		licenseQueryRepo, containerQueryRepo, containerCmdRepo,
	)
	if err != nil {
		log.Fatal(err)
	}
}
