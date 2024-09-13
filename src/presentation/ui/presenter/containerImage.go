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

type ContainerImagePresenter struct {
	containerImageService *service.ContainerImageService
}

func NewContainerImagePresenter(
	persistentDbSvc *db.PersistentDatabaseService,
	trailDbSvc *db.TrailDatabaseService,
) *ContainerImagePresenter {
	return &ContainerImagePresenter{
		containerImageService: service.NewContainerImageService(
			persistentDbSvc, trailDbSvc,
		),
	}
}

func (presenter *ContainerImagePresenter) Handler(c echo.Context) error {
	responseOutput := presenter.containerImageService.Read()
	if responseOutput.Status != service.Success {
		return nil
	}

	imageEntities, assertOk := responseOutput.Body.([]entity.ContainerImage)
	if !assertOk {
		return nil
	}

	pageContent := page.ContainerImageIndex(imageEntities)
	return uiHelper.Render(c, pageContent, http.StatusOK)
}
