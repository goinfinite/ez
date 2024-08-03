package uiMiddleware

import (
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/speedianet/control/src/domain/repository"
	"github.com/speedianet/control/src/domain/useCase"
	"github.com/speedianet/control/src/domain/valueObject"
	"github.com/speedianet/control/src/infra"
	"github.com/speedianet/control/src/infra/db"
)

func getAccountIdFromAccessToken(
	authQueryRepo repository.AuthQueryRepo,
	accessTokenValue valueObject.AccessTokenValue,
	ipAddress valueObject.IpAddress,
) (valueObject.AccountId, error) {
	trustedIpsRaw := strings.Split(os.Getenv("TRUSTED_IPS"), ",")
	var trustedIps []valueObject.IpAddress
	for _, trustedIp := range trustedIpsRaw {
		ipAddress, err := valueObject.NewIpAddress(trustedIp)
		if err != nil {
			continue
		}
		trustedIps = append(trustedIps, ipAddress)
	}

	accessTokenDetails, err := useCase.ReadAccessTokenDetails(
		authQueryRepo, accessTokenValue, trustedIps, ipAddress,
	)
	if err != nil {
		return valueObject.AccountId(0), err
	}

	return accessTokenDetails.AccountId, nil
}

func shouldSkipUiAuthentication(req *http.Request) bool {
	urlSkipRegex := regexp.MustCompile(`^/(api|\_|login)/`)
	return urlSkipRegex.MatchString(req.URL.Path)
}

func Authentication(
	persistentDbSvc *db.PersistentDatabaseService,
) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if shouldSkipUiAuthentication(c.Request()) {
				return next(c)
			}

			rawAccessToken := ""
			accessTokenCookie, err := c.Cookie("control-access-token")
			if err == nil {
				rawAccessToken = accessTokenCookie.Value
			}

			loginPath := "/_/#/login"

			if rawAccessToken == "" {
				rawAccessToken = c.Request().Header.Get("Authorization")
				if rawAccessToken == "" {
					return c.Redirect(http.StatusTemporaryRedirect, loginPath)
				}
				tokenWithoutPrefix := rawAccessToken[7:]
				rawAccessToken = tokenWithoutPrefix
			}

			accessTokenValue, err := valueObject.NewAccessTokenValue(rawAccessToken)
			if err != nil {
				return c.Redirect(http.StatusTemporaryRedirect, loginPath)
			}

			userIpAddress, err := valueObject.NewIpAddress(c.RealIP())
			if err != nil {
				return c.Redirect(http.StatusTemporaryRedirect, loginPath)
			}

			authQueryRepo := infra.NewAuthQueryRepo(persistentDbSvc)
			_, err = getAccountIdFromAccessToken(
				authQueryRepo, accessTokenValue, userIpAddress,
			)
			if err != nil {
				return c.Redirect(http.StatusTemporaryRedirect, loginPath)
			}
			return next(c)
		}
	}
}
