package service

import (
	"errors"

	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/useCase"
	"github.com/goinfinite/ez/src/domain/valueObject"
	voHelper "github.com/goinfinite/ez/src/domain/valueObject/helper"
	marketplaceInfra "github.com/goinfinite/ez/src/infra/marketplace"
)

type MarketplaceService struct {
}

func NewMarketplaceService() *MarketplaceService {
	return &MarketplaceService{}
}

func (service *MarketplaceService) Read(input map[string]interface{}) ServiceOutput {
	var itemSlugPtr *valueObject.MarketplaceItemSlug
	if input["itemSlug"] != nil {
		itemSlug, err := valueObject.NewMarketplaceItemSlug(input["itemSlug"])
		if err != nil {
			return NewServiceOutput(UserError, err)
		}
		itemSlugPtr = &itemSlug
	}

	var itemNamePtr *valueObject.MarketplaceItemName
	if input["itemName"] != nil {
		itemName, err := valueObject.NewMarketplaceItemName(input["itemName"])
		if err != nil {
			return NewServiceOutput(UserError, err)
		}
		itemNamePtr = &itemName
	}

	var itemTypePtr *valueObject.MarketplaceItemType
	if input["itemType"] != nil {
		itemType, err := valueObject.NewMarketplaceItemType(input["itemType"])
		if err != nil {
			return NewServiceOutput(UserError, err)
		}
		itemTypePtr = &itemType
	}

	paginationDto := useCase.MarketplaceDefaultPagination
	if input["pageNumber"] != nil {
		pageNumber, err := voHelper.InterfaceToUint32(input["pageNumber"])
		if err != nil {
			return NewServiceOutput(UserError, errors.New("InvalidPageNumber"))
		}
		paginationDto.PageNumber = pageNumber
	}

	if input["itemsPerPage"] != nil {
		itemsPerPage, err := voHelper.InterfaceToUint16(input["itemsPerPage"])
		if err != nil {
			return NewServiceOutput(UserError, errors.New("InvalidItemsPerPage"))
		}
		paginationDto.ItemsPerPage = itemsPerPage
	}

	if input["sortBy"] != nil {
		sortBy, err := valueObject.NewPaginationSortBy(input["sortBy"])
		if err != nil {
			return NewServiceOutput(UserError, err)
		}
		paginationDto.SortBy = &sortBy
	}

	if input["sortDirection"] != nil {
		sortDirection, err := valueObject.NewPaginationSortDirection(input["sortDirection"])
		if err != nil {
			return NewServiceOutput(UserError, err)
		}
		paginationDto.SortDirection = &sortDirection
	}

	if input["lastSeenId"] != nil {
		lastSeenId, err := valueObject.NewPaginationLastSeenId(input["lastSeenId"])
		if err != nil {
			return NewServiceOutput(UserError, err)
		}
		paginationDto.LastSeenId = &lastSeenId
	}

	readDto := dto.ReadMarketplaceItemsRequest{
		Pagination: paginationDto,
		ItemSlug:   itemSlugPtr,
		ItemName:   itemNamePtr,
		ItemType:   itemTypePtr,
	}

	marketplaceQueryRepo := marketplaceInfra.NewMarketplaceQueryRepo()

	responseDto, err := useCase.ReadMarketplaceItems(marketplaceQueryRepo, readDto)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, responseDto)
}
