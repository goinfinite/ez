package useCase

import (
	"github.com/goinfinite/fleet/src/domain/entity"
	"github.com/goinfinite/fleet/src/domain/repository"
)

func GetO11yOverview(
	o11yQueryRepo repository.O11yQueryRepo,
) (entity.O11yOverview, error) {
	return o11yQueryRepo.GetOverview()
}
