package cliMiddleware

import (
	"log"
	"math/rand"

	"github.com/speedianet/control/src/domain/useCase"
	"github.com/speedianet/control/src/infra"
	"github.com/speedianet/control/src/infra/db"
)

func SporadicLicenseValidation(persistentDbSvc *db.PersistentDatabaseService) {
	shouldRun := rand.Intn(30) == 0
	if !shouldRun {
		return
	}

	licenseQueryRepo := infra.NewLicenseQueryRepo(persistentDbSvc)
	licenseCmdRepo := infra.NewLicenseCmdRepo(persistentDbSvc)

	err := useCase.LicenseValidation(licenseQueryRepo, licenseCmdRepo)
	if err != nil {
		log.Fatal(err)
	}
}
