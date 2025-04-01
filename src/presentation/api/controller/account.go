package apiController

import (
	"errors"
	"net/http"
	"time"

	"github.com/goinfinite/ez/src/domain/useCase"
	"github.com/goinfinite/ez/src/domain/valueObject"
	voHelper "github.com/goinfinite/ez/src/domain/valueObject/helper"
	"github.com/goinfinite/ez/src/infra"
	"github.com/goinfinite/ez/src/infra/db"
	apiHelper "github.com/goinfinite/ez/src/presentation/api/helper"
	"github.com/goinfinite/ez/src/presentation/service"
	"github.com/labstack/echo/v4"
)

type AccountController struct {
	accountService  *service.AccountService
	persistentDbSvc *db.PersistentDatabaseService
}

func NewAccountController(
	persistentDbSvc *db.PersistentDatabaseService,
	trailDbSvc *db.TrailDatabaseService,
) *AccountController {
	return &AccountController{
		accountService:  service.NewAccountService(persistentDbSvc, trailDbSvc),
		persistentDbSvc: persistentDbSvc,
	}
}

// ReadAccounts  godoc
// @Summary      ReadAccounts
// @Description  List accounts.
// @Tags         account
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200 {array} entity.Account
// @Param        id query  string  false  "AccountId"
// @Param        username query  string  false  "AccountUsername"
// @Param        pageNumber query  uint  false  "PageNumber (Pagination)"
// @Param        itemsPerPage query  uint  false  "ItemsPerPage (Pagination)"
// @Param        sortBy query  string  false  "SortBy (Pagination)"
// @Param        sortDirection query  string  false  "SortDirection (Pagination)"
// @Param        lastSeenId query  string  false  "LastSeenId (Pagination)"
// @Router       /v1/account/ [get]
func (controller *AccountController) Read(c echo.Context) error {
	requestBody, err := apiHelper.ReadRequestBody(c)
	if err != nil {
		return err
	}

	return apiHelper.ServiceResponseWrapper(
		c, controller.accountService.Read(requestBody),
	)
}

func (controller *AccountController) accountQuotaFactory(
	quota interface{}, withDefaults bool,
) (accountQuota valueObject.AccountQuota, err error) {
	quotaMap, quotaMapOk := quota.(map[string]interface{})
	if !quotaMapOk {
		return accountQuota, errors.New("InvalidQuotaStructure")
	}

	accountQuota = valueObject.NewAccountQuotaWithDefaultValues()
	if !withDefaults {
		accountQuota = valueObject.NewAccountQuotaWithBlankValues()
	}

	millicores := accountQuota.Millicores
	if quotaMap["millicores"] != nil {
		millicores, err = valueObject.NewMillicores(quotaMap["millicores"])
		if err != nil {
			return accountQuota, err
		}
	}

	if quotaMap["cpuCores"] != nil {
		cpuCoresUint, err := voHelper.InterfaceToFloat64(quotaMap["cpuCores"])
		if err != nil {
			return accountQuota, err
		}

		millicores, err = valueObject.NewMillicores(cpuCoresUint * 1000)
		if err != nil {
			return accountQuota, err
		}
	}

	memoryBytes := accountQuota.MemoryBytes
	if quotaMap["memoryBytes"] != nil {
		memoryBytes, err = valueObject.NewByte(quotaMap["memoryBytes"])
		if err != nil {
			return accountQuota, err
		}
	}

	if quotaMap["memoryMebibytes"] != nil {
		memoryBytes, err = valueObject.NewMebibyte(quotaMap["memoryMebibytes"])
		if err != nil {
			return accountQuota, err
		}
	}

	if quotaMap["memoryGibibytes"] != nil {
		memoryBytes, err = valueObject.NewGibibyte(quotaMap["memoryGibibytes"])
		if err != nil {
			return accountQuota, err
		}
	}

	storageBytes := accountQuota.StorageBytes
	if quotaMap["diskBytes"] != nil {
		quotaMap["storageBytes"] = quotaMap["diskBytes"]
	}

	if quotaMap["storageBytes"] != nil {
		storageBytes, err = valueObject.NewByte(quotaMap["storageBytes"])
		if err != nil {
			return accountQuota, err
		}
	}

	if quotaMap["diskMebibytes"] != nil {
		quotaMap["storageMebibytes"] = quotaMap["diskMebibytes"]
	}

	if quotaMap["storageMebibytes"] != nil {
		storageBytes, err = valueObject.NewMebibyte(quotaMap["storageMebibytes"])
		if err != nil {
			return accountQuota, err
		}
	}

	if quotaMap["diskGibibytes"] != nil {
		quotaMap["storageGibibytes"] = quotaMap["diskGibibytes"]
	}

	if quotaMap["storageGibibytes"] != nil {
		storageBytes, err = valueObject.NewGibibyte(quotaMap["storageGibibytes"])
		if err != nil {
			return accountQuota, err
		}
	}

	storageInodes := accountQuota.StorageInodes
	if quotaMap["inodes"] != nil {
		quotaMap["storageInodes"] = quotaMap["inodes"]
	}

	if quotaMap["storageInodes"] != nil {
		storageInodes, err = voHelper.InterfaceToUint64(quotaMap["storageInodes"])
		if err != nil {
			return accountQuota, errors.New("InvalidStorageInodes")
		}
	}

	storagePerformanceUnits := accountQuota.StoragePerformanceUnits
	if quotaMap["storagePerformanceUnits"] != nil {
		storagePerformanceUnits, err = valueObject.NewStoragePerformanceUnits(
			quotaMap["storagePerformanceUnits"],
		)
		if err != nil {
			return accountQuota, err
		}
	}

	return valueObject.NewAccountQuota(
		millicores, memoryBytes, storageBytes, storageInodes, storagePerformanceUnits,
	), nil
}

// CreateAccount godoc
// @Summary      CreateAccount
// @Description  Create a new account.
// @Tags         account
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        createDto 	  body    dto.CreateAccount  true  "Human-readable fields ('cpuCores', 'memoryMebibytes' etc) will be converted to their technical counterpart ('millicores' etc) automatically."
// @Success      201 {object} object{} "AccountCreated"
// @Router       /v1/account/ [post]
func (controller *AccountController) Create(c echo.Context) error {
	requestBody, err := apiHelper.ReadRequestBody(c)
	if err != nil {
		return err
	}

	if requestBody["quota"] != nil {
		quota, err := controller.accountQuotaFactory(requestBody["quota"], true)
		if err != nil {
			return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
		}
		requestBody["quota"] = quota
	}

	return apiHelper.ServiceResponseWrapper(
		c, controller.accountService.Create(requestBody),
	)
}

// UpdateAccount godoc
// @Summary      UpdateAccount
// @Description  Update an account.
// @Tags         account
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        updateDto 	  body dto.UpdateAccount  true  "Only 'accountId' is required. Human-readable fields ('cpuCores', 'memoryMebibytes' etc) will be converted to their technical counterpart ('millicores' etc) automatically."
// @Success      200 {object} object{} "AccountUpdated message or NewKeyString"
// @Router       /v1/account/ [put]
func (controller *AccountController) Update(c echo.Context) error {
	requestBody, err := apiHelper.ReadRequestBody(c)
	if err != nil {
		return err
	}

	if requestBody["quota"] != nil {
		quota, err := controller.accountQuotaFactory(requestBody["quota"], true)
		if err != nil {
			return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
		}
		requestBody["quota"] = quota
	}

	return apiHelper.ServiceResponseWrapper(
		c, controller.accountService.Update(requestBody),
	)
}

// DeleteAccount godoc
// @Summary      DeleteAccount
// @Description  Delete an account.
// @Tags         account
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        accountId 	  path   string  true  "AccountId"
// @Success      200 {object} object{} "AccountDeleted"
// @Router       /v1/account/{accountId}/ [delete]
func (controller *AccountController) Delete(c echo.Context) error {
	requestBody, err := apiHelper.ReadRequestBody(c)
	if err != nil {
		return err
	}

	return apiHelper.ServiceResponseWrapper(
		c, controller.accountService.Delete(requestBody),
	)
}

func (controller *AccountController) AutoRefreshAccountQuotas() {
	taskInterval := time.Duration(useCase.AutoRefreshAccountQuotasTimeIntervalSecs) * time.Second
	timer := time.NewTicker(taskInterval)
	defer timer.Stop()

	accountQueryRepo := infra.NewAccountQueryRepo(controller.persistentDbSvc)
	accountCmdRepo := infra.NewAccountCmdRepo(controller.persistentDbSvc)
	for range timer.C {
		_ = useCase.RefreshAccountQuotas(accountQueryRepo, accountCmdRepo)
	}
}
