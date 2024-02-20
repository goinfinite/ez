package sharedMiddleware

import (
	"log"

	"github.com/speedianet/control/src/domain/useCase"
	"github.com/speedianet/control/src/infra"
	"github.com/speedianet/control/src/infra/db"
)

func InvalidLicenseBlocker(persistentDbSvc *db.PersistentDatabaseService) {
	licenseQueryRepo := infra.NewLicenseQueryRepo(persistentDbSvc)
	containerQueryRepo := infra.NewContainerQueryRepo(persistentDbSvc)
	containerCmdRepo := infra.NewContainerCmdRepo(persistentDbSvc)

	err := useCase.InvalidLicenseBlocker(
		licenseQueryRepo,
		containerQueryRepo,
		containerCmdRepo,
	)
	if err != nil {
		log.Fatal(err)
	}
}
