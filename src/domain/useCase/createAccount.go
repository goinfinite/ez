package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/repository"
	"github.com/speedianet/control/src/domain/valueObject"
)

func CreateAccount(
	accountQueryRepo repository.AccountQueryRepo,
	accountCmdRepo repository.AccountCmdRepo,
	securityCmdRepo repository.SecurityCmdRepo,
	createDto dto.CreateAccount,
) error {
	_, err := accountQueryRepo.ReadByUsername(createDto.Username)
	if err == nil {
		return errors.New("AccountAlreadyExists")
	}

	defaultQuota := valueObject.NewAccountQuotaWithDefaultValues()
	if createDto.Quota == nil {
		createDto.Quota = &defaultQuota
	}

	err = accountCmdRepo.Create(createDto)
	if err != nil {
		log.Printf("CreateAccountError: %s", err)
		return errors.New("CreateAccountInfraError")
	}

	eventType, _ := valueObject.NewSecurityEventType("account-created")
	eventDetails, _ := valueObject.NewSecurityEventDetails(
		"Username: " + createDto.Username.String(),
	)
	createSecurityEventDto := dto.NewCreateSecurityEvent(
		eventType, &eventDetails, &createDto.IpAddress, nil,
	)
	CreateSecurityEvent(securityCmdRepo, createSecurityEventDto)

	return nil
}
