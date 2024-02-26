package apiMiddleware

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func AddTrailingSlash() echo.MiddlewareFunc {
	urlSkipRegex := regexp.MustCompile(
		`^/api/v\d{1,2}/(swagger|auth|health)`,
	)

	return middleware.AddTrailingSlashWithConfig(middleware.TrailingSlashConfig{
		RedirectCode: http.StatusTemporaryRedirect,
		Skipper: func(c echo.Context) bool {
			urlPath := c.Request().URL.Path
			isNotApi := !strings.HasPrefix(urlPath, "/_/api/")

			return isNotApi || urlSkipRegex.MatchString(urlPath)
		},
	})
}
