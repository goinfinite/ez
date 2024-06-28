package useCase

import (
	"errors"
	"log"
	"time"

	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/repository"
	"github.com/speedianet/control/src/domain/valueObject"
)

func GetSessionToken(
	authQueryRepo repository.AuthQueryRepo,
	authCmdRepo repository.AuthCmdRepo,
	accountQueryRepo repository.AccountQueryRepo,
	login dto.Login,
	ipAddress valueObject.IpAddress,
) (accessToken entity.AccessToken, err error) {
	isLoginValid := authQueryRepo.IsLoginValid(login)

	if !isLoginValid {
		log.Printf(
			"Login failed for '%v' from '%v'.",
			login.Username.String(),
			ipAddress.String(),
		)
		return accessToken, errors.New("InvalidCredentials")
	}

	accountDetails, err := accountQueryRepo.GetByUsername(login.Username)
	if err != nil {
		return accessToken, errors.New("AccountNotFound")
	}

	accountId := accountDetails.Id
	expiresIn := valueObject.NewUnixTimeAfterNow(3 * time.Hour)

	return authCmdRepo.GenerateSessionToken(accountId, expiresIn, ipAddress)
}
