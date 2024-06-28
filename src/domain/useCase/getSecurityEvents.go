package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/repository"
)

const MaxFailedLoginAttemptsPerIpAddress uint = 5

func GetSecurityEvents(
	securityQueryRepo repository.SecurityQueryRepo,
	getDto dto.GetSecurityEvents,
) ([]entity.SecurityEvent, error) {
	securityEvents, err := securityQueryRepo.GetEvents(getDto)
	if err != nil {
		log.Printf("GetSecurityEventsError: %s", err)
		return nil, errors.New("GetSecurityEventsInfraError")
	}

	return securityEvents, nil
}
