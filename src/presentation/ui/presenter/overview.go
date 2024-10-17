package presenter

import (
	"log/slog"
	"net/http"

	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/entity"
	"github.com/goinfinite/ez/src/domain/valueObject"
	"github.com/goinfinite/ez/src/infra/db"
	"github.com/goinfinite/ez/src/presentation/service"
	componentContainer "github.com/goinfinite/ez/src/presentation/ui/component/container"
	componentForm "github.com/goinfinite/ez/src/presentation/ui/component/form"
	uiHelper "github.com/goinfinite/ez/src/presentation/ui/helper"
	"github.com/goinfinite/ez/src/presentation/ui/page"
	presenterHelper "github.com/goinfinite/ez/src/presentation/ui/presenter/helper"
	"github.com/labstack/echo/v4"
)

type OverviewPresenter struct {
	persistentDbSvc *db.PersistentDatabaseService
	transientDbSvc  *db.TransientDatabaseService
	trailDbSvc      *db.TrailDatabaseService
}

func NewOverviewPresenter(
	persistentDbSvc *db.PersistentDatabaseService,
	transientDbSvc *db.TransientDatabaseService,
	trailDbSvc *db.TrailDatabaseService,
) *OverviewPresenter {
	return &OverviewPresenter{
		persistentDbSvc: persistentDbSvc,
		transientDbSvc:  transientDbSvc,
		trailDbSvc:      trailDbSvc,
	}
}

func (presenter *OverviewPresenter) transformContainerImagesIntoSearchableItems() []componentForm.SearchableSelectItem {
	searchableSelectItems := []componentForm.SearchableSelectItem{}

	containerImageService := service.NewContainerImageService(
		presenter.persistentDbSvc, presenter.trailDbSvc,
	)

	readImagesServiceOutput := containerImageService.Read()
	if readImagesServiceOutput.Status != service.Success {
		slog.Debug("ReadImagesFailure")
		return nil
	}

	containerImageEntities, assertOk := readImagesServiceOutput.Body.([]entity.ContainerImage)
	if !assertOk {
		slog.Debug("AssertImagesFailure")
		return nil
	}

	for _, imageEntity := range containerImageEntities {
		searchableTextSerialized := imageEntity.JsonSerialize()
		htmlLabel := componentContainer.ImageTaggedSummary(imageEntity)

		searchableSelectItem := componentForm.SearchableSelectItem{
			Label:          imageEntity.ImageAddress.String(),
			Value:          imageEntity.ImageAddress.String(),
			SearchableText: &searchableTextSerialized,
			HtmlLabel:      &htmlLabel,
		}
		searchableSelectItems = append(searchableSelectItems, searchableSelectItem)
	}

	return searchableSelectItems
}

func (presenter *OverviewPresenter) Handler(c echo.Context) error {
	containerService := service.NewContainerService(
		presenter.persistentDbSvc, presenter.trailDbSvc,
	)

	readContainersServiceOutput := containerService.ReadWithMetrics()
	if readContainersServiceOutput.Status != service.Success {
		slog.Debug("ReadContainersFailure")
		return nil
	}

	containerEntities, assertOk := readContainersServiceOutput.Body.([]dto.ContainerWithMetrics)
	if !assertOk {
		slog.Debug("AssertContainersFailure")
		return nil
	}

	containerSummaries := presenterHelper.ReadContainerSummaries(
		presenter.persistentDbSvc, presenter.trailDbSvc,
	)

	containerIdSummariesMap := map[valueObject.ContainerId]componentContainer.ContainerSummary{}
	for _, containerSummary := range containerSummaries {
		containerIdSummariesMap[containerSummary.ContainerId] = containerSummary
	}

	marketplaceRequestBody := map[string]interface{}{
		"sortBy":       "id",
		"itemsPerPage": 100,
	}
	marketplaceService := service.NewMarketplaceService()

	readMarketplaceServiceOutput := marketplaceService.Read(marketplaceRequestBody)
	if readMarketplaceServiceOutput.Status != service.Success {
		slog.Debug("ReadMarketplaceFailure")
		return nil
	}

	marketplaceResponseDto, assertOk := readMarketplaceServiceOutput.Body.(dto.ReadMarketplaceItemsResponse)
	if !assertOk {
		slog.Debug("AssertMarketplaceFailure")
		return nil
	}

	unprocessedCreateContainerModalDto := page.CreateContainerModalUnprocessedDto{
		MarketplaceItems:              marketplaceResponseDto.Items,
		ContainerImageSearchableItems: presenter.transformContainerImagesIntoSearchableItems(),
	}

	pageContent := page.OverviewIndex(
		containerEntities, containerIdSummariesMap, unprocessedCreateContainerModalDto,
	)

	return uiHelper.Render(c, pageContent, http.StatusOK)
}
