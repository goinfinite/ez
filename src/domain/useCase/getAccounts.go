package useCase

import (
	"github.com/speedianet/sfm/src/domain/entity"
	"github.com/speedianet/sfm/src/domain/repository"
)

func GetAccounts(
	accQueryRepo repository.AccQueryRepo,
) ([]entity.Account, error) {
	return accQueryRepo.Get()
}
