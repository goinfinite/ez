package useCase

import (
	"log"

	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/repository"
)

func ReadSecurityEvents(
	securityQueryRepo repository.SecurityQueryRepo,
	readDto dto.ReadSecurityEvents,
) (securityEvents []entity.SecurityEvent) {
	securityEvents, err := securityQueryRepo.ReadEvents(readDto)
	if err != nil {
		log.Printf("ReadSecurityEventsInfraError: %s", err)
	}

	return securityEvents
}
