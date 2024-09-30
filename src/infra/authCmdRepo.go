package infra

import (
	"errors"
	"os"
	"time"

	"github.com/goinfinite/ez/src/domain/entity"
	"github.com/goinfinite/ez/src/domain/valueObject"
	"github.com/golang-jwt/jwt"
)

type AuthCmdRepo struct {
}

func (repo AuthCmdRepo) CreateSessionToken(
	accountId valueObject.AccountId,
	expiresIn valueObject.UnixTime,
	ipAddress valueObject.IpAddress,
) (accessToken entity.AccessToken, err error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	apiUrl, err := os.Hostname()
	if err != nil {
		apiUrl = "localhost"
	}

	now := time.Now()
	tokenExpiration := time.Unix(expiresIn.Read(), 0)

	claims := jwt.MapClaims{
		"iss":        apiUrl,
		"iat":        now.Unix(),
		"nbf":        now.Unix(),
		"exp":        tokenExpiration.Unix(),
		"accountId":  accountId.Uint64(),
		"originalIp": ipAddress.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStrUnparsed, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return accessToken, errors.New("SessionTokenGenerationError")
	}

	tokenType, _ := valueObject.NewAccessTokenType("sessionToken")
	tokenStr, err := valueObject.NewAccessTokenValue(tokenStrUnparsed)
	if err != nil {
		return accessToken, errors.New("SessionTokenGenerationError")
	}

	return entity.NewAccessToken(
		tokenType,
		expiresIn,
		tokenStr,
	), nil
}
