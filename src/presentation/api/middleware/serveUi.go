package apiMiddleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/speedianet/control/src/presentation/ui"
)

func ServeUi() echo.MiddlewareFunc {
	uiFs := ui.UiFs()

	return middleware.StaticWithConfig(middleware.StaticConfig{
		HTML5:      true,
		Root:       "",
		Filesystem: uiFs,
	})
}
