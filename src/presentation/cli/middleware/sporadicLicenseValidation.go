package cliMiddleware

import (
	"log"
	"math/rand"

	"github.com/speedianet/control/src/domain/useCase"
	"github.com/speedianet/control/src/infra"
	"github.com/speedianet/control/src/infra/db"
)

func SporadicLicenseValidation(persistDbSvc *db.PersistentDatabaseService) {
	shouldRun := rand.Intn(30) == 0
	if !shouldRun {
		return
	}

	licenseQueryRepo := infra.NewLicenseQueryRepo(persistDbSvc)
	licenseCmdRepo := infra.NewLicenseCmdRepo(persistDbSvc)

	err := useCase.LicenseValidation(licenseQueryRepo, licenseCmdRepo)
	if err != nil {
		log.Fatal(err)
	}
}
