package presenter

import (
	"github.com/labstack/echo/v4"

	"net/http"

	"github.com/goinfinite/ez/src/domain/entity"
	"github.com/goinfinite/ez/src/infra/db"
	"github.com/goinfinite/ez/src/presentation/service"
	"github.com/goinfinite/ez/src/presentation/ui/page"
	presenterHelper "github.com/goinfinite/ez/src/presentation/ui/presenter/helper"
)

type ContainerProfilePresenter struct {
	containerProfileService *service.ContainerProfileService
}

func NewContainerProfilePresenter(
	persistentDbSvc *db.PersistentDatabaseService,
	trailDbSvc *db.TrailDatabaseService,
) *ContainerProfilePresenter {
	return &ContainerProfilePresenter{
		containerProfileService: service.NewContainerProfileService(
			persistentDbSvc, trailDbSvc,
		),
	}
}

func (presenter *ContainerProfilePresenter) Handler(c echo.Context) error {
	responseOutput := presenter.containerProfileService.Read()
	if responseOutput.Status != service.Success {
		return nil
	}

	profileEntities, assertOk := responseOutput.Body.([]entity.ContainerProfile)
	if !assertOk {
		return nil
	}

	pageContent := page.ContainerProfileIndex(profileEntities)
	return presenterHelper.Render(c, pageContent, http.StatusOK)
}
