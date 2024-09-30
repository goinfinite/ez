package repository

import (
	"github.com/goinfinite/ez/src/domain/entity"
	"github.com/goinfinite/ez/src/domain/valueObject"
)

type AuthCmdRepo interface {
	CreateSessionToken(
		accountId valueObject.AccountId,
		expiresIn valueObject.UnixTime,
		ipAddress valueObject.IpAddress,
	) (entity.AccessToken, error)
}
