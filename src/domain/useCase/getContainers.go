package useCase

import (
	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/repository"
)

func GetContainers(
	containerQueryRepo repository.ContainerQueryRepo,
) ([]entity.Container, error) {
	return containerQueryRepo.Get()
}
