package useCase

import (
	"errors"
	"time"

	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/repository"
	"github.com/speedianet/control/src/domain/valueObject"
)

const MaxFailedLoginAttemptsPerIpAddress uint = 3
const FailedLoginAttemptsInterval time.Duration = 15 * time.Minute
const SessionTokenExpiresIn time.Duration = 3 * time.Hour

func GenerateSessionToken(
	authQueryRepo repository.AuthQueryRepo,
	authCmdRepo repository.AuthCmdRepo,
	accountQueryRepo repository.AccountQueryRepo,
	securityQueryRepo repository.SecurityQueryRepo,
	securityCmdRepo repository.SecurityCmdRepo,
	loginDto dto.Login,
) (accessToken entity.AccessToken, err error) {
	eventType, _ := valueObject.NewSecurityEventType("failed-login")
	failedAttemptsIntervalStartsAt := valueObject.NewUnixTimeBeforeNow(
		FailedLoginAttemptsInterval,
	)
	readSecurityEventsDto := dto.NewReadSecurityEvents(
		&eventType, loginDto.IpAddress, nil, &failedAttemptsIntervalStartsAt,
	)

	failedLoginAttempts, err := ReadSecurityEvents(securityQueryRepo, readSecurityEventsDto)
	if err != nil {
		return accessToken, err
	}
	failedAttemptsCount := uint(len(failedLoginAttempts))
	if failedAttemptsCount >= MaxFailedLoginAttemptsPerIpAddress {
		return accessToken, errors.New("MaxFailedLoginAttemptsReached")
	}

	if !authQueryRepo.IsLoginValid(loginDto) {
		eventDetails, _ := valueObject.NewSecurityEventDetails(
			"Username: " + loginDto.Username.String(),
		)
		createSecurityEventDto := dto.NewCreateSecurityEvent(
			eventType, &eventDetails, loginDto.IpAddress, nil,
		)
		err = CreateSecurityEvent(securityCmdRepo, createSecurityEventDto)
		if err != nil {
			return accessToken, err
		}

		return accessToken, errors.New("InvalidCredentials")
	}

	accountDetails, err := accountQueryRepo.ReadByUsername(loginDto.Username)
	if err != nil {
		return accessToken, errors.New("AccountNotFound")
	}

	eventType, _ = valueObject.NewSecurityEventType("successful-login")
	createSecurityEventDto := dto.NewCreateSecurityEvent(
		eventType, nil, loginDto.IpAddress, &accountDetails.Id,
	)
	err = CreateSecurityEvent(securityCmdRepo, createSecurityEventDto)
	if err != nil {
		return accessToken, err
	}

	expiresIn := valueObject.NewUnixTimeAfterNow(SessionTokenExpiresIn)

	return authCmdRepo.GenerateSessionToken(
		accountDetails.Id, expiresIn, *loginDto.IpAddress,
	)
}
