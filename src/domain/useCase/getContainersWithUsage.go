package useCase

import (
	"github.com/goinfinite/fleet/src/domain/dto"
	"github.com/goinfinite/fleet/src/domain/repository"
)

func GetContainersWithUsage(
	containerQueryRepo repository.ContainerQueryRepo,
) ([]dto.ContainerWithUsage, error) {
	return containerQueryRepo.GetWithUsage()
}
