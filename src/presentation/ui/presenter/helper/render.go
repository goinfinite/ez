package presenterHelper

import (
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"

	"github.com/goinfinite/ez/src/presentation/ui/layout"
)

func Render(c echo.Context, pageContent templ.Component, statusCode int) error {
	c.Response().Writer.WriteHeader(statusCode)
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	currentUrl := c.Request().URL.String()

	return layout.
		MainLayout(pageContent, currentUrl).
		Render(c.Request().Context(), c.Response().Writer)
}
