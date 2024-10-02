package repository

import (
	"github.com/goinfinite/ez/src/domain/dto"
)

type MarketplaceQueryRepo interface {
	Read(dto.ReadMarketplaceItemsRequest) (dto.ReadMarketplaceItemsResponse, error)
}
