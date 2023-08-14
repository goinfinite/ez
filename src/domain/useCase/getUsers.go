package useCase

import (
	"github.com/speedianet/sfm/src/domain/entity"
	"github.com/speedianet/sfm/src/domain/repository"
)

func GetUsers(
	accQueryRepo repository.AccQueryRepo,
) ([]entity.AccountDetails, error) {
	return accQueryRepo.Get()
}
