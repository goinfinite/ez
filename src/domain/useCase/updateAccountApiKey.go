package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/repository"
	"github.com/speedianet/control/src/domain/valueObject"
)

func UpdateAccountApiKey(
	accQueryRepo repository.AccQueryRepo,
	accCmdRepo repository.AccCmdRepo,
	updateAccountDto dto.UpdateAccount,
) (valueObject.AccessTokenValue, error) {
	_, err := accQueryRepo.GetById(updateAccountDto.AccountId)
	if err != nil {
		return "", errors.New("AccountNotFound")
	}

	newKey, err := accCmdRepo.UpdateApiKey(updateAccountDto.AccountId)
	if err != nil {
		log.Printf("UpdateAccountApiKeyError: %s", err)
		return "", errors.New("UpdateAccountApiKeyInfraError")
	}

	log.Printf("AccountId '%v' api key updated.", updateAccountDto.AccountId)

	return newKey, nil
}
