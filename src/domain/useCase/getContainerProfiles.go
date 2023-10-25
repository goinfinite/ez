package useCase

import (
	"github.com/goinfinite/fleet/src/domain/entity"
	"github.com/goinfinite/fleet/src/domain/repository"
)

func GetContainerProfiles(
	containerProfileQueryRepo repository.ContainerProfileQueryRepo,
) ([]entity.ContainerProfile, error) {
	return containerProfileQueryRepo.Get()
}
