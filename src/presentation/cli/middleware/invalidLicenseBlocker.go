package cliMiddleware

import (
	"log"

	"github.com/goinfinite/ez/src/domain/useCase"
	"github.com/goinfinite/ez/src/infra"
	"github.com/goinfinite/ez/src/infra/db"
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
