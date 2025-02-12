package apiController

import (
	"mime/multipart"
	"net/http"
	"strings"

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

// ReadContainerImageArchives	 godoc
// @Summary      ReadContainerImageArchives
// @Description  List container image archive files.
// @Tags         containerImageArchive
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        imageId query  string  false  "ImageId"
// @Param        accountId query  uint  false  "AccountId"
// @Param        createdBeforeAt query  string  false  "CreatedBeforeAt"
// @Param        createdAfterAt query  string  false  "CreatedAfterAt"
// @Param        pageNumber query  uint  false  "PageNumber (Pagination)"
// @Param        itemsPerPage query  uint  false  "ItemsPerPage (Pagination)"
// @Param        sortBy query  string  false  "SortBy (Pagination)"
// @Param        sortDirection query  string  false  "SortDirection (Pagination)"
// @Param        lastSeenId query  string  false  "LastSeenId (Pagination)"
// @Success      200 {object} dto.ReadContainerImageArchivesResponse
// @Router       /v1/container/image/archive/ [get]
func (controller *ContainerImageController) ReadArchives(c echo.Context) error {
	requestBody := map[string]interface{}{}
	queryParameters := []string{
		"imageId", "accountId", "createdBeforeAt", "createdAfterAt",
		"pageNumber", "itemsPerPage", "sortBy", "sortDirection", "lastSeenId",
	}
	for _, paramName := range queryParameters {
		paramValue := c.QueryParam(paramName)
		if paramValue == "" {
			continue
		}

		requestBody[paramName] = strings.Trim(paramValue, "\"")
	}

	return apiHelper.ServiceResponseWrapper(
		c, controller.containerImageService.ReadArchives(requestBody, &c.Request().Host),
	)
}

// DownloadContainerImageArchive	 godoc
// @Summary      DownloadContainerImageArchive
// @Description  Download a container image archive file.
// @Tags         containerImageArchive
// @Accept       json
// @Produce      octet-stream
// @Security     Bearer
// @Param        accountId 	  path   string  true  "AccountId"
// @Param        imageId 	  path   string  true  "ImageId"
// @Success      200 file file "ContainerImageArchive"
// @Router       /v1/container/image/archive/{accountId}/{imageId}/ [get]
func (controller *ContainerImageController) ReadArchive(c echo.Context) error {
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

	readDto := dto.NewReadContainerImageArchive(accountId, imageId)

	archiveFile, err := useCase.ReadContainerImageArchive(containerImageQueryRepo, readDto)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return c.Attachment(
		archiveFile.UnixFilePath.String(),
		archiveFile.UnixFilePath.ReadFileName().String(),
	)
}

// CreateContainerImageArchive	 godoc
// @Summary      CreateContainerImageArchive
// @Description  Export a container image to a file. This is an asynchronous operation.
// @Tags         containerImageArchive
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        createContainerImageArchiveDto 	  body    dto.CreateContainerImageArchive  true  "CreateContainerImageArchiveDto"
// @Success      201 {object} object{} "ContainerImageArchiveCreationScheduled"
// @Router       /v1/container/image/archive/ [post]
func (controller *ContainerImageController) CreateArchive(c echo.Context) error {
	requestBody, err := apiHelper.ReadRequestBody(c)
	if err != nil {
		return err
	}

	return apiHelper.ServiceResponseWrapper(
		c, controller.containerImageService.CreateArchive(requestBody, true),
	)
}

type FailedArchiveFileUpload struct {
	FileName   string `json:"fileName"`
	FailReason string `json:"failReason"`
}

// ImportContainerImageArchives	godoc
// @Summary      ImportContainerImageArchives
// @Description  Import container images from archive files.
// @Tags         containerImageArchive
// @Accept       mpfd
// @Produce      json
// @Security     Bearer
// @Param        accountId		formData	string	false	"AccountId"
// @Param        archiveFiles	formData	file	true	"ArchiveFiles"
// @Success      201 string string "ContainerImageArchivesImported"
// @Success      207 {object} []FailedArchiveFileUpload "ContainerImageArchivesPartiallyImported"
// @Router       /v1/container/image/archive/import/ [post]
func (controller *ContainerImageController) ImportArchive(c echo.Context) error {
	requestBody, err := apiHelper.ReadRequestBody(c)
	if err != nil {
		return err
	}

	requiredParams := []string{"operatorAccountId", "operatorIpAddress"}
	err = serviceHelper.RequiredParamsInspector(requestBody, requiredParams)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
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
		importDto := dto.NewImportContainerImageArchive(
			accountId, archiveFile, operatorAccountId, operatorIpAddress,
		)

		_, err = useCase.ImportContainerImageArchive(
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
		c, http.StatusCreated, "ContainerImageArchivesImported",
	)
}

// DeleteContainerImageArchive	 godoc
// @Summary      DeleteContainerImageArchive
// @Description  Delete a container image archive file.
// @Tags         containerImageArchive
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        accountId 	  path   string  true  "AccountId"
// @Param        imageId 	  path   string  true  "ImageId"
// @Success      200 {object} object{} "ContainerImageArchiveDeleted"
// @Router       /v1/container/image/archive/{accountId}/{imageId}/ [delete]
func (controller *ContainerImageController) DeleteArchive(c echo.Context) error {
	requestBody := map[string]interface{}{
		"accountId":         c.Param("accountId"),
		"imageId":           c.Param("imageId"),
		"operatorAccountId": c.Get("accountId"),
		"operatorIpAddress": c.RealIP(),
	}

	return apiHelper.ServiceResponseWrapper(
		c, controller.containerImageService.DeleteArchive(requestBody),
	)
}
