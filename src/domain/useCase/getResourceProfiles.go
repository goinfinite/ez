package useCase

import (
	"github.com/speedianet/sfm/src/domain/entity"
	"github.com/speedianet/sfm/src/domain/repository"
)

func GetResourceProfiles(
	resourceProfileQueryRepo repository.ResourceProfileQueryRepo,
) ([]entity.ResourceProfile, error) {
	return resourceProfileQueryRepo.Get()
}
