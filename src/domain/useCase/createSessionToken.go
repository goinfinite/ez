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

func getFailedLoginAttemptsCount(
	activityRecordQueryRepo repository.ActivityRecordQueryRepo,
	loginDto dto.Login,
) uint {
	secLevel, _ := valueObject.NewActivityRecordLevel("SEC")
	recordCode, _ := valueObject.NewActivityRecordCode("LoginFailed")
	failedAttemptsIntervalStartsAt := valueObject.NewUnixTimeBeforeNow(
		FailedLoginAttemptsInterval,
	)
	readActivityRecordsDto := dto.NewReadActivityRecords(
		&secLevel, &recordCode, nil, loginDto.IpAddress, nil, nil, nil,
		nil, nil, nil, &failedAttemptsIntervalStartsAt,
	)

	failedLoginAttempts := ReadActivityRecords(
		activityRecordQueryRepo, readActivityRecordsDto,
	)

	return uint(len(failedLoginAttempts))
}

func CreateSessionToken(
	authQueryRepo repository.AuthQueryRepo,
	authCmdRepo repository.AuthCmdRepo,
	accountQueryRepo repository.AccountQueryRepo,
	activityRecordQueryRepo repository.ActivityRecordQueryRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	loginDto dto.Login,
) (accessToken entity.AccessToken, err error) {
	failedAttemptsCount := getFailedLoginAttemptsCount(activityRecordQueryRepo, loginDto)
	if failedAttemptsCount >= MaxFailedLoginAttemptsPerIpAddress {
		return accessToken, errors.New("MaxFailedLoginAttemptsReached")
	}

	if !authQueryRepo.IsLoginValid(loginDto) {
		recordCode, _ := valueObject.NewActivityRecordCode("LoginFailed")
		NewCreateSecurityActivityRecord(activityRecordCmdRepo).CreateSessionToken(loginDto, recordCode)

		return accessToken, errors.New("InvalidCredentials")
	}

	accountEntity, err := accountQueryRepo.ReadByUsername(loginDto.Username)
	if err != nil {
		return accessToken, errors.New("AccountNotFound")
	}

	recordCode, _ := valueObject.NewActivityRecordCode("LoginSuccessful")
	NewCreateSecurityActivityRecord(activityRecordCmdRepo).CreateSessionToken(loginDto, recordCode)

	expiresIn := valueObject.NewUnixTimeAfterNow(SessionTokenExpiresIn)

	return authCmdRepo.CreateSessionToken(
		accountEntity.Id, expiresIn, *loginDto.IpAddress,
	)
}
