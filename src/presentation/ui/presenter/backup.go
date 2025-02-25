package presenter

import (
	"net/http"

	"github.com/goinfinite/ez/src/infra/db"
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

func (presenter *BackupPresenter) Handler(c echo.Context) (err error) {
	pageContent := page.BackupIndex(
		page.BackupTasksInputDto{},
		page.BackupJobsInputDto{},
		page.BackupDestinationsInputDto{},
	)

	return uiHelper.Render(c, pageContent, http.StatusOK)
}
