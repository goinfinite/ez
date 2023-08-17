package useCase

import (
	"github.com/speedianet/sfm/src/domain/entity"
	"github.com/speedianet/sfm/src/domain/repository"
)

func GetContainers(
	containerQueryRepo repository.ContainerQueryRepo,
) ([]entity.Container, error) {
	return containerQueryRepo.Get()
}
