package presenter

import (
	"log/slog"
	"maps"
	"net/http"

	"github.com/goinfinite/ez/src/domain/dto"
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
	requestParamsMap := uiHelper.ReadRequestParser(
		c, entitiesNamePrefix, dto.ReadAccountsRequest{},
	)
	serviceRequestBody := paginationMap
	maps.Copy(serviceRequestBody, requestParamsMap)

	readAccountsRequestDto, err := presenter.accountService.ReadAccountsRequestFactory(
		serviceRequestBody,
	)
	if err != nil {
		slog.Debug("ReadAccountsRequestFactoryFailure", slog.Any("error", err))
		return nil
	}

	readAccountsServiceOutput := presenter.accountService.Read(serviceRequestBody)
	if readAccountsServiceOutput.Status != service.Success {
		slog.Debug("ReadAccountsFailure")
		return nil
	}

	readAccountsResponseDto, assertOk := readAccountsServiceOutput.Body.(dto.ReadAccountsResponse)
	if !assertOk {
		slog.Debug("ReadAccountsRequestTypeAssertFailure")
		return nil
	}

	pageContent := page.AccountsIndex(readAccountsRequestDto, readAccountsResponseDto)
	return presenterHelper.Render(c, pageContent, http.StatusOK)
}
