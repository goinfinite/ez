package presenter

import (
	"net/http"

	"github.com/goinfinite/ez/src/presentation/ui/layout"
	"github.com/labstack/echo/v4"
)

type LoginPresenter struct {
}

func NewLoginPresenter() *LoginPresenter {
	return &LoginPresenter{}
}

func (presenter *LoginPresenter) Handler(c echo.Context) error {
	c.Response().Writer.WriteHeader(http.StatusOK)
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)

	return layout.Login().
		Render(c.Request().Context(), c.Response().Writer)
}
