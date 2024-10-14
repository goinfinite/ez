package apiController

import (
	"errors"
	"net/http"
	"time"

	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/useCase"
	"github.com/goinfinite/ez/src/domain/valueObject"
	voHelper "github.com/goinfinite/ez/src/domain/valueObject/helper"
	"github.com/goinfinite/ez/src/infra"
	apiHelper "github.com/goinfinite/ez/src/presentation/api/helper"
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

	var itemSlugPtr *valueObject.MarketplaceItemSlug
	if requestBody["itemSlug"] != nil {
		itemSlug, err := valueObject.NewMarketplaceItemSlug(requestBody["itemSlug"])
		if err != nil {
			return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
		}
		itemSlugPtr = &itemSlug
	}

	var itemNamePtr *valueObject.MarketplaceItemName
	if requestBody["itemName"] != nil {
		itemName, err := valueObject.NewMarketplaceItemName(requestBody["itemName"])
		if err != nil {
			return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
		}
		itemNamePtr = &itemName
	}

	var itemTypePtr *valueObject.MarketplaceItemType
	if requestBody["itemType"] != nil {
		itemType, err := valueObject.NewMarketplaceItemType(requestBody["itemType"])
		if err != nil {
			return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
		}
		itemTypePtr = &itemType
	}

	paginationDto := useCase.MarketplaceDefaultPagination
	if requestBody["pageNumber"] != nil {
		pageNumber, err := voHelper.InterfaceToUint32(requestBody["pageNumber"])
		if err != nil {
			return apiHelper.ResponseWrapper(c, http.StatusBadRequest, errors.New("InvalidPageNumber"))
		}
		paginationDto.PageNumber = pageNumber
	}

	if requestBody["itemsPerPage"] != nil {
		itemsPerPage, err := voHelper.InterfaceToUint16(requestBody["itemsPerPage"])
		if err != nil {
			return apiHelper.ResponseWrapper(c, http.StatusBadRequest, errors.New("InvalidItemsPerPage"))
		}
		paginationDto.ItemsPerPage = itemsPerPage
	}

	if requestBody["sortBy"] != nil {
		sortBy, err := valueObject.NewPaginationSortBy(requestBody["sortBy"])
		if err != nil {
			return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err)
		}
		paginationDto.SortBy = &sortBy
	}

	if requestBody["sortDirection"] != nil {
		sortDirection, err := valueObject.NewPaginationSortDirection(requestBody["sortDirection"])
		if err != nil {
			return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err)
		}
		paginationDto.SortDirection = &sortDirection
	}

	if requestBody["lastSeenId"] != nil {
		lastSeenId, err := valueObject.NewPaginationLastSeenId(requestBody["lastSeenId"])
		if err != nil {
			return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err)
		}
		paginationDto.LastSeenId = &lastSeenId
	}

	readDto := dto.ReadMarketplaceItemsRequest{
		Pagination: paginationDto,
		ItemSlug:   itemSlugPtr,
		ItemName:   itemNamePtr,
		ItemType:   itemTypePtr,
	}

	marketplaceQueryRepo := infra.NewMarketplaceQueryRepo()

	marketplaceItemsList, err := useCase.ReadMarketplaceItems(marketplaceQueryRepo, readDto)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err)
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, marketplaceItemsList)
}

func (controller *MarketplaceController) RefreshMarketplace() {
	taskInterval := time.Duration(24) * time.Hour
	timer := time.NewTicker(taskInterval)
	defer timer.Stop()

	marketplaceCmdRepo := infra.NewMarketplaceCmdRepo()
	for range timer.C {
		useCase.RefreshMarketplace(marketplaceCmdRepo)
	}
}
