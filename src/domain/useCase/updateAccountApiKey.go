package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/repository"
	"github.com/speedianet/control/src/domain/valueObject"
)

func UpdateAccountApiKey(
	accountQueryRepo repository.AccountQueryRepo,
	accountCmdRepo repository.AccountCmdRepo,
	updateAccountDto dto.UpdateAccount,
) (valueObject.AccessTokenValue, error) {
	_, err := accountQueryRepo.GetById(updateAccountDto.AccountId)
	if err != nil {
		return "", errors.New("AccountNotFound")
	}

	newKey, err := accountCmdRepo.UpdateApiKey(updateAccountDto.AccountId)
	if err != nil {
		log.Printf("UpdateAccountApiKeyError: %s", err)
		return "", errors.New("UpdateAccountApiKeyInfraError")
	}

	log.Printf("AccountId '%v' api key updated.", updateAccountDto.AccountId)

	return newKey, nil
}
