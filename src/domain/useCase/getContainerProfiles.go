package useCase

import (
	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/repository"
)

func GetContainerProfiles(
	containerProfileQueryRepo repository.ContainerProfileQueryRepo,
) ([]entity.ContainerProfile, error) {
	return containerProfileQueryRepo.Get()
}
