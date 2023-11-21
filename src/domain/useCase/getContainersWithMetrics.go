package useCase

import (
	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/repository"
)

func GetContainersWithMetrics(
	containerQueryRepo repository.ContainerQueryRepo,
) ([]dto.ContainerWithMetrics, error) {
	return containerQueryRepo.GetWithMetrics()
}
