package apiController

import (
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/speedianet/control/src/domain/useCase"
	"github.com/speedianet/control/src/domain/valueObject"
	"github.com/speedianet/control/src/infra"
	"github.com/speedianet/control/src/infra/db"
	apiHelper "github.com/speedianet/control/src/presentation/api/helper"
	"github.com/speedianet/control/src/presentation/service"
)

type AccountController struct {
	accountService  *service.AccountService
	persistentDbSvc *db.PersistentDatabaseService
}

func NewAccountController(
	persistentDbSvc *db.PersistentDatabaseService,
) *AccountController {
	return &AccountController{
		accountService:  service.NewAccountService(persistentDbSvc),
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
// @Router       /v1/account/ [get]
func (controller *AccountController) Read(c echo.Context) error {
	return apiHelper.ServiceResponseWrapper(c, controller.accountService.Read())
}

func (controller *AccountController) accountQuotaFactory(
	quota interface{}, withDefaults bool,
) (accountQuota valueObject.AccountQuota, err error) {
	quotaMap, quotaMapOk := quota.(map[string]interface{})
	if !quotaMapOk {
		return valueObject.AccountQuota{}, errors.New("InvalidQuotaStructure")
	}

	accountQuota = valueObject.NewAccountQuotaWithDefaultValues()
	if !withDefaults {
		accountQuota = valueObject.NewAccountQuotaWithBlankValues()
	}

	cpuCores := accountQuota.CpuCores
	if quotaMap["cpuCores"] != nil {
		cpuCores, err = valueObject.NewCpuCoresCount(quotaMap["cpuCores"])
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

	diskBytes := accountQuota.DiskBytes
	if quotaMap["diskBytes"] != nil {
		diskBytes, err = valueObject.NewByte(quotaMap["diskBytes"])
		if err != nil {
			return accountQuota, err
		}
	}

	inodes := accountQuota.Inodes
	if quotaMap["inodes"] != nil {
		inodes, err = valueObject.NewInodesCount(quotaMap["inodes"])
		if err != nil {
			return accountQuota, err
		}
	}

	return valueObject.NewAccountQuota(cpuCores, memoryBytes, diskBytes, inodes), nil
}

// CreateAccount godoc
// @Summary      CreateAccount
// @Description  Create a new account.
// @Tags         account
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        createDto 	  body    dto.CreateAccount  true  "NewAccount"
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

	requestBody["ipAddress"] = c.RealIP()

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
// @Param        updateDto 	  body dto.UpdateAccount  true  "UpdateAccount (Only accountId is required.)"
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

	requestBody["ipAddress"] = c.RealIP()

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
	requestBody := map[string]interface{}{
		"accountId": c.Param("accountId"),
		"ipAddress": c.RealIP(),
	}

	return apiHelper.ServiceResponseWrapper(
		c, controller.accountService.Delete(requestBody),
	)
}

func (controller *AccountController) AutoUpdateAccountsQuotaUsage() {
	taskInterval := time.Duration(15) * time.Minute
	timer := time.NewTicker(taskInterval)
	defer timer.Stop()

	accountQueryRepo := infra.NewAccountQueryRepo(controller.persistentDbSvc)
	accountCmdRepo := infra.NewAccountCmdRepo(controller.persistentDbSvc)
	for range timer.C {
		useCase.AutoUpdateAccountsQuotaUsage(accountQueryRepo, accountCmdRepo)
	}
}
