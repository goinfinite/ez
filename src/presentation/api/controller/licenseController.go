package apiController

import (
	"log"
	"net/http"
	"time"

	"github.com/goinfinite/ez/src/domain/useCase"
	"github.com/goinfinite/ez/src/infra"
	"github.com/goinfinite/ez/src/infra/db"
	apiHelper "github.com/goinfinite/ez/src/presentation/api/helper"
	"github.com/labstack/echo/v4"
)

// ReadLicenseInfo	 godoc
// @Summary      ReadLicenseInfo
// @Description  Get license info.
// @Tags         license
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200 {object} entity.LicenseInfo
// @Router       /v1/license/ [get]
func ReadLicenseInfoController(c echo.Context) error {
	persistentDbSvc := c.Get("persistentDbSvc").(*db.PersistentDatabaseService)
	transientDbSvc := c.Get("transientDbSvc").(*db.TransientDatabaseService)
	licenseQueryRepo := infra.NewLicenseQueryRepo(persistentDbSvc, transientDbSvc)
	licenseStatus, err := useCase.ReadLicenseInfo(licenseQueryRepo)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, licenseStatus)
}

func AutoLicenseValidationController(
	persistentDbSvc *db.PersistentDatabaseService,
	transientDbSvc *db.TransientDatabaseService,
) {
	validationIntervalHours := 24 / useCase.LicenseValidationsPerDay

	taskInterval := time.Duration(validationIntervalHours) * time.Hour
	timer := time.NewTicker(taskInterval)
	defer timer.Stop()

	licenseQueryRepo := infra.NewLicenseQueryRepo(persistentDbSvc, transientDbSvc)
	licenseCmdRepo := infra.NewLicenseCmdRepo(persistentDbSvc, transientDbSvc)

	for range timer.C {
		err := useCase.LicenseValidation(licenseQueryRepo, licenseCmdRepo)
		if err != nil {
			log.Fatal(err)
		}
	}
}
