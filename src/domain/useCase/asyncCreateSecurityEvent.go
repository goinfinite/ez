package useCase

import (
	"log"

	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/repository"
)

func AsyncCreateSecurityEvent(
	securityCmdRepo repository.SecurityCmdRepo,
	createDto dto.CreateSecurityEvent,
) {
	wasIpAddressSet := createDto.IpAddress != nil
	wasAccountIdSet := createDto.AccountId != nil
	if !wasIpAddressSet && !wasAccountIdSet {
		log.Printf("SecurityEventWithoutIpAddressOrAccountIdNotCreated: %v", createDto)
		return
	}

	go func() {
		err := securityCmdRepo.CreateEvent(createDto)
		if err != nil {
			log.Printf("CreateSecurityEventInfraError: %v", err)
		}
	}()
}
