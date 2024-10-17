package marketplaceInfra

import (
	"testing"

	testHelpers "github.com/goinfinite/ez/src/devUtils"
	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/useCase"
	"github.com/goinfinite/ez/src/domain/valueObject"
)

func TestMarketplaceQueryRepo(t *testing.T) {
	testHelpers.LoadEnvVars()

	t.Run("Read", func(t *testing.T) {
		itemType, _ := valueObject.NewMarketplaceItemType("app")

		paginationDto := useCase.MarketplaceDefaultPagination
		sortBy, _ := valueObject.NewPaginationSortBy("id")
		sortDirection, _ := valueObject.NewPaginationSortDirection("desc")
		paginationDto.SortBy = &sortBy
		paginationDto.SortDirection = &sortDirection

		readDto := dto.ReadMarketplaceItemsRequest{
			Pagination: paginationDto,
			ItemType:   &itemType,
		}

		marketplaceQueryRepo := NewMarketplaceQueryRepo()
		responseDto, err := marketplaceQueryRepo.Read(readDto)
		if err != nil {
			t.Errorf("ReadMarketplaceItemsError: %v", err)
			return
		}

		if len(responseDto.Items) == 0 {
			t.Errorf("NoItemsFound")
		}
	})
}
