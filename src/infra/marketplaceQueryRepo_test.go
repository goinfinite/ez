package infra

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

		readDto := dto.ReadMarketplaceItemsRequest{
			Pagination: useCase.MarketplaceDefaultPagination,
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
