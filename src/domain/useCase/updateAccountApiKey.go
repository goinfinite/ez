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
	securityCmdRepo repository.SecurityCmdRepo,
	updateDto dto.UpdateAccount,
) (newKey valueObject.AccessTokenValue, err error) {
	_, err = accountQueryRepo.ReadById(updateDto.AccountId)
	if err != nil {
		return newKey, errors.New("AccountNotFound")
	}

	newKey, err = accountCmdRepo.UpdateApiKey(updateDto.AccountId)
	if err != nil {
		log.Printf("UpdateAccountApiKeyError: %s", err)
		return newKey, errors.New("UpdateAccountApiKeyInfraError")
	}

	eventType, _ := valueObject.NewSecurityEventType("account-api-key-updated")
	createSecurityEventDto := dto.NewCreateSecurityEvent(
		eventType, nil, &updateDto.IpAddress, &updateDto.AccountId,
	)
	AsyncCreateSecurityEvent(securityCmdRepo, createSecurityEventDto)

	return newKey, nil
}
