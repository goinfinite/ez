package marketplaceInfra

import (
	"errors"
	"os"

	infraEnvs "github.com/goinfinite/ez/src/infra/envs"
	infraHelper "github.com/goinfinite/ez/src/infra/helper"
)

type MarketplaceCmdRepo struct {
}

func NewMarketplaceCmdRepo() *MarketplaceCmdRepo {
	return &MarketplaceCmdRepo{}
}

func (repo *MarketplaceCmdRepo) Refresh() error {
	_, err := os.Stat(infraEnvs.MarketplaceDir)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}

		_, err = infraHelper.RunCmdWithSubShell(
			"cd " + infraEnvs.InfiniteEzMainDir + ";" +
				"git clone https://github.com/goinfinite/ez-marketplace.git marketplace",
		)
		if err != nil {
			return errors.New("CloneMarketplaceRepoError: " + err.Error())
		}
	}

	_, err = infraHelper.RunCmdWithSubShell(
		"cd " + infraEnvs.MarketplaceDir + ";" +
			"git clean -f -d; git reset --hard HEAD; git pull",
	)
	return err
}
