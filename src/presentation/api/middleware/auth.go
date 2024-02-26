package apiMiddleware

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
	accessToken valueObject.AccessTokenStr,
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

	accessTokenDetails, err := useCase.GetAccessTokenDetails(
		authQueryRepo,
		accessToken,
		trustedIps,
		ipAddress,
	)
	if err != nil {
		return valueObject.AccountId(0), err
	}

	return accessTokenDetails.AccountId, nil
}

func Auth(apiBasePath string) echo.MiddlewareFunc {
	urlSkipRegex := regexp.MustCompile(
		`^` + apiBasePath + `/v\d{1,2}/(swagger|auth|health)`,
	)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			urlPath := c.Request().URL.Path
			isNotApi := !strings.HasPrefix(urlPath, apiBasePath)

			if isNotApi || urlSkipRegex.MatchString(urlPath) {
				return next(c)
			}

			token := c.Request().Header.Get("Authorization")
			if token == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, map[string]interface{}{
					"status": http.StatusUnauthorized,
					"body":   "MissingAuthToken",
				})
			}

			persistentDbSvc := c.Get("persistentDbSvc").(*db.PersistentDatabaseService)
			authQueryRepo := infra.NewAuthQueryRepo(persistentDbSvc)
			tokenWithoutPrefix := token[7:]
			accountId, err := getAccountIdFromAccessToken(
				authQueryRepo,
				valueObject.AccessTokenStr(tokenWithoutPrefix),
				valueObject.NewIpAddressPanic(c.RealIP()),
			)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, map[string]interface{}{
					"status": http.StatusUnauthorized,
					"body":   err.Error(),
				})
			}

			c.Set("accountId", accountId.String())
			return next(c)
		}
	}
}
