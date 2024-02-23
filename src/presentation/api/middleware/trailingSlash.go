package apiMiddleware

import (
	"net/http"
	"regexp"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func TrailingSlash() echo.MiddlewareFunc {
	trailingSlashSkipRegex := regexp.MustCompile(
		`^/(v\d{1,2}/(swagger|auth|health)|_)`,
	)

	return middleware.AddTrailingSlashWithConfig(middleware.TrailingSlashConfig{
		RedirectCode: http.StatusTemporaryRedirect,
		Skipper: func(c echo.Context) bool {
			return trailingSlashSkipRegex.MatchString(c.Request().URL.Path)
		},
	})
}
