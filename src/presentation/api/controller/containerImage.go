package apiController

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/useCase"
	"github.com/speedianet/control/src/domain/valueObject"
	"github.com/speedianet/control/src/infra"
	"github.com/speedianet/control/src/infra/db"
	apiHelper "github.com/speedianet/control/src/presentation/api/helper"
	"github.com/speedianet/control/src/presentation/service"
)

type ContainerImageController struct {
	persistentDbSvc       *db.PersistentDatabaseService
	containerImageService *service.ContainerImageService
}

func NewContainerImageController(
	persistentDbSvc *db.PersistentDatabaseService,
	trailDbSvc *db.TrailDatabaseService,
) *ContainerImageController {
	return &ContainerImageController{
		persistentDbSvc:       persistentDbSvc,
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
// @Description  Create a new container snapshot image. This is an asynchronous operation.
// @Tags         containerImage
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        createContainerSnapshotImageDto 	  body    dto.CreateContainerSnapshotImage  true  "CreateContainerSnapshotImageDto"
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
// @Description  Export a container image to a file. This is an asynchronous operation.
// @Tags         containerImage
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        exportContainerImageDto 	  body    dto.ExportContainerImage  true "ExportContainerImageDto"
// @Success      201 {object} object{} "ContainerImageExportScheduled"
// @Router       /v1/container/image/export/ [post]
func (controller *ContainerImageController) Export(c echo.Context) error {
	requestBody, err := apiHelper.ReadRequestBody(c)
	if err != nil {
		return err
	}

	return apiHelper.ServiceResponseWrapper(
		c, controller.containerImageService.Export(requestBody, true),
	)
}

// ReadContainerImageArchiveFiles	 godoc
// @Summary      ReadContainerImageArchiveFiles
// @Description  List container image archive files.
// @Tags         containerImage
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200 {array} entity.ContainerImageArchiveFile
// @Router       /v1/container/image/archive/ [get]
func (controller *ContainerImageController) ReadArchiveFiles(c echo.Context) error {
	return apiHelper.ServiceResponseWrapper(
		c, controller.containerImageService.ReadArchiveFiles(),
	)
}

// DownloadContainerImageArchiveFile	 godoc
// @Summary      DownloadContainerImageArchiveFile
// @Description  Download a container image archive file.
// @Tags         containerImage
// @Accept       json
// @Produce      octet-stream
// @Security     Bearer
// @Param        accountId 	  path   string  true  "AccountId"
// @Param        imageId 	  path   string  true  "ImageId"
// @Success      200 file file "ContainerImageArchiveFile"
// @Router       /v1/container/image/archive/{accountId}/{imageId}/ [get]
func (controller *ContainerImageController) ReadArchiveFile(c echo.Context) error {
	if c.Param("accountId") == "" {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, "AccountIdRequired")
	}
	accountId, err := valueObject.NewAccountId(c.Param("accountId"))
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
	}

	if c.Param("imageId") == "" {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, "ImageIdRequired")
	}
	imageId, err := valueObject.NewContainerImageId(c.Param("imageId"))
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
	}

	containerImageQueryRepo := infra.NewContainerImageQueryRepo(controller.persistentDbSvc)

	readDto := dto.NewReadContainerImageArchiveFile(accountId, imageId)

	archiveFile, err := useCase.ReadContainerImageArchiveFile(containerImageQueryRepo, readDto)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return c.Attachment(
		archiveFile.UnixFilePath.String(),
		archiveFile.UnixFilePath.ReadFileName().String(),
	)
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
