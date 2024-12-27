package apiController

import (
	"strings"

	"github.com/goinfinite/ez/src/infra/db"
	apiHelper "github.com/goinfinite/ez/src/presentation/api/helper"
	"github.com/goinfinite/ez/src/presentation/service"
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
// @Param        createBackupDestinationDto 	  body    dto.CreateBackupDestination  true  "CreateBackupDestination"
// @Success      201 {object} object{} "BackupDestinationCreated"
// @Router       /v1/backup/destination/ [post]
func (controller *BackupController) CreateDestination(c echo.Context) error {
	requestBody, err := apiHelper.ReadRequestBody(c)
	if err != nil {
		return err
	}

	if requestBody["accountId"] == nil {
		requestBody["accountId"] = requestBody["operatorAccountId"]
	}

	return apiHelper.ServiceResponseWrapper(
		c, controller.backupService.CreateDestination(requestBody),
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
