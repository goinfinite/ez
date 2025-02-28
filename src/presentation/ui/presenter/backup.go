package presenter

import (
	"log/slog"
	"maps"
	"net/http"

	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/infra/db"
	"github.com/goinfinite/ez/src/presentation/service"
	uiHelper "github.com/goinfinite/ez/src/presentation/ui/helper"
	"github.com/goinfinite/ez/src/presentation/ui/page"
	"github.com/labstack/echo/v4"
)

type BackupPresenter struct {
	persistentDbSvc *db.PersistentDatabaseService
	transientDbSvc  *db.TransientDatabaseService
	trailDbSvc      *db.TrailDatabaseService
}

func NewBackupPresenter(
	persistentDbSvc *db.PersistentDatabaseService,
	transientDbSvc *db.TransientDatabaseService,
	trailDbSvc *db.TrailDatabaseService,
) *BackupPresenter {
	return &BackupPresenter{
		persistentDbSvc: persistentDbSvc,
		transientDbSvc:  transientDbSvc,
		trailDbSvc:      trailDbSvc,
	}
}

func (presenter *BackupPresenter) ReadTasks(
	echoContext echo.Context,
	backupService *service.BackupService,
) (tasksResponseDto dto.ReadBackupTasksResponse) {
	paginationMap := uiHelper.PaginationParser(echoContext, "task", "id")
	requestParamsMap := uiHelper.ReadRequestParser(
		echoContext, "task", dto.ReadBackupTasksRequest{},
	)
	serviceRequestBody := paginationMap
	maps.Copy(serviceRequestBody, requestParamsMap)

	readBackupTasksServiceOutput := backupService.ReadTask(serviceRequestBody)
	if readBackupTasksServiceOutput.Status != service.Success {
		slog.Debug("ReadBackupTasksFailure", slog.Any("serviceOutput", readBackupTasksServiceOutput))
		return tasksResponseDto
	}

	var assertOk bool
	tasksResponseDto, assertOk = readBackupTasksServiceOutput.Body.(dto.ReadBackupTasksResponse)
	if !assertOk {
		slog.Debug("AssertBackupTasksResponseFailure")
		return tasksResponseDto
	}

	return tasksResponseDto
}

func (presenter *BackupPresenter) ReadTaskArchives(
	echoContext echo.Context,
	backupService *service.BackupService,
) (archivesResponseDto dto.ReadBackupTaskArchivesResponse) {
	paginationMap := uiHelper.PaginationParser(echoContext, "archive", "id")
	requestParamsMap := uiHelper.ReadRequestParser(
		echoContext, "archive", dto.ReadBackupTaskArchivesRequest{},
	)
	serviceRequestBody := paginationMap
	maps.Copy(serviceRequestBody, requestParamsMap)

	readBackupTaskArchivesServiceOutput := backupService.ReadTaskArchive(
		serviceRequestBody, &echoContext.Request().Host,
	)
	if readBackupTaskArchivesServiceOutput.Status != service.Success {
		slog.Debug(
			"ReadBackupTaskArchivesFailure",
			slog.Any("serviceOutput", readBackupTaskArchivesServiceOutput),
		)
		return archivesResponseDto
	}

	var assertOk bool
	archivesResponseDto, assertOk = readBackupTaskArchivesServiceOutput.Body.(dto.ReadBackupTaskArchivesResponse)
	if !assertOk {
		slog.Debug("AssertBackupTaskArchivesResponseFailure")
		return archivesResponseDto
	}

	return archivesResponseDto
}

func (presenter *BackupPresenter) ReadJobs(
	echoContext echo.Context,
	backupService *service.BackupService,
) (jobsResponseDto dto.ReadBackupJobsResponse) {
	paginationMap := uiHelper.PaginationParser(echoContext, "job", "id")
	requestParamsMap := uiHelper.ReadRequestParser(
		echoContext, "job", dto.ReadBackupJobsRequest{},
	)
	serviceRequestBody := paginationMap
	maps.Copy(serviceRequestBody, requestParamsMap)

	readBackupJobsServiceOutput := backupService.ReadJob(serviceRequestBody)
	if readBackupJobsServiceOutput.Status != service.Success {
		slog.Debug("ReadBackupJobsFailure", slog.Any("serviceOutput", readBackupJobsServiceOutput))
		return jobsResponseDto
	}

	var assertOk bool
	jobsResponseDto, assertOk = readBackupJobsServiceOutput.Body.(dto.ReadBackupJobsResponse)
	if !assertOk {
		slog.Debug("AssertBackupJobsResponseFailure")
		return jobsResponseDto
	}

	return jobsResponseDto
}

func (presenter *BackupPresenter) Handler(c echo.Context) (err error) {
	backupService := service.NewBackupService(
		presenter.persistentDbSvc, presenter.trailDbSvc,
	)

	tasksResponseDto := presenter.ReadTasks(c, backupService)
	archivesResponseDto := presenter.ReadTaskArchives(c, backupService)
	jobsResponseDto := presenter.ReadJobs(c, backupService)

	pageContent := page.BackupIndex(
		tasksResponseDto,
		archivesResponseDto,
		jobsResponseDto,
		dto.ReadBackupDestinationsResponse{},
	)

	return uiHelper.Render(c, pageContent, http.StatusOK)
}
