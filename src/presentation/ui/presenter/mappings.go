package presenter

import (
	"github.com/labstack/echo/v4"

	"net/http"

	"github.com/goinfinite/ez/src/domain/entity"
	"github.com/goinfinite/ez/src/infra/db"
	"github.com/goinfinite/ez/src/presentation/service"
	"github.com/goinfinite/ez/src/presentation/ui/page"
	presenterHelper "github.com/goinfinite/ez/src/presentation/ui/presenter/helper"
)

type MappingsPresenter struct {
	persistentDbSvc *db.PersistentDatabaseService
	trailDbSvc      *db.TrailDatabaseService
}

func NewMappingsPresenter(
	persistentDbSvc *db.PersistentDatabaseService,
	trailDbSvc *db.TrailDatabaseService,
) *MappingsPresenter {
	return &MappingsPresenter{
		persistentDbSvc: persistentDbSvc,
		trailDbSvc:      trailDbSvc,
	}
}

func (presenter *MappingsPresenter) Handler(c echo.Context) error {
	readMappingsServiceOutput := service.
		NewMappingService(presenter.persistentDbSvc, presenter.trailDbSvc).Read()
	if readMappingsServiceOutput.Status != service.Success {
		return nil
	}

	mappingsList, assertOk := readMappingsServiceOutput.Body.([]entity.Mapping)
	if !assertOk {
		return nil
	}

	accountSelectPairs := presenterHelper.ReadAccountSelectLabelValuePairs(
		presenter.persistentDbSvc, presenter.trailDbSvc,
	)

	pageContent := page.MappingsIndex(mappingsList, accountSelectPairs)
	return presenterHelper.Render(c, pageContent, http.StatusOK)
}
