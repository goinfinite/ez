package useCase

import (
	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/repository"
)

func GetContainersWithUsage(
	containerQueryRepo repository.ContainerQueryRepo,
) ([]dto.ContainerWithUsage, error) {
	return containerQueryRepo.GetWithUsage()
}
