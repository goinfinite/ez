package apiController

import (
	"github.com/labstack/echo/v4"
	"github.com/speedianet/control/src/infra/db"
	apiHelper "github.com/speedianet/control/src/presentation/api/helper"
	"github.com/speedianet/control/src/presentation/service"
)

type ContainerImageController struct {
	containerImageService *service.ContainerImageService
}

func NewContainerImageController(
	persistentDbSvc *db.PersistentDatabaseService,
) *ContainerImageController {
	return &ContainerImageController{
		containerImageService: service.NewContainerImageService(persistentDbSvc),
	}
}

// ReadContainerImages	 godoc
// @Summary      ReadContainerImages
// @Description  List container images.
// @Tags         container
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200 {array} entity.ContainerImage
// @Router       /v1/container/image/ [get]
func (controller *ContainerImageController) Read(c echo.Context) error {
	return apiHelper.ServiceResponseWrapper(c, controller.containerImageService.Read())
}
