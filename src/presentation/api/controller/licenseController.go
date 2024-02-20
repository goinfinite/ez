package apiController

import (
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/speedianet/control/src/domain/useCase"
	"github.com/speedianet/control/src/infra"
	"github.com/speedianet/control/src/infra/db"
	apiHelper "github.com/speedianet/control/src/presentation/api/helper"
)

// GetLicenseInfo	 godoc
// @Summary      GetLicenseInfo
// @Description  Get license info.
// @Tags         license
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200 {object} entity.LicenseInfo
// @Router       /license/ [get]
func GetLicenseInfoController(c echo.Context) error {
	persistentDbSvc := c.Get("persistentDbSvc").(*db.PersistentDatabaseService)
	licenseQueryRepo := infra.NewLicenseQueryRepo(persistentDbSvc)
	licenseStatus, err := useCase.GetLicenseInfo(licenseQueryRepo)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, licenseStatus)
}

func AutoLicenseValidationController(persistentDbSvc *db.PersistentDatabaseService) {
	validationIntervalHours := 24 / useCase.LicenseValidationsPerDay

	taskInterval := time.Duration(validationIntervalHours) * time.Hour
	timer := time.NewTicker(taskInterval)
	defer timer.Stop()

	licenseQueryRepo := infra.NewLicenseQueryRepo(persistentDbSvc)
	licenseCmdRepo := infra.NewLicenseCmdRepo(persistentDbSvc)

	for range timer.C {
		err := useCase.LicenseValidation(licenseQueryRepo, licenseCmdRepo)
		if err != nil {
			log.Fatal(err)
		}
	}
}
