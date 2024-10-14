package marketplaceInfra

import (
	"testing"

	testHelpers "github.com/goinfinite/ez/src/devUtils"
)

func TestMarketplaceCmdRepo(t *testing.T) {
	testHelpers.LoadEnvVars()

	t.Run("Refresh", func(t *testing.T) {
		marketplaceCmdRepo := NewMarketplaceCmdRepo()
		err := marketplaceCmdRepo.Refresh()
		if err != nil {
			t.Errorf("RefreshMarketplaceError: %v", err)
		}
	})
}
