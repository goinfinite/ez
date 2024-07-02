package useCase

import (
	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/repository"
)

func ReadContainersWithMetrics(
	containerQueryRepo repository.ContainerQueryRepo,
) ([]dto.ContainerWithMetrics, error) {
	return containerQueryRepo.ReadWithMetrics()
}
