package presenter

import (
	"log/slog"
	"net/http"

	"github.com/a-h/templ"
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

func (presenter *OverviewPresenter) transformContainerProfilesIntoIntoSearchableItems() []componentForm.SearchableSelectItem {
	searchableSelectItems := []componentForm.SearchableSelectItem{}

	containerProfileService := service.NewContainerService(
		presenter.persistentDbSvc, presenter.trailDbSvc,
	)

	readProfilesServiceOutput := containerProfileService.Read()
	if readProfilesServiceOutput.Status != service.Success {
		slog.Debug("ReadContainerProfileFailure")
		return nil
	}

	profileEntities, assertOk := readProfilesServiceOutput.Body.([]entity.ContainerProfile)
	if !assertOk {
		slog.Debug("AssertContainerProfileFailure")
		return nil
	}

	for _, profileEntity := range profileEntities {
		searchableTextSerialized := profileEntity.JsonSerialize()
		htmlLabel := componentContainer.ProfileTaggedSummary(profileEntity)

		searchableSelectItem := componentForm.SearchableSelectItem{
			Label:          profileEntity.Name.String(),
			Value:          profileEntity.Id.String(),
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

	var appCarouselItems, frameworkCarouselItems, stackCarouselItems []templ.Component
	for _, itemEntity := range marketplaceResponseDto.Items {
		switch itemEntity.Type.String() {
		case "app":
			appCarouselItems = append(
				appCarouselItems, page.MarketplaceCarouselItem(itemEntity),
			)
		case "framework":
			frameworkCarouselItems = append(
				frameworkCarouselItems, page.MarketplaceCarouselItem(itemEntity),
			)
		case "stack":
			stackCarouselItems = append(
				stackCarouselItems, page.MarketplaceCarouselItem(itemEntity),
			)
		}
	}

	createContainerModalDto := page.CreateContainerModalDto{
		AppMarketplaceCarouselItems:       appCarouselItems,
		FrameworkMarketplaceCarouselItems: frameworkCarouselItems,
		StackMarketplaceCarouselItems:     stackCarouselItems,
		ContainerImageSearchableItems:     presenter.transformContainerImagesIntoSearchableItems(),
		ContainerProfileSearchableItems:   presenter.transformContainerProfilesIntoIntoSearchableItems(),
	}

	pageContent := page.OverviewIndex(
		containerEntities, containerIdSummariesMap, createContainerModalDto,
	)

	return uiHelper.Render(c, pageContent, http.StatusOK)
}
