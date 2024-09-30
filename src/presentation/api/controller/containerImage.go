package apiController

import (
	"mime/multipart"
	"net/http"

	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/useCase"
	"github.com/goinfinite/ez/src/domain/valueObject"
	"github.com/goinfinite/ez/src/infra"
	"github.com/goinfinite/ez/src/infra/db"
	apiHelper "github.com/goinfinite/ez/src/presentation/api/helper"
	"github.com/goinfinite/ez/src/presentation/service"
	serviceHelper "github.com/goinfinite/ez/src/presentation/service/helper"
	"github.com/labstack/echo/v4"
)

type ContainerImageController struct {
	persistentDbSvc       *db.PersistentDatabaseService
	trailDbSvc            *db.TrailDatabaseService
	containerImageService *service.ContainerImageService
}

func NewContainerImageController(
	persistentDbSvc *db.PersistentDatabaseService,
	trailDbSvc *db.TrailDatabaseService,
) *ContainerImageController {
	return &ContainerImageController{
		persistentDbSvc:       persistentDbSvc,
		trailDbSvc:            trailDbSvc,
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
// @Param        createContainerSnapshotImageDto 	  body    dto.CreateContainerSnapshotImage  true  "Container's owner account must have enough quota to store the snapshot image (and/or archive).<br />'shouldCreateArchive' and 'shouldDiscardImage' are optional and default to false if not provided.<br/>'shouldDiscardImage' is only effective when 'shouldCreateArchive' is true and it will delete the snapshot image after creating the archive file.<br /> 'archiveCompressionFormat' is optional and defaults to 'br' if not provided. Although it's possible to provide other values, it's recommended to use 'br' for best speed/compression ratio."
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
		"accountId":         c.Param("accountId"),
		"imageId":           c.Param("imageId"),
		"operatorAccountId": c.Get("accountId"),
		"operatorIpAddress": c.RealIP(),
	}

	return apiHelper.ServiceResponseWrapper(
		c, controller.containerImageService.Delete(requestBody),
	)
}

// ReadContainerImageArchiveFiles	 godoc
// @Summary      ReadContainerImageArchiveFiles
// @Description  List container image archive files.
// @Tags         containerImageArchive
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200 {array} entity.ContainerImageArchiveFile
// @Router       /v1/container/image/archive/ [get]
func (controller *ContainerImageController) ReadArchiveFiles(c echo.Context) error {
	return apiHelper.ServiceResponseWrapper(
		c, controller.containerImageService.ReadArchiveFiles(&c.Request().Host),
	)
}

// DownloadContainerImageArchiveFile	 godoc
// @Summary      DownloadContainerImageArchiveFile
// @Description  Download a container image archive file.
// @Tags         containerImageArchive
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

// CreateContainerImageArchiveFile	 godoc
// @Summary      CreateContainerImageArchiveFile
// @Description  Export a container image to a file. This is an asynchronous operation.
// @Tags         containerImageArchive
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        createContainerImageArchiveFileDto 	  body    dto.CreateContainerImageArchiveFile  true  "CreateContainerImageArchiveFileDto"
// @Success      201 {object} object{} "ContainerImageArchiveFileCreationScheduled"
// @Router       /v1/container/image/archive/ [post]
func (controller *ContainerImageController) CreateArchiveFile(c echo.Context) error {
	requestBody, err := apiHelper.ReadRequestBody(c)
	if err != nil {
		return err
	}

	if requestBody["accountId"] == nil {
		requestBody["accountId"] = requestBody["operatorAccountId"]
	}

	return apiHelper.ServiceResponseWrapper(
		c, controller.containerImageService.CreateArchiveFile(requestBody, true),
	)
}

type FailedArchiveFileUpload struct {
	FileName   string `json:"fileName"`
	FailReason string `json:"failReason"`
}

// ImportContainerImageArchiveFiles	godoc
// @Summary      ImportContainerImageArchiveFiles
// @Description  Import container images from archive files.
// @Tags         containerImageArchive
// @Accept       mpfd
// @Produce      json
// @Security     Bearer
// @Param        accountId		formData	string	false	"AccountId"
// @Param        archiveFiles	formData	file	true	"ArchiveFiles"
// @Success      201 string string "ContainerImageArchiveFilesImported"
// @Success      207 {object} []FailedArchiveFileUpload "ContainerImageArchiveFilesPartiallyImported"
// @Router       /v1/container/image/archive/import/ [post]
func (controller *ContainerImageController) ImportArchiveFile(c echo.Context) error {
	requestBody, err := apiHelper.ReadRequestBody(c)
	if err != nil {
		return err
	}

	requiredParams := []string{"operatorAccountId", "operatorIpAddress"}
	err = serviceHelper.RequiredParamsInspector(requestBody, requiredParams)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
	}

	if requestBody["accountId"] == nil {
		requestBody["accountId"] = requestBody["operatorAccountId"]
	}
	accountId, err := valueObject.NewAccountId(requestBody["accountId"])
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
	}

	if requestBody["archiveFiles"] == nil {
		if requestBody["archiveFile"] == nil {
			return apiHelper.ResponseWrapper(c, http.StatusBadRequest, "ArchiveFilesRequired")
		}
		requestBody["archiveFiles"] = []interface{}{requestBody["archiveFile"]}
	}

	archiveFiles, assertOk := requestBody["archiveFiles"].([]*multipart.FileHeader)
	if !assertOk {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, "ArchiveFilesNotValid")
	}

	operatorAccountId, err := valueObject.NewAccountId(requestBody["operatorAccountId"])
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
	}

	operatorIpAddress, err := valueObject.NewIpAddress(requestBody["operatorIpAddress"])
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
	}

	containerImageCmdRepo := infra.NewContainerImageCmdRepo(controller.persistentDbSvc)
	accountQueryRepo := infra.NewAccountQueryRepo(controller.persistentDbSvc)
	activityRecordCmdRepo := infra.NewActivityRecordCmdRepo(controller.trailDbSvc)

	failedUploads := []FailedArchiveFileUpload{}
	for _, archiveFile := range archiveFiles {
		importDto := dto.NewImportContainerImageArchiveFile(
			accountId, archiveFile, operatorAccountId, operatorIpAddress,
		)

		_, err = useCase.ImportContainerImageArchiveFile(
			containerImageCmdRepo, accountQueryRepo, activityRecordCmdRepo, importDto,
		)
		if err != nil {
			failedUpload := FailedArchiveFileUpload{
				FileName:   archiveFile.Filename,
				FailReason: err.Error(),
			}
			failedUploads = append(failedUploads, failedUpload)
			continue
		}
	}

	if len(failedUploads) > 0 {
		return apiHelper.ResponseWrapper(c, http.StatusMultiStatus, failedUploads)
	}

	return apiHelper.ResponseWrapper(
		c, http.StatusCreated, "ContainerImageArchiveFilesImported",
	)
}

// DeleteContainerImageArchiveFile	 godoc
// @Summary      DeleteContainerImageArchiveFile
// @Description  Delete a container image archive file.
// @Tags         containerImageArchive
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        accountId 	  path   string  true  "AccountId"
// @Param        imageId 	  path   string  true  "ImageId"
// @Success      200 {object} object{} "ContainerImageArchiveFileDeleted"
// @Router       /v1/container/image/archive/{accountId}/{imageId}/ [delete]
func (controller *ContainerImageController) DeleteArchiveFile(c echo.Context) error {
	requestBody := map[string]interface{}{
		"accountId":         c.Param("accountId"),
		"imageId":           c.Param("imageId"),
		"operatorAccountId": c.Get("accountId"),
		"operatorIpAddress": c.RealIP(),
	}

	return apiHelper.ServiceResponseWrapper(
		c, controller.containerImageService.DeleteArchiveFile(requestBody),
	)
}
