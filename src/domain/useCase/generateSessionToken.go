package useCase

import (
	"errors"
	"time"

	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/repository"
	"github.com/speedianet/control/src/domain/valueObject"
)

const MaxFailedLoginAttemptsPerIpAddress uint = 5

func GenerateSessionToken(
	authQueryRepo repository.AuthQueryRepo,
	authCmdRepo repository.AuthCmdRepo,
	accountQueryRepo repository.AccountQueryRepo,
	loginDto dto.Login,
) (accessToken entity.AccessToken, err error) {
	if !authQueryRepo.IsLoginValid(loginDto) {
		return accessToken, errors.New("InvalidCredentials")
	}

	accountDetails, err := accountQueryRepo.GetByUsername(loginDto.Username)
	if err != nil {
		return accessToken, errors.New("AccountNotFound")
	}

	accountId := accountDetails.Id
	expiresIn := valueObject.NewUnixTimeAfterNow(3 * time.Hour)

	return authCmdRepo.GenerateSessionToken(accountId, expiresIn, *loginDto.IpAddress)
}
