package useCase

import (
	"errors"
	"log/slog"

	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/repository"
)

func ReadO11yOverview(
	o11yQueryRepo repository.O11yQueryRepo,
) (entity.O11yOverview, error) {
	overview, err := o11yQueryRepo.ReadOverview()
	if err != nil {
		slog.Error("ReadO11yOverviewInfraError", slog.Any("error", err))
		return overview, errors.New("ReadO11yOverviewInfraError")
	}

	return overview, nil
}
