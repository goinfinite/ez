package apiController

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/infra/db"
	apiHelper "github.com/speedianet/control/src/presentation/api/helper"
	"github.com/speedianet/control/src/presentation/service"
)

type ContainerImageController struct {
	containerImageService *service.ContainerImageService
}

func NewContainerImageController(
	persistentDbSvc *db.PersistentDatabaseService,
	trailDbSvc *db.TrailDatabaseService,
) *ContainerImageController {
	return &ContainerImageController{
		containerImageService: service.NewContainerImageService(persistentDbSvc, trailDbSvc),
	}
}

// ReadContainerImages	 godoc
// @Summary      ReadContainerImages
// @Description  List container images.
// @Tags         containerImage
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200 {array} entity.ContainerImage
// @Router       /v1/container/image/ [get]
func (controller *ContainerImageController) Read(c echo.Context) error {
	return apiHelper.ServiceResponseWrapper(c, controller.containerImageService.Read())
}

// CreateContainerSnapshotImage	 godoc
// @Summary      CreateContainerSnapshotImage
// @Description  Create a new container snapshot image.
// @Tags         containerImage
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        createContainerSnapshotImageDto 	  body    dto.CreateContainerSnapshotImage  true  "Asynchronous Snapshot Image Creation"
// @Success      201 {object} object{} "ContainerSnapshotImageCreationScheduled"
// @Router       /v1/container/image/snapshot/ [post]
func (controller *ContainerImageController) CreateSnapshot(c echo.Context) error {
	requestBody, err := apiHelper.ReadRequestBody(c)
	if err != nil {
		return err
	}

	return apiHelper.ServiceResponseWrapper(
		c, controller.containerImageService.CreateSnapshot(requestBody, true),
	)
}

// ExportContainerImage	 godoc
// @Summary      ExportContainerImage
// @Description  Export a container image.
// @Tags         containerImage
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        exportContainerImageDto 	  body    dto.ExportContainerImage  true "ExportContainerImageDto"
// @Success      302 {string} string "Redirect to the archive image file download URL."
// @Router       /v1/container/image/export/ [post]
func (controller *ContainerImageController) Export(c echo.Context) error {
	requestBody, err := apiHelper.ReadRequestBody(c)
	if err != nil {
		return err
	}

	serviceResponse := controller.containerImageService.Export(requestBody)
	if serviceResponse.Status != service.Success {
		return apiHelper.ServiceResponseWrapper(c, serviceResponse)
	}

	archiveFile, assertOk := serviceResponse.Body.(entity.ContainerImageArchiveFile)
	if !assertOk {
		return apiHelper.ServiceResponseWrapper(c, serviceResponse)
	}

	return c.Redirect(http.StatusFound, archiveFile.DownloadUrl.String())
}

// DeleteContainerImage godoc
// @Summary      DeleteContainerImage
// @Description  Delete a container image.
// @Tags         containerImage
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        accountId 	  path   string  true  "AccountId"
// @Param        imageId 	  path   string  true  "ImageId"
// @Success      200 {object} object{} "ContainerImageDeleted"
// @Router       /v1/container/image/{accountId}/{imageId}/ [delete]
func (controller *ContainerImageController) Delete(c echo.Context) error {
	requestBody := map[string]interface{}{
		"accountId": c.Param("accountId"),
		"imageId":   c.Param("imageId"),
	}

	return apiHelper.ServiceResponseWrapper(
		c, controller.containerImageService.Delete(requestBody),
	)
}
