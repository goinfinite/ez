package repository

import (
	"github.com/goinfinite/fleet/src/domain/entity"
	"github.com/goinfinite/fleet/src/domain/valueObject"
)

type AuthCmdRepo interface {
	GenerateSessionToken(
		accountId valueObject.AccountId,
		expiresIn valueObject.UnixTime,
		ipAddress valueObject.IpAddress,
	) (entity.AccessToken, error)
}
