package infra

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/speedianet/sfm/src/domain/entity"
	"github.com/speedianet/sfm/src/domain/valueObject"
)

type AuthCmdRepo struct {
}

func (repo AuthCmdRepo) GenerateSessionToken(
	userId valueObject.UserId,
	expiresIn valueObject.UnixTime,
	ipAddress valueObject.IpAddress,
) entity.AccessToken {
	jwtSecret := os.Getenv("JWT_SECRET")
	apiURL := os.Getenv("VIRTUAL_HOST")

	now := time.Now()
	tokenExpiration := time.Unix(expiresIn.Get(), 0)

	claims := jwt.MapClaims{
		"iss":        apiURL,
		"iat":        now.Unix(),
		"nbf":        now.Unix(),
		"exp":        tokenExpiration.Unix(),
		"userId":     userId.Get(),
		"originalIp": ipAddress.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStrUnparsed, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		panic("SessionTokenGenerationError")
	}

	tokenType := valueObject.NewAccessTokenTypePanic("sessionToken")
	tokenStr := valueObject.NewAccessTokenStrPanic(tokenStrUnparsed)

	return entity.NewAccessToken(
		tokenType,
		expiresIn,
		tokenStr,
	)
}
