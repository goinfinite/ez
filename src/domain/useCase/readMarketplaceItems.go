package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/repository"
)

var MarketplaceDefaultPagination dto.Pagination = dto.Pagination{
	PageNumber:   0,
	ItemsPerPage: 5,
}

func ReadMarketplaceItems(
	marketplaceQueryRepo repository.MarketplaceQueryRepo,
	readDto dto.ReadMarketplaceItemsRequest,
) (responseDto dto.ReadMarketplaceItemsResponse, err error) {
	responseDto, err = marketplaceQueryRepo.Read(readDto)
	if err != nil {
		slog.Error("ReadMarketplaceItemsInfraError", slog.Any("error", err))
		return responseDto, errors.New("ReadMarketplaceItemsInfraError")
	}

	return responseDto, nil
}
