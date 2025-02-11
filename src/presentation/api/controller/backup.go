package apiController

import (
	"net/http"
	"strings"

	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/useCase"
	"github.com/goinfinite/ez/src/domain/valueObject"
	backupInfra "github.com/goinfinite/ez/src/infra/backup"
	"github.com/goinfinite/ez/src/infra/db"
	apiHelper "github.com/goinfinite/ez/src/presentation/api/helper"
	"github.com/goinfinite/ez/src/presentation/service"
	sharedHelper "github.com/goinfinite/ez/src/presentation/shared/helper"
	"github.com/labstack/echo/v4"
)

type BackupController struct {
	backupService   *service.BackupService
	persistentDbSvc *db.PersistentDatabaseService
	trailDbSvc      *db.TrailDatabaseService
}

func NewBackupController(
	persistentDbSvc *db.PersistentDatabaseService,
	trailDbSvc *db.TrailDatabaseService,
) *BackupController {
	return &BackupController{
		backupService: service.NewBackupService(
			persistentDbSvc, trailDbSvc,
		),
		persistentDbSvc: persistentDbSvc,
		trailDbSvc:      trailDbSvc,
	}
}

// ReadBackupDestinations	 godoc
// @Summary      ReadBackupsDestinations
// @Description  List backups destinations.
// @Tags         backup
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        destinationId query  string  false  "BackupDestinationId"
// @Param        accountId query  uint  false  "BackupAccountId"
// @Param        destinationName query  string  false  "BackupDestinationName"
// @Param        destinationType query  string  false  "BackupDestinationType"
// @Param        objectStorageProvider query  string  false  "ObjectStorageProvider"
// @Param        remoteHostType query  string  false  "RemoteHostType"
// @Param        remoteHostname query  string  false  "RemoteHostname"
// @Param        createdBeforeAt query  string  false  "CreatedBeforeAt"
// @Param        createdAfterAt query  string  false  "CreatedAfterAt"
// @Param        pageNumber query  uint  false  "PageNumber (Pagination)"
// @Param        itemsPerPage query  uint  false  "ItemsPerPage (Pagination)"
// @Param        sortBy query  string  false  "SortBy (Pagination)"
// @Param        sortDirection query  string  false  "SortDirection (Pagination)"
// @Param        lastSeenId query  string  false  "LastSeenId (Pagination)"
// @Success      200 {object} dto.ReadBackupDestinationsResponse
// @Router       /v1/backup/destination/ [get]
func (controller *BackupController) ReadDestination(c echo.Context) error {
	requestBody := map[string]interface{}{}
	queryParameters := []string{
		"destinationId", "accountId", "destinationName", "destinationType",
		"objectStorageProvider", "remoteHostType", "remoteHostname",
		"createdBeforeAt", "createdAfterAt", "pageNumber", "itemsPerPage",
		"sortBy", "sortDirection", "lastSeenId",
	}
	for _, paramName := range queryParameters {
		paramValue := c.QueryParam(paramName)
		if paramValue == "" {
			continue
		}

		requestBody[paramName] = strings.Trim(paramValue, "\"")
	}

	return apiHelper.ServiceResponseWrapper(
		c, controller.backupService.ReadDestination(requestBody),
	)
}

// CreateBackupDestination	 godoc
// @Summary      CreateBackupDestination
// @Description  Create a backup destination.
// @Tags         backup
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        createBackupDestinationDto 	  body    dto.CreateBackupDestinationRequest  true  "CreateBackupDestination"
// @Success      201 {object} dto.CreateBackupDestinationResponse
// @Failure      500 {string} string "CreateBackupDestinationInfraError"
// @Router       /v1/backup/destination/ [post]
func (controller *BackupController) CreateDestination(c echo.Context) error {
	requestBody, err := apiHelper.ReadRequestBody(c)
	if err != nil {
		return err
	}

	return apiHelper.ServiceResponseWrapper(
		c, controller.backupService.CreateDestination(requestBody),
	)
}

// UpdateBackupDestination	 godoc
// @Summary      UpdateBackupDestination
// @Description  Update a backup destination.
// @Tags         backup
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        updateBackupDestinationDto 	  body    dto.UpdateBackupDestination  true  "UpdateBackupDestination"
// @Success      200 {object} object{} "BackupDestinationUpdated"
// @Router       /v1/backup/destination/ [put]
func (controller *BackupController) UpdateDestination(c echo.Context) error {
	requestBody, err := apiHelper.ReadRequestBody(c)
	if err != nil {
		return err
	}

	return apiHelper.ServiceResponseWrapper(
		c, controller.backupService.UpdateDestination(requestBody),
	)
}

// DeleteBackupDestination	 godoc
// @Summary      DeleteBackupDestination
// @Description  Delete a backup destination.
// @Tags         backup
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        accountId 	  path   string  true  "AccountId"
// @Param        destinationId path  string  true  "BackupDestinationId"
// @Success      200 {object} object{} "BackupDestinationDeleted"
// @Router       /v1/backup/destination/{accountId}/{destinationId} [delete]
func (controller *BackupController) DeleteDestination(c echo.Context) error {
	requestBody := map[string]interface{}{
		"accountId":         c.Param("accountId"),
		"destinationId":     c.Param("destinationId"),
		"operatorAccountId": c.Get("accountId"),
		"operatorIpAddress": c.RealIP(),
	}

	return apiHelper.ServiceResponseWrapper(
		c, controller.backupService.DeleteDestination(requestBody),
	)
}

// ReadBackupJobs	 godoc
// @Summary      ReadBackupJobs
// @Description  List backup jobs.
// @Tags         backup
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        jobId query  string  false  "BackupJobId"
// @Param        jobStatus query  bool  false  "BackupJobStatus"
// @Param        accountId query  uint  false  "BackupAccountId"
// @Param        destinationId query  string  false  "BackupDestinationId"
// @Param        retentionStrategy query  string  false  "RetentionStrategy"
// @Param        archiveCompressionFormat query  string  false  "ArchiveCompressionFormat"
// @Param        lastRunStatus query  string  false  "LastRunStatus"
// @Param        lastRunBeforeAt query  string  false  "LastRunBeforeAt"
// @Param        lastRunAfterAt query  string  false  "LastRunAfterAt"
// @Param        nextRunBeforeAt query  string  false  "NextRunBeforeAt"
// @Param        nextRunAfterAt query  string  false  "NextRunAfterAt"
// @Param        createdBeforeAt query  string  false  "CreatedBeforeAt"
// @Param        createdAfterAt query  string  false  "CreatedAfterAt"
// @Param        pageNumber query  uint  false  "PageNumber (Pagination)"
// @Param        itemsPerPage query  uint  false  "ItemsPerPage (Pagination)"
// @Param        sortBy query  string  false  "SortBy (Pagination)"
// @Param        sortDirection query  string  false  "SortDirection (Pagination)"
// @Param        lastSeenId query  string  false  "LastSeenId (Pagination)"
// @Success      200 {object} dto.ReadBackupJobsResponse
// @Router       /v1/backup/job/ [get]
func (controller *BackupController) ReadJob(c echo.Context) error {
	requestBody := map[string]interface{}{}
	queryParameters := []string{
		"jobId", "jobStatus", "accountId", "destinationId", "retentionStrategy",
		"archiveCompressionFormat", "lastRunStatus", "lastRunBeforeAt",
		"lastRunAfterAt", "nextRunBeforeAt", "nextRunAfterAt",
		"createdBeforeAt", "createdAfterAt", "pageNumber", "itemsPerPage",
		"sortBy", "sortDirection", "lastSeenId",
	}
	for _, paramName := range queryParameters {
		paramValue := c.QueryParam(paramName)
		if paramValue == "" {
			continue
		}

		requestBody[paramName] = strings.Trim(paramValue, "\"")
	}

	return apiHelper.ServiceResponseWrapper(
		c, controller.backupService.ReadJob(requestBody),
	)
}

// CreateBackupJob	 godoc
// @Summary      CreateBackupJob
// @Description  Create a backup destination.
// @Tags         backup
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        createBackupJobDto 	  body    dto.CreateBackupJob  true  "CreateBackupJob"
// @Success      201 {object} object{} "BackupJobCreated"
// @Router       /v1/backup/job/ [post]
func (controller *BackupController) CreateJob(c echo.Context) error {
	requestBody, err := apiHelper.ReadRequestBody(c)
	if err != nil {
		return err
	}

	if requestBody["destinationIds"] != nil {
		requestBody["destinationIds"] = sharedHelper.StringSliceValueObjectParser(
			requestBody["destinationIds"], valueObject.NewBackupDestinationId,
		)
	}

	return apiHelper.ServiceResponseWrapper(
		c, controller.backupService.CreateJob(requestBody),
	)
}

// UpdateBackupJob	 godoc
// @Summary      UpdateBackupJob
// @Description  Update a backup job.
// @Tags         backup
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        updateBackupJobDto 	  body    dto.UpdateBackupJob  true  "UpdateBackupJob"
// @Success      200 {object} object{} "BackupJobUpdated"
// @Router       /v1/backup/job/ [put]
func (controller *BackupController) UpdateJob(c echo.Context) error {
	requestBody, err := apiHelper.ReadRequestBody(c)
	if err != nil {
		return err
	}

	return apiHelper.ServiceResponseWrapper(
		c, controller.backupService.UpdateJob(requestBody),
	)
}

// DeleteBackupJob	 godoc
// @Summary      DeleteBackupJob
// @Description  Delete a backup job.
// @Tags         backup
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        accountId 	  path   string  true  "AccountId"
// @Param        jobId path  string  true  "BackupJobId"
// @Success      200 {object} object{} "BackupJobDeleted"
// @Router       /v1/backup/job/{accountId}/{jobId} [delete]
func (controller *BackupController) DeleteJob(c echo.Context) error {
	requestBody := map[string]interface{}{
		"accountId":         c.Param("accountId"),
		"jobId":             c.Param("jobId"),
		"operatorAccountId": c.Get("accountId"),
		"operatorIpAddress": c.RealIP(),
	}

	return apiHelper.ServiceResponseWrapper(
		c, controller.backupService.DeleteJob(requestBody),
	)
}

// RunBackupJob	 godoc
// @Summary      RunBackupJob
// @Description  Run a backup job.
// @Tags         backup
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        accountId 	  path   string  true  "AccountId"
// @Param        jobId path  string  true  "BackupJobId"
// @Success      201 {object} object{} "BackupTaskCreated"
// @Router       /v1/backup/job/{accountId}/{jobId}/run/ [post]
func (controller *BackupController) RunJob(c echo.Context) error {
	requestBody := map[string]interface{}{
		"accountId":         c.Param("accountId"),
		"jobId":             c.Param("jobId"),
		"operatorAccountId": c.Get("accountId"),
		"operatorIpAddress": c.RealIP(),
	}

	return apiHelper.ServiceResponseWrapper(
		c, controller.backupService.RunJob(requestBody),
	)
}

// ReadBackupTasks	 godoc
// @Summary      ReadBackupTasks
// @Description  List backup tasks.
// @Tags         backup
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        taskId query  string  false  "BackupTaskId"
// @Param        accountId query  uint  false  "BackupAccountId"
// @Param        jobId query  string  false  "BackupJobId"
// @Param        destinationId query  string  false  "BackupDestinationId"
// @Param        taskStatus query  string  false  "BackupTaskStatus"
// @Param        retentionStrategy query  string  false  "RetentionStrategy"
// @Param        containerId query  string  false  "ContainerId"
// @Param        startedBeforeAt query  string  false  "StartedBeforeAt"
// @Param        startedAfterAt query  string  false  "StartedAfterAt"
// @Param        finishedBeforeAt query  string  false  "FinishedBeforeAt"
// @Param        finishedAfterAt query  string  false  "FinishedAfterAt"
// @Param        createdBeforeAt query  string  false  "CreatedBeforeAt"
// @Param        createdAfterAt query  string  false  "CreatedAfterAt"
// @Param        pageNumber query  uint  false  "PageNumber (Pagination)"
// @Param        itemsPerPage query  uint  false  "ItemsPerPage (Pagination)"
// @Param        sortBy query  string  false  "SortBy (Pagination)"
// @Param        sortDirection query  string  false  "SortDirection (Pagination)"
// @Param        lastSeenId query  string  false  "LastSeenId (Pagination)"
// @Success      200 {object} dto.ReadBackupTasksResponse
// @Router       /v1/backup/task/ [get]
func (controller *BackupController) ReadTask(c echo.Context) error {
	requestBody := map[string]interface{}{}
	queryParameters := []string{
		"taskId", "accountId", "jobId", "destinationId", "taskStatus",
		"retentionStrategy", "containerId", "startedBeforeAt", "startedAfterAt",
		"finishedBeforeAt", "finishedAfterAt", "createdBeforeAt", "createdAfterAt",
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
		c, controller.backupService.ReadTask(requestBody),
	)
}

// DeleteBackupTask	 godoc
// @Summary      DeleteBackupTask
// @Description  Delete a backup task and its files if "shouldDiscardFiles" is true.
// @Tags         backup
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        taskId path  string  true  "BackupTaskId"
// @Param        shouldDiscardFiles query  string  false  "ShouldDiscardFiles (bool)"
// @Success      200 {object} object{} "BackupTaskDeleted"
// @Router       /v1/backup/task/{taskId}/ [delete]
func (controller *BackupController) DeleteTask(c echo.Context) error {
	requestBody := map[string]interface{}{
		"taskId":            c.Param("taskId"),
		"operatorAccountId": c.Get("accountId"),
		"operatorIpAddress": c.RealIP(),
	}
	if c.QueryParam("shouldDiscardFiles") != "" {
		requestBody["shouldDiscardFiles"] = c.QueryParam("shouldDiscardFiles")
	}

	return apiHelper.ServiceResponseWrapper(
		c, controller.backupService.DeleteTask(requestBody),
	)
}

// ReadBackupTaskArchives	 godoc
// @Summary      ReadBackupTaskArchives
// @Description  List backup tasks archives.
// @Tags         backup
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        archiveId query  string  false  "BackupTaskArchiveId"
// @Param        accountId query  uint  false  "BackupAccountId"
// @Param        taskId query  string  false  "BackupTaskId"
// @Param        createdBeforeAt query  string  false  "CreatedBeforeAt"
// @Param        createdAfterAt query  string  false  "CreatedAfterAt"
// @Param        pageNumber query  uint  false  "PageNumber (Pagination)"
// @Param        itemsPerPage query  uint  false  "ItemsPerPage (Pagination)"
// @Param        sortBy query  string  false  "SortBy (Pagination)"
// @Param        sortDirection query  string  false  "SortDirection (Pagination)"
// @Param        lastSeenId query  string  false  "LastSeenId (Pagination)"
// @Success      200 {object} dto.ReadBackupTaskArchivesResponse
// @Router       /v1/backup/task/archive/ [get]
func (controller *BackupController) ReadTaskArchives(c echo.Context) error {
	requestBody := map[string]interface{}{}
	queryParameters := []string{
		"archiveId", "accountId", "taskId", "createdBeforeAt", "createdAfterAt",
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
		c, controller.backupService.ReadTaskArchive(requestBody, &c.Request().Host),
	)
}

// DownloadBackupTaskArchiveFile	 godoc
// @Summary      DownloadBackupTaskArchiveFile
// @Description  Download a backup task archive file.
// @Tags         backup
// @Accept       json
// @Produce      octet-stream
// @Security     Bearer
// @Param        archiveId 	  path   string  true  "ArchiveId"
// @Success      200 file file "BackupTaskArchiveFile"
// @Router       /v1/backup/task/archive/{archiveId}/ [get]
func (controller *BackupController) ReadTaskArchive(c echo.Context) error {
	if c.Param("archiveId") == "" {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, "ArchiveIdRequired")
	}
	archiveId, err := valueObject.NewBackupTaskArchiveId(c.Param("archiveId"))
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
	}

	backupQueryRepo := backupInfra.NewBackupQueryRepo(controller.persistentDbSvc)
	requestDto := dto.ReadBackupTaskArchivesRequest{
		ArchiveId: &archiveId,
	}

	archiveFile, err := useCase.ReadBackupTaskArchive(backupQueryRepo, requestDto)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return c.Attachment(
		archiveFile.UnixFilePath.String(),
		archiveFile.UnixFilePath.ReadFileName().String(),
	)
}

// CreateBackupTaskArchive	godoc
// @Summary      CreateBackupTaskArchive
// @Description  Schedule a backup task archive creation.
// @Tags         backup
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        createBackupTaskArchiveDto	body	dto.CreateBackupTaskArchive	true	"CreateBackupTaskArchive"
// @Success      201	{object}	object{}	"BackupTaskArchiveCreationScheduled"
// @Router       /v1/backup/task/archive/	[post]
func (controller *BackupController) CreateTaskArchive(c echo.Context) error {
	requestBody, err := apiHelper.ReadRequestBody(c)
	if err != nil {
		return err
	}

	if requestBody["containerAccountIds"] == nil {
		if requestBody["containerAccountId"] != nil {
			requestBody["containerAccountIds"] = requestBody["containerAccountId"]
		}
	}

	if requestBody["containerAccountIds"] != nil {
		requestBody["containerAccountIds"] = sharedHelper.StringSliceValueObjectParser(
			requestBody["containerAccountIds"], valueObject.NewAccountId,
		)
	}

	if requestBody["containerIds"] == nil {
		if requestBody["containerId"] != nil {
			requestBody["containerIds"] = requestBody["containerId"]
		}
	}

	if requestBody["containerIds"] != nil {
		requestBody["containerIds"] = sharedHelper.StringSliceValueObjectParser(
			requestBody["containerIds"], valueObject.NewContainerId,
		)
	}

	if requestBody["exceptContainerAccountIds"] == nil {
		if requestBody["exceptContainerAccountId"] != nil {
			requestBody["exceptContainerAccountIds"] = requestBody["exceptContainerAccountId"]
		}
	}

	if requestBody["exceptContainerAccountIds"] != nil {
		requestBody["exceptContainerAccountIds"] = sharedHelper.StringSliceValueObjectParser(
			requestBody["exceptContainerAccountIds"], valueObject.NewAccountId,
		)
	}

	if requestBody["exceptContainerIds"] == nil {
		if requestBody["exceptContainerId"] != nil {
			requestBody["exceptContainerIds"] = requestBody["exceptContainerId"]
		}
	}

	if requestBody["exceptContainerIds"] != nil {
		requestBody["exceptContainerIds"] = sharedHelper.StringSliceValueObjectParser(
			requestBody["exceptContainerIds"], valueObject.NewContainerId,
		)
	}

	return apiHelper.ServiceResponseWrapper(
		c, controller.backupService.CreateTaskArchive(requestBody, true),
	)
}

// DeleteBackupTaskArchive	 godoc
// @Summary      DeleteBackupTaskArchive
// @Description  Delete a backup task archive.
// @Tags         backup
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        archiveId path  string  true  "BackupTaskArchiveId"
// @Success      200 {object} object{} "BackupTaskArchiveDeleted"
// @Router       /v1/backup/task/archive/{archiveId}/ [delete]
func (controller *BackupController) DeleteTaskArchive(c echo.Context) error {
	if c.Param("archiveId") == "" {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, "ArchiveIdRequired")
	}
	requestBody := map[string]interface{}{
		"archiveId":         c.Param("archiveId"),
		"operatorAccountId": c.Get("accountId"),
		"operatorIpAddress": c.RealIP(),
	}

	return apiHelper.ServiceResponseWrapper(
		c, controller.backupService.DeleteTaskArchive(requestBody),
	)
}
