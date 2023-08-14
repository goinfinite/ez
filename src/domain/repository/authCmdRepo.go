package repository

import (
	"github.com/speedianet/sfm/src/domain/entity"
	"github.com/speedianet/sfm/src/domain/valueObject"
)

type AuthCmdRepo interface {
	GenerateSessionToken(
		userId valueObject.UserId,
		expiresIn valueObject.UnixTime,
		ipAddress valueObject.IpAddress,
	) entity.AccessToken
}
