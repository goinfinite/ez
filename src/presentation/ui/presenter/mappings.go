package presenter

import (
	"github.com/labstack/echo/v4"

	"net/http"

	"github.com/goinfinite/ez/src/domain/entity"
	"github.com/goinfinite/ez/src/infra/db"
	"github.com/goinfinite/ez/src/presentation/service"
	componentForm "github.com/goinfinite/ez/src/presentation/ui/component/form"
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

func (presenter *MappingsPresenter) readMappingsSelectLabelValuePairs(
	mappingsList []entity.Mapping,
) []componentForm.SelectLabelValuePair {
	selectLabelValuePairs := []componentForm.SelectLabelValuePair{}

	defaultHostname := "*"
	for _, mappingEntity := range mappingsList {
		mappingHostname := defaultHostname
		if mappingEntity.Hostname != nil {
			mappingHostname = mappingEntity.Hostname.String()
		}

		selectLabelValuePair := componentForm.SelectLabelValuePair{
			Label: mappingHostname + " (#" + mappingEntity.Id.String() + ")",
			Value: mappingEntity.Id.String(),
		}
		selectLabelValuePairs = append(selectLabelValuePairs, selectLabelValuePair)
	}

	return selectLabelValuePairs
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

	mappingsSelectPairs := presenter.readMappingsSelectLabelValuePairs(mappingsList)

	accountsSelectPairs := presenterHelper.ReadAccountSelectLabelValuePairs(
		presenter.persistentDbSvc, presenter.trailDbSvc,
	)

	containersSelectPairs := presenterHelper.ReadContainerSelectLabelValuePairs(
		presenter.persistentDbSvc, presenter.trailDbSvc,
	)

	pageContent := page.MappingsIndex(
		mappingsList, mappingsSelectPairs, accountsSelectPairs, containersSelectPairs,
	)
	return presenterHelper.Render(c, pageContent, http.StatusOK)
}
