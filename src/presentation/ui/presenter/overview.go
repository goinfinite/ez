package presenter

import (
	"log/slog"
	"net/http"

	"github.com/a-h/templ"
	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/entity"
	"github.com/goinfinite/ez/src/domain/valueObject"
	voHelper "github.com/goinfinite/ez/src/domain/valueObject/helper"
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
		htmlLabel := componentContainer.ContainerTaggedSummary(containerSummary, "")

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

func (presenter *OverviewPresenter) ReadContainers(c echo.Context) (
	responseDto dto.ReadContainersResponse,
) {
	containersPageNumber := uint16(0)
	if c.QueryParam("containersPageNumber") != "" {
		containersPageNumber, _ = voHelper.InterfaceToUint16(
			c.QueryParam("containersPageNumber"),
		)
	}

	containersItemsPerPage := uint16(5)
	if c.QueryParam("containersItemsPerPage") != "" {
		containersItemsPerPage, _ = voHelper.InterfaceToUint16(
			c.QueryParam("containersItemsPerPage"),
		)
	}

	containersSortByStr := "hostname"
	if c.QueryParam("containersSortBy") != "" {
		containersSortBy, err := valueObject.NewPaginationSortBy(
			c.QueryParam("containersSortBy"),
		)
		if err == nil {
			containersSortByStr = containersSortBy.String()
		}
	}

	containersSortDirectionStr := "asc"
	if c.QueryParam("containersSortDirection") != "" {
		containersSortDirection, err := valueObject.NewPaginationSortDirection(
			c.QueryParam("containersSortDirection"),
		)
		if err == nil {
			containersSortDirectionStr = containersSortDirection.String()
		}
	}

	containerService := service.NewContainerService(
		presenter.persistentDbSvc, presenter.trailDbSvc,
	)

	readContainersRequestBody := map[string]interface{}{
		"pageNumber":    containersPageNumber,
		"itemsPerPage":  containersItemsPerPage,
		"sortBy":        containersSortByStr,
		"sortDirection": containersSortDirectionStr,
		"withMetrics":   true,
	}

	if c.QueryParam("containersContainerId") != "" {
		readContainersRequestBody["containerId"] = c.QueryParam("containersContainerId")
	}

	if c.QueryParam("containersAccountId") != "" {
		readContainersRequestBody["containerAccountId"] = c.QueryParam("containersAccountId")
	}

	if c.QueryParam("containersHostname") != "" {
		readContainersRequestBody["containerHostname"] = c.QueryParam("containersHostname")
	}

	if c.QueryParam("containersStatus") != "" {
		readContainersRequestBody["containerStatus"] = c.QueryParam("containersStatus")
	}

	if c.QueryParam("containersImageId") != "" {
		readContainersRequestBody["containerImageId"] = c.QueryParam("containersImageId")
	}

	if c.QueryParam("containersImageAddress") != "" {
		readContainersRequestBody["containerImageAddress"] = c.QueryParam("containersImageAddress")
	}

	if c.QueryParam("containersImageHash") != "" {
		readContainersRequestBody["containerImageHash"] = c.QueryParam("containersImageHash")
	}

	if c.QueryParam("containersRestartPolicy") != "" {
		readContainersRequestBody["containerRestartPolicy"] = c.QueryParam("containersRestartPolicy")
	}

	if c.QueryParam("containersProfileId") != "" {
		readContainersRequestBody["containerProfileId"] = c.QueryParam("containersProfileId")
	}

	readContainersServiceOutput := containerService.Read(readContainersRequestBody)
	if readContainersServiceOutput.Status != service.Success {
		slog.Debug("ReadContainersFailure")
		return responseDto
	}

	containersResponseDto, assertOk := readContainersServiceOutput.Body.(dto.ReadContainersResponse)
	if !assertOk {
		slog.Debug("AssertContainersResponseFailure")
		return responseDto
	}

	return containersResponseDto
}

func (presenter *OverviewPresenter) Handler(c echo.Context) (err error) {
	containersResponseDto := presenter.ReadContainers(c)

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
		containersResponseDto, containerIdSummariesMap, createContainerModalDto,
	)

	return uiHelper.Render(c, pageContent, http.StatusOK)
}
