package presenter

import (
	"log/slog"
	"net/http"

	"github.com/goinfinite/ez/src/domain/entity"
	"github.com/goinfinite/ez/src/infra/db"
	"github.com/goinfinite/ez/src/presentation/service"
	uiHelper "github.com/goinfinite/ez/src/presentation/ui/helper"
	"github.com/goinfinite/ez/src/presentation/ui/page"
	"github.com/labstack/echo/v4"
)

type OverviewPresenter struct {
	persistentDbSvc *db.PersistentDatabaseService
	transientDbSvc  *db.TransientDatabaseService
	trailDbSvc      *db.TrailDatabaseService
}

func NewOverviewPresenter(
	persistentDbSvc *db.PersistentDatabaseService,
	transientDbSvc *db.TransientDatabaseService,
	trailDbSvc *db.TrailDatabaseService,
) *OverviewPresenter {
	return &OverviewPresenter{
		persistentDbSvc: persistentDbSvc,
		transientDbSvc:  transientDbSvc,
		trailDbSvc:      trailDbSvc,
	}
}

func (presenter *OverviewPresenter) Handler(c echo.Context) error {
	containerService := service.NewContainerService(
		presenter.persistentDbSvc, presenter.trailDbSvc,
	)

	readContainersServiceOutput := containerService.Read()
	if readContainersServiceOutput.Status != service.Success {
		slog.Debug("OverviewPresenterReadContainersFailure")
		return nil
	}

	containerEntities, assertOk := readContainersServiceOutput.Body.([]entity.Container)
	if !assertOk {
		slog.Debug("OverviewPresenterAssertContainersFailure")
		return nil
	}

	pageContent := page.OverviewIndex(containerEntities)
	return uiHelper.Render(c, pageContent, http.StatusOK)
}
