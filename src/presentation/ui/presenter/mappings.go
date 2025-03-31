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
	mappingService *service.MappingService
}

func NewMappingsPresenter(
	persistentDbSvc *db.PersistentDatabaseService,
	trailDbSvc *db.TrailDatabaseService,
) *MappingsPresenter {
	return &MappingsPresenter{
		mappingService: service.NewMappingService(persistentDbSvc, trailDbSvc),
	}
}

func (presenter *MappingsPresenter) Handler(c echo.Context) error {
	readMappingsServiceOutput := presenter.mappingService.Read()
	if readMappingsServiceOutput.Status != service.Success {
		return nil
	}

	mappingsList, assertOk := readMappingsServiceOutput.Body.([]entity.Mapping)
	if !assertOk {
		return nil
	}

	pageContent := page.MappingsIndex(mappingsList)
	return presenterHelper.Render(c, pageContent, http.StatusOK)
}
