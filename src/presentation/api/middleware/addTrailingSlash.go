package apiMiddleware

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func AddTrailingSlash(apiBasePath string) echo.MiddlewareFunc {
	urlSkipRegex := regexp.MustCompile(
		`^` + apiBasePath + `(/v\d{1,2}/(auth|health)|/swagger)`,
	)

	return middleware.AddTrailingSlashWithConfig(middleware.TrailingSlashConfig{
		RedirectCode: http.StatusTemporaryRedirect,
		Skipper: func(c echo.Context) bool {
			urlPath := c.Request().URL.Path
			isNotApi := !strings.HasPrefix(urlPath, apiBasePath)

			return isNotApi || urlSkipRegex.MatchString(urlPath)
		},
	})
}
