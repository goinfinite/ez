package apiMiddleware

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func IsSkippableApiCall(req *http.Request, apiBasePath string) bool {
	urlPath := req.URL.Path
	isNotApi := !strings.HasPrefix(urlPath, apiBasePath)
	if isNotApi {
		return true
	}

	urlSkipRegex := regexp.MustCompile(
		`^` + apiBasePath + `(/v\d{1,2}/(auth|health)|/swagger)`,
	)
	return !urlSkipRegex.MatchString(urlPath)
}

func AddTrailingSlash(apiBasePath string) echo.MiddlewareFunc {
	return middleware.AddTrailingSlashWithConfig(middleware.TrailingSlashConfig{
		RedirectCode: http.StatusTemporaryRedirect,
		Skipper: func(c echo.Context) bool {
			return IsSkippableApiCall(c.Request(), apiBasePath)
		},
	})
}
