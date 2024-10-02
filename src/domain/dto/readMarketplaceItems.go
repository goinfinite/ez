package dto

import (
	"github.com/goinfinite/ez/src/domain/entity"
	"github.com/goinfinite/ez/src/domain/valueObject"
)

type ReadMarketplaceItemsRequest struct {
	Pagination Pagination                       `json:"pagination"`
	ItemSlug   *valueObject.MarketplaceItemSlug `json:"itemSlug,omitempty"`
	ItemName   *valueObject.MarketplaceItemName `json:"itemName,omitempty"`
	ItemType   *valueObject.MarketplaceItemType `json:"itemType,omitempty"`
}

type ReadMarketplaceItemsResponse struct {
	Pagination Pagination               `json:"pagination"`
	Items      []entity.MarketplaceItem `json:"items"`
}
