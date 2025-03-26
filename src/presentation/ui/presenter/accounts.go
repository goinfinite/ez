package presenter

import (
	"net/http"

	"github.com/goinfinite/ez/src/domain/entity"
	"github.com/goinfinite/ez/src/infra/db"
	"github.com/goinfinite/ez/src/presentation/service"
	"github.com/goinfinite/ez/src/presentation/ui/page"
	presenterHelper "github.com/goinfinite/ez/src/presentation/ui/presenter/helper"
	"github.com/labstack/echo/v4"
)

type AccountsPresenter struct {
	accountService *service.AccountService
}

func NewAccountsPresenter(
	persistentDbSvc *db.PersistentDatabaseService,
	trailDbSvc *db.TrailDatabaseService,
) *AccountsPresenter {
	return &AccountsPresenter{
		accountService: service.NewAccountService(persistentDbSvc, trailDbSvc),
	}
}

func (presenter *AccountsPresenter) Handler(c echo.Context) error {
	responseOutput := presenter.accountService.Read()
	if responseOutput.Status != service.Success {
		return nil
	}

	accountList, assertOk := responseOutput.Body.([]entity.Account)
	if !assertOk {
		return nil
	}

	pageContent := page.AccountsIndex(accountList)
	return presenterHelper.Render(c, pageContent, http.StatusOK)
}
