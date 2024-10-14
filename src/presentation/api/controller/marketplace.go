package apiController

import (
	"time"

	"github.com/goinfinite/ez/src/domain/useCase"
	marketplaceInfra "github.com/goinfinite/ez/src/infra/marketplace"
	apiHelper "github.com/goinfinite/ez/src/presentation/api/helper"
	"github.com/goinfinite/ez/src/presentation/service"
	"github.com/labstack/echo/v4"
)

type MarketplaceController struct {
}

func NewMarketplaceController() *MarketplaceController {
	return &MarketplaceController{}
}

// ReadMarketplaceItems	 godoc
// @Summary      ReadMarketplaceItems
// @Description  List marketplace items.
// @Tags         marketplace
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        itemSlug query  string  false  "Slug"
// @Param        itemName query  string  false  "Name"
// @Param        itemType query  string  false  "Type"
// @Param        pageNumber query  uint  false  "PageNumber (Pagination)"
// @Param        itemsPerPage query  uint  false  "ItemsPerPage (Pagination)"
// @Param        sortBy query  string  false  "SortBy (Pagination)"
// @Param        sortDirection query  string  false  "SortDirection (Pagination)"
// @Param        lastSeenId query  string  false  "LastSeenId (Pagination)"
// @Success      200 {object} dto.ReadMarketplaceItemsResponse
// @Router       /v1/marketplace/ [get]
func (controller *MarketplaceController) Read(c echo.Context) error {
	requestBody := map[string]interface{}{}
	queryParameters := []string{
		"itemSlug", "itemName", "itemType",
		"pageNumber", "itemsPerPage", "sortBy", "sortDirection", "lastSeenId",
	}
	for _, paramName := range queryParameters {
		paramValue := c.QueryParam(paramName)
		if paramValue == "" {
			continue
		}

		requestBody[paramName] = paramValue
	}

	marketplaceService := service.NewMarketplaceService()
	return apiHelper.ServiceResponseWrapper(
		c, marketplaceService.Read(requestBody),
	)
}

func (controller *MarketplaceController) RefreshMarketplace() {
	taskInterval := time.Duration(24) * time.Hour
	timer := time.NewTicker(taskInterval)
	defer timer.Stop()

	marketplaceCmdRepo := marketplaceInfra.NewMarketplaceCmdRepo()
	for range timer.C {
		useCase.RefreshMarketplace(marketplaceCmdRepo)
	}
}
