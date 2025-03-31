package presenter

import (
	"github.com/labstack/echo/v4"

	"net/http"

	"github.com/goinfinite/ez/src/presentation/ui/page"
	presenterHelper "github.com/goinfinite/ez/src/presentation/ui/presenter/helper"
)

type MappingsPresenter struct {
}

func NewMappingsPresenter() *MappingsPresenter {
	return &MappingsPresenter{}
}

func (presenter *MappingsPresenter) Handler(c echo.Context) error {
	pageContent := page.MappingsIndex()
	return presenterHelper.Render(c, pageContent, http.StatusOK)
}
