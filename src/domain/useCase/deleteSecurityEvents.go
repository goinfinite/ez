package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/repository"
)

func DeleteSecurityEvents(
	securityCmdRepo repository.SecurityCmdRepo,
	deleteDto dto.DeleteSecurityEvents,
) error {
	err := securityCmdRepo.DeleteEvents(deleteDto)
	if err != nil {
		log.Printf("DeleteSecurityEventsError: %v", err)
		return errors.New("DeleteSecurityEventsInfraError")
	}

	return nil
}
