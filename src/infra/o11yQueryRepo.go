package infra

import (
	"github.com/goinfinite/fleet/src/domain/entity"
	o11yInfra "github.com/goinfinite/fleet/src/infra/o11y"
)

type O11yQueryRepo struct {
}

func (repo O11yQueryRepo) GetOverview() (entity.O11yOverview, error) {
	getOverviewRepo := o11yInfra.GetOverview{}
	return getOverviewRepo.Get()
}
