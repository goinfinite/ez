package useCase

import (
	"github.com/goinfinite/fleet/src/domain/entity"
	"github.com/goinfinite/fleet/src/domain/repository"
)

func GetContainers(
	containerQueryRepo repository.ContainerQueryRepo,
) ([]entity.Container, error) {
	return containerQueryRepo.Get()
}
