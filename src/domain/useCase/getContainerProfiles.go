package useCase

import (
	"github.com/speedianet/sfm/src/domain/entity"
	"github.com/speedianet/sfm/src/domain/repository"
)

func GetContainerProfiles(
	containerProfileQueryRepo repository.ContainerProfileQueryRepo,
) ([]entity.ContainerProfile, error) {
	return containerProfileQueryRepo.Get()
}
