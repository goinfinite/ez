package infra

import (
	"github.com/goinfinite/ez/src/domain/dto"
)

type MarketplaceQueryRepo struct {
}

func NewMarketplaceQueryRepo() *MarketplaceQueryRepo {
	return &MarketplaceQueryRepo{}
}

func (repo *MarketplaceQueryRepo) Read(
	readDto dto.ReadMarketplaceItemsRequest,
) (responseDto dto.ReadMarketplaceItemsResponse, err error) {
	return responseDto, nil
}
