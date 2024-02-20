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
	persistDbSvc := c.Get("persistDbSvc").(*db.PersistentDatabaseService)
	licenseQueryRepo := infra.NewLicenseQueryRepo(persistDbSvc)
	licenseStatus, err := useCase.GetLicenseInfo(licenseQueryRepo)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, licenseStatus)
}

func AutoLicenseValidationController(persistDbSvc *db.PersistentDatabaseService) {
	validationIntervalHours := 24 / useCase.LicenseValidationsPerDay

	taskInterval := time.Duration(validationIntervalHours) * time.Hour
	timer := time.NewTicker(taskInterval)
	defer timer.Stop()

	licenseQueryRepo := infra.NewLicenseQueryRepo(persistDbSvc)
	licenseCmdRepo := infra.NewLicenseCmdRepo(persistDbSvc)

	for range timer.C {
		err := useCase.LicenseValidation(licenseQueryRepo, licenseCmdRepo)
		if err != nil {
			log.Fatal(err)
		}
	}
}
