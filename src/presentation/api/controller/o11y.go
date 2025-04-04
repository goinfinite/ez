package apiController

import (
	"github.com/goinfinite/ez/src/infra/db"
	apiHelper "github.com/goinfinite/ez/src/presentation/api/helper"
	"github.com/goinfinite/ez/src/presentation/service"
	"github.com/labstack/echo/v4"
)

type O11yController struct {
	o11yService *service.O11yService
}

func NewO11yController(
	transientDbSvc *db.TransientDatabaseService,
) *O11yController {
	return &O11yController{o11yService: service.NewO11yService(transientDbSvc)}
}

// O11yOverview  godoc
// @Summary      O11yOverview
// @Description  Show system information and resource usage.
// @Tags         o11y
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200 {object} entity.O11yOverview
// @Router       /v1/o11y/overview/ [get]
func (controller *O11yController) ReadOverview(c echo.Context) error {
	return apiHelper.ServiceResponseWrapper(c, controller.o11yService.ReadOverview())
}
