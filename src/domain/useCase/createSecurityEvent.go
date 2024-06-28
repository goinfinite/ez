package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/repository"
)

func CreateSecurityEvent(
	securityCmdRepo repository.SecurityCmdRepo,
	createDto dto.CreateSecurityEvent,
) error {
	wasIpAddressSet := createDto.IpAddress != nil
	wasAccountIdSet := createDto.AccountId != nil
	if !wasIpAddressSet && !wasAccountIdSet {
		return errors.New("SecurityEventRequiresIpAddressOrAccountId")
	}

	err := securityCmdRepo.CreateEvent(createDto)
	if err != nil {
		log.Printf("CreateSecurityEventError: %v", err)
		return errors.New("CreateSecurityEventInfraError")
	}

	return nil
}
