package useCase

import (
	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/repository"
)

func ReadContainers(
	containerQueryRepo repository.ContainerQueryRepo,
) ([]entity.Container, error) {
	return containerQueryRepo.Read()
}
