package presenter

import (
	"maps"
	"net/http"

	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/valueObject"
	"github.com/goinfinite/ez/src/infra/db"
	"github.com/goinfinite/ez/src/presentation/service"
	uiHelper "github.com/goinfinite/ez/src/presentation/ui/helper"
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
	entitiesNamePrefix := "accounts"
	paginationMap := uiHelper.PaginationParser(c, entitiesNamePrefix, "id")
	if c.QueryParam(entitiesNamePrefix+"SortDirection") == "" {
		paginationMap["sortDirection"] = valueObject.PaginationSortDirectionDesc.String()
	}
	requestParamsMap := uiHelper.ReadRequestParser(
		c, entitiesNamePrefix, dto.ReadAccountsRequest{},
	)
	serviceRequestBody := paginationMap
	maps.Copy(serviceRequestBody, requestParamsMap)

	responseOutput := presenter.accountService.Read(serviceRequestBody)
	if responseOutput.Status != service.Success {
		return nil
	}

	readAccountsResponseDto, assertOk := responseOutput.Body.(dto.ReadAccountsResponse)
	if !assertOk {
		return nil
	}

	pageContent := page.AccountsIndex(readAccountsResponseDto.Accounts)
	return presenterHelper.Render(c, pageContent, http.StatusOK)
}
