package apiController

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/speedianet/control/src/domain/useCase"
	"github.com/speedianet/control/src/infra/db"
	o11yInfra "github.com/speedianet/control/src/infra/o11y"
	apiHelper "github.com/speedianet/control/src/presentation/api/helper"
)

// O11yOverview  godoc
// @Summary      O11yOverview
// @Description  Show system information and resource usage.
// @Tags         o11y
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200 {object} entity.O11yOverview
// @Router       /v1/o11y/overview/ [get]
func O11yOverviewController(c echo.Context) error {
	transientDbSvc := c.Get("transientDbSvc").(*db.TransientDatabaseService)
	o11yQueryRepo := o11yInfra.NewO11yQueryRepo(transientDbSvc)
	o11yOverview, err := useCase.GetO11yOverview(o11yQueryRepo)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, o11yOverview)
}
