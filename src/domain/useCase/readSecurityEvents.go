package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/repository"
)

func ReadSecurityEvents(
	securityQueryRepo repository.SecurityQueryRepo,
	readDto dto.ReadSecurityEvents,
) ([]entity.SecurityEvent, error) {
	securityEvents, err := securityQueryRepo.ReadEvents(readDto)
	if err != nil {
		log.Printf("ReadSecurityEventsError: %s", err)
		return nil, errors.New("ReadSecurityEventsInfraError")
	}

	return securityEvents, nil
}
