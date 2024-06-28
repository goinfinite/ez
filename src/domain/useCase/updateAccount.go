package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/repository"
)

func UpdateAccount(
	accountQueryRepo repository.AccountQueryRepo,
	accountCmdRepo repository.AccountCmdRepo,
	updateAccountDto dto.UpdateAccount,
) error {
	_, err := accountQueryRepo.GetById(updateAccountDto.AccountId)
	if err != nil {
		return errors.New("AccountNotFound")
	}

	if updateAccountDto.Password != nil {
		err = accountCmdRepo.UpdatePassword(
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
		err = accountCmdRepo.UpdateQuota(
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
