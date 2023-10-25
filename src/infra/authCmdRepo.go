package infra

import (
	"errors"
	"os"
	"time"

	"github.com/goinfinite/fleet/src/domain/entity"
	"github.com/goinfinite/fleet/src/domain/valueObject"
	"github.com/golang-jwt/jwt"
)

type AuthCmdRepo struct {
}

func (repo AuthCmdRepo) GenerateSessionToken(
	accountId valueObject.AccountId,
	expiresIn valueObject.UnixTime,
	ipAddress valueObject.IpAddress,
) (entity.AccessToken, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	apiUrl, err := os.Hostname()
	if err != nil {
		apiUrl = "localhost"
	}

	now := time.Now()
	tokenExpiration := time.Unix(expiresIn.Get(), 0)

	claims := jwt.MapClaims{
		"iss":        apiUrl,
		"iat":        now.Unix(),
		"nbf":        now.Unix(),
		"exp":        tokenExpiration.Unix(),
		"accountId":  accountId.Get(),
		"originalIp": ipAddress.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStrUnparsed, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return entity.AccessToken{}, errors.New("SessionTokenGenerationError")
	}

	tokenType := valueObject.NewAccessTokenTypePanic("sessionToken")
	tokenStr := valueObject.NewAccessTokenStrPanic(tokenStrUnparsed)

	return entity.NewAccessToken(
		tokenType,
		expiresIn,
		tokenStr,
	), nil
}
