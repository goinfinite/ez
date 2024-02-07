package apiController

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/speedianet/control/src/domain/useCase"
	"github.com/speedianet/control/src/infra"
	"github.com/speedianet/control/src/infra/db"
	apiHelper "github.com/speedianet/control/src/presentation/api/helper"
)

// GetLicenseStatus	 godoc
// @Summary      GetLicenseStatus
// @Description  Get license status.
// @Tags         license
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200 {object} entity.LicenseStatus
// @Router       /license/ [get]
func GetLicenseStatusController(c echo.Context) error {
	dbSvc := c.Get("dbSvc").(*db.DatabaseService)
	licenseQueryRepo := infra.NewLicenseQueryRepo(dbSvc)
	licenseStatus, err := useCase.GetLicenseStatus(licenseQueryRepo)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, licenseStatus)
}
