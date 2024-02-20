package sharedMiddleware

import (
	"log"

	"github.com/speedianet/control/src/domain/useCase"
	"github.com/speedianet/control/src/infra"
	"github.com/speedianet/control/src/infra/db"
)

func InvalidLicenseBlocker(persistDbSvc *db.PersistentDatabaseService) {
	licenseQueryRepo := infra.NewLicenseQueryRepo(persistDbSvc)
	containerQueryRepo := infra.NewContainerQueryRepo(persistDbSvc)
	containerCmdRepo := infra.NewContainerCmdRepo(persistDbSvc)

	err := useCase.InvalidLicenseBlocker(
		licenseQueryRepo,
		containerQueryRepo,
		containerCmdRepo,
	)
	if err != nil {
		log.Fatal(err)
	}
}
