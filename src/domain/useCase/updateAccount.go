package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/repository"
	"github.com/speedianet/control/src/domain/valueObject"
)

func UpdateAccount(
	accountQueryRepo repository.AccountQueryRepo,
	accountCmdRepo repository.AccountCmdRepo,
	securityCmdRepo repository.SecurityCmdRepo,
	updateDto dto.UpdateAccount,
) error {
	_, err := accountQueryRepo.ReadById(updateDto.AccountId)
	if err != nil {
		return errors.New("AccountNotFound")
	}

	if updateDto.Password != nil {
		err = accountCmdRepo.UpdatePassword(updateDto.AccountId, *updateDto.Password)
		if err != nil {
			log.Printf("UpdateAccountPasswordError: %s", err)
			return errors.New("UpdateAccountPasswordInfraError")
		}

		eventType, _ := valueObject.NewSecurityEventType("account-password-updated")
		createSecurityEventDto := dto.NewCreateSecurityEvent(
			eventType, nil, &updateDto.IpAddress, &updateDto.AccountId,
		)
		AsyncCreateSecurityEvent(securityCmdRepo, createSecurityEventDto)
	}

	if updateDto.Quota == nil {
		return nil
	}

	err = accountCmdRepo.UpdateQuota(updateDto.AccountId, *updateDto.Quota)
	if err != nil {
		log.Printf("UpdateAccountQuotaError: %s", err)
		return errors.New("UpdateAccountQuotaInfraError")
	}

	eventType, _ := valueObject.NewSecurityEventType("account-quota-updated")
	createSecurityEventDto := dto.NewCreateSecurityEvent(
		eventType, nil, &updateDto.IpAddress, &updateDto.AccountId,
	)
	AsyncCreateSecurityEvent(securityCmdRepo, createSecurityEventDto)

	return nil
}
