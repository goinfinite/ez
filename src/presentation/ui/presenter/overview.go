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

func (presenter *OverviewPresenter) readContainerImageSearchableItems() []componentForm.SearchableSelectItem {
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

func (presenter *OverviewPresenter) readContainerProfileSearchableItems() []componentForm.SearchableSelectItem {
	searchableSelectItems := []componentForm.SearchableSelectItem{}

	containerProfileService := service.NewContainerProfileService(
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

func (presenter *OverviewPresenter) readAccountSelectLabelValuePairs() []componentForm.SelectLabelValuePair {
	selectLabelValuePairs := []componentForm.SelectLabelValuePair{}

	accountService := service.NewAccountService(
		presenter.persistentDbSvc, presenter.trailDbSvc,
	)

	readAccountsServiceOutput := accountService.Read()
	if readAccountsServiceOutput.Status != service.Success {
		slog.Debug("ReadAccountsFailure")
		return nil
	}

	accountEntities, assertOk := readAccountsServiceOutput.Body.([]entity.Account)
	if !assertOk {
		slog.Debug("AssertAccountsFailure")
		return nil
	}

	for _, accountEntity := range accountEntities {
		selectLabelValuePair := componentForm.SelectLabelValuePair{
			Label: accountEntity.Username.String(),
			Value: accountEntity.Id.String(),
		}
		selectLabelValuePairs = append(selectLabelValuePairs, selectLabelValuePair)
	}

	return selectLabelValuePairs
}

func (presenter *OverviewPresenter) transformContainerSummariesIntoSearchableItems(
	containerSummaries []componentContainer.ContainerSummary,
) []componentForm.SearchableSelectItem {
	searchableSelectItems := []componentForm.SearchableSelectItem{}

	for _, containerSummary := range containerSummaries {
		searchableTextSerialized := containerSummary.JsonSerialize()
		htmlLabel := componentContainer.ContainerTaggedSummary(containerSummary)

		searchableSelectItem := componentForm.SearchableSelectItem{
			Label:          containerSummary.Hostname.String(),
			Value:          containerSummary.ContainerId.String(),
			SearchableText: &searchableTextSerialized,
			HtmlLabel:      &htmlLabel,
		}
		searchableSelectItems = append(searchableSelectItems, searchableSelectItem)
	}

	return searchableSelectItems
}

func (presenter *OverviewPresenter) ReadCreateContainerModalDto(
	containerSummaries []componentContainer.ContainerSummary,
) (createDto page.CreateContainerModalDto) {
	marketplaceRequestBody := map[string]interface{}{
		"sortBy":       "id",
		"itemsPerPage": 100,
	}
	marketplaceService := service.NewMarketplaceService()

	readMarketplaceServiceOutput := marketplaceService.Read(marketplaceRequestBody)
	if readMarketplaceServiceOutput.Status != service.Success {
		slog.Debug("ReadMarketplaceFailure")
		return createDto
	}

	marketplaceResponseDto, assertOk := readMarketplaceServiceOutput.Body.(dto.ReadMarketplaceItemsResponse)
	if !assertOk {
		slog.Debug("AssertMarketplaceFailure")
		return createDto
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

	return page.CreateContainerModalDto{
		AppMarketplaceCarouselItems:       appCarouselItems,
		FrameworkMarketplaceCarouselItems: frameworkCarouselItems,
		StackMarketplaceCarouselItems:     stackCarouselItems,
		ContainerImageSearchableItems:     presenter.readContainerImageSearchableItems(),
		ContainerProfileSearchableItems:   presenter.readContainerProfileSearchableItems(),
		AccountSelectLabelValuePairs:      presenter.readAccountSelectLabelValuePairs(),
		ContainerSummarySearchableItems: presenter.transformContainerSummariesIntoSearchableItems(
			containerSummaries,
		),
	}
}

func (presenter *OverviewPresenter) Handler(c echo.Context) error {
	containerService := service.NewContainerService(
		presenter.persistentDbSvc, presenter.trailDbSvc,
	)

	readContainersRequestBody := map[string]interface{}{
		"withMetrics": true,
	}

	readContainersServiceOutput := containerService.Read(readContainersRequestBody)
	if readContainersServiceOutput.Status != service.Success {
		slog.Debug("ReadContainersFailure")
		return nil
	}

	containersResponseDto, assertOk := readContainersServiceOutput.Body.(dto.ReadContainersResponse)
	if !assertOk {
		slog.Debug("AssertContainersResponseFailure")
		return nil
	}

	containerSummaries := presenterHelper.ReadContainerSummaries(
		presenter.persistentDbSvc, presenter.trailDbSvc,
	)

	containerIdSummariesMap := map[valueObject.ContainerId]componentContainer.ContainerSummary{}
	for _, containerSummary := range containerSummaries {
		containerIdSummariesMap[containerSummary.ContainerId] = containerSummary
	}

	createContainerModalDto := presenter.ReadCreateContainerModalDto(
		containerSummaries,
	)

	pageContent := page.OverviewIndex(
		containersResponseDto.ContainersWithMetrics, containerIdSummariesMap,
		createContainerModalDto,
	)

	return uiHelper.Render(c, pageContent, http.StatusOK)
}
