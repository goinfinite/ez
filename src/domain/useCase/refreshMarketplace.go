package useCase

import (
	"log/slog"

	"github.com/goinfinite/ez/src/domain/repository"
)

func RefreshMarketplace(
	marketplaceCmdRepo repository.MarketplaceCmdRepo,
) {
	err := marketplaceCmdRepo.Refresh()
	if err != nil {
		slog.Error("RefreshMarketplaceInfraError", slog.Any("error", err))
	}
}
