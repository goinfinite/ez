package useCase

import (
	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/repository"
)

func ReadO11yOverview(
	o11yQueryRepo repository.O11yQueryRepo,
) (entity.O11yOverview, error) {
	return o11yQueryRepo.ReadOverview()
}