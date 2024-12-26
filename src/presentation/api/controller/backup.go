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
// @Param        backupType query  string  false  "BackupType"
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
		"jobId", "jobStatus", "accountId", "destinationId", "backupType",
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
