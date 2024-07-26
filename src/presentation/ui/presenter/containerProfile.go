package presenter

import (
	"github.com/labstack/echo/v4"

	"net/http"

	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/infra/db"
	"github.com/speedianet/control/src/presentation/service"
	uiHelper "github.com/speedianet/control/src/presentation/ui/helper"
	"github.com/speedianet/control/src/presentation/ui/page"
)

type ContainerProfilePresenter struct {
	containerProfileService *service.ContainerProfileService
}

func NewContainerProfilePresenter(
	persistentDbSvc *db.PersistentDatabaseService,
) *ContainerProfilePresenter {
	return &ContainerProfilePresenter{
		containerProfileService: service.NewContainerProfileService(persistentDbSvc),
	}
}

func (presenter *ContainerProfilePresenter) Handler(c echo.Context) error {
	responseOutput := presenter.containerProfileService.Read()
	if responseOutput.Status != service.Success {
		return nil
	}

	contactEntities, assertOk := responseOutput.Body.([]entity.ContainerProfile)
	if !assertOk {
		return nil
	}

	pageContent := page.ContainerProfileIndex(contactEntities)
	return uiHelper.Render(c, pageContent, http.StatusOK)
}