package presenter

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/infra/db"
	"github.com/speedianet/control/src/presentation/service"
	"github.com/speedianet/control/src/presentation/ui/layout"
)

type FooterPresenter struct {
	persistentDbSvc *db.PersistentDatabaseService
	transientDbSvc  *db.TransientDatabaseService
	trailDbSvc      *db.TrailDatabaseService
}

func NewFooterPresenter(
	persistentDbSvc *db.PersistentDatabaseService,
	transientDbSvc *db.TransientDatabaseService,
	trailDbSvc *db.TrailDatabaseService,
) *FooterPresenter {
	return &FooterPresenter{
		persistentDbSvc: persistentDbSvc,
		transientDbSvc:  transientDbSvc,
		trailDbSvc:      trailDbSvc,
	}
}

func (presenter *FooterPresenter) Handler(c echo.Context) error {
	o11yService := service.NewO11yService(presenter.transientDbSvc)

	o11yServiceOutput := o11yService.ReadOverview()
	if o11yServiceOutput.Status != service.Success {
		return nil
	}

	o11yOverviewEntity, assertOk := o11yServiceOutput.Body.(entity.O11yOverview)
	if !assertOk {
		return nil
	}

	scheduledTaskService := service.NewScheduledTaskService(presenter.persistentDbSvc)

	scheduledTaskServiceOutput := scheduledTaskService.Read()
	if scheduledTaskServiceOutput.Status != service.Success {
		return nil
	}

	scheduledTaskEntities, assertOk := scheduledTaskServiceOutput.Body.([]entity.ScheduledTask)
	if !assertOk {
		return nil
	}

	c.Response().Writer.WriteHeader(http.StatusOK)
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)

	return layout.Footer(o11yOverviewEntity, scheduledTaskEntities).
		Render(c.Request().Context(), c.Response().Writer)
}
