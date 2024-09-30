package cliMiddleware

import (
	"log"
	"math/rand"

	"github.com/goinfinite/ez/src/domain/useCase"
	"github.com/goinfinite/ez/src/infra"
	"github.com/goinfinite/ez/src/infra/db"
)

func SporadicLicenseValidation(
	persistentDbSvc *db.PersistentDatabaseService,
	transientDbSvc *db.TransientDatabaseService,
) {
	shouldRun := rand.Intn(30) == 0
	if !shouldRun {
		return
	}

	licenseQueryRepo := infra.NewLicenseQueryRepo(persistentDbSvc, transientDbSvc)
	licenseCmdRepo := infra.NewLicenseCmdRepo(persistentDbSvc, transientDbSvc)

	err := useCase.LicenseValidation(licenseQueryRepo, licenseCmdRepo)
	if err != nil {
		log.Fatal(err)
	}
}
