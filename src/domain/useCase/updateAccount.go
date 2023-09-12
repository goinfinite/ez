package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/sfm/src/domain/dto"
	"github.com/speedianet/sfm/src/domain/repository"
)

func UpdateAccount(
	accQueryRepo repository.AccQueryRepo,
	accCmdRepo repository.AccCmdRepo,
	updateAccountDto dto.UpdateAccount,
) error {
	_, err := accQueryRepo.GetById(updateAccountDto.AccountId)
	if err != nil {
		return errors.New("AccountNotFound")
	}

	if updateAccountDto.Password != nil {
		err = accCmdRepo.UpdatePassword(
			updateAccountDto.AccountId,
			*updateAccountDto.Password,
		)
		if err != nil {
			log.Printf("UpdateAccountPasswordError: %s", err)
			return errors.New("UpdateAccountPasswordInfraError")
		}

		log.Printf("AccountId '%v' password updated.", updateAccountDto.AccountId)
	}

	if updateAccountDto.Quota != nil {
		err = accCmdRepo.UpdateQuota(
			updateAccountDto.AccountId,
			*updateAccountDto.Quota,
		)
		if err != nil {
			log.Printf("UpdateAccountQuotaError: %s", err)
			return errors.New("UpdateAccountQuotaInfraError")
		}

		log.Printf("AccountId '%v' quota updated.", updateAccountDto.AccountId)
	}

	return nil
}
