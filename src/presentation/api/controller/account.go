package apiController

import (
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/useCase"
	"github.com/speedianet/control/src/domain/valueObject"
	"github.com/speedianet/control/src/infra"
	"github.com/speedianet/control/src/infra/db"
	apiHelper "github.com/speedianet/control/src/presentation/api/helper"
)

type AccountController struct {
	persistentDbSvc *db.PersistentDatabaseService
}

func NewAccountController(
	persistentDbSvc *db.PersistentDatabaseService,
) *AccountController {
	return &AccountController{persistentDbSvc: persistentDbSvc}
}

// GetAccounts	 godoc
// @Summary      ReadAccounts
// @Description  List accounts.
// @Tags         account
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200 {array} entity.Account
// @Router       /v1/account/ [get]
func (controller *AccountController) Read(c echo.Context) error {
	accountQueryRepo := infra.NewAccountQueryRepo(controller.persistentDbSvc)
	accsList, err := useCase.GetAccounts(accountQueryRepo)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, accsList)
}

func (controller *AccountController) accQuotaFactory(
	quota interface{},
	withDefaults bool,
) (valueObject.AccountQuota, error) {
	quotaMap, quotaMapOk := quota.(map[string]interface{})
	if !quotaMapOk {
		return valueObject.AccountQuota{}, errors.New("InvalidQuotaStructure")
	}

	accQuota := valueObject.NewAccountQuotaWithDefaultValues()
	if !withDefaults {
		accQuota = valueObject.NewAccountQuotaWithBlankValues()
	}

	cpuCores := accQuota.CpuCores
	if quotaMap["cpuCores"] != nil {
		cpuCores = valueObject.NewCpuCoresCountPanic(quotaMap["cpuCores"])
	}

	memoryBytes := accQuota.MemoryBytes
	if quotaMap["memoryBytes"] != nil {
		memoryBytes = valueObject.NewBytePanic(quotaMap["memoryBytes"])
	}

	diskBytes := accQuota.DiskBytes
	if quotaMap["diskBytes"] != nil {
		diskBytes = valueObject.NewBytePanic(quotaMap["diskBytes"])
	}

	inodes := accQuota.Inodes
	if quotaMap["inodes"] != nil {
		inodes = valueObject.NewInodesCountPanic(quotaMap["inodes"])
	}

	return valueObject.NewAccountQuota(cpuCores, memoryBytes, diskBytes, inodes), nil
}

// CreateAccount	 godoc
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

	requiredParams := []string{"username", "password"}
	apiHelper.CheckMissingParams(requestBody, requiredParams)

	var quotaPtr *valueObject.AccountQuota
	if requestBody["quota"] != nil {
		quota, err := controller.accQuotaFactory(requestBody["quota"], true)
		if err != nil {
			return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
		}
		quotaPtr = &quota
	}

	username, err := valueObject.NewUsername(requestBody["username"])
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
	}

	password, err := valueObject.NewPassword(requestBody["password"])
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
	}

	ipAddress, err := valueObject.NewIpAddress(c.RealIP())
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
	}

	createDto := dto.NewCreateAccount(username, password, quotaPtr, &ipAddress)

	accountQueryRepo := infra.NewAccountQueryRepo(controller.persistentDbSvc)
	accountCmdRepo := infra.NewAccountCmdRepo(controller.persistentDbSvc)
	securityCmdRepo := infra.NewSecurityCmdRepo(controller.persistentDbSvc)

	err = useCase.CreateAccount(
		accountQueryRepo, accountCmdRepo, securityCmdRepo, createDto,
	)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusCreated, "AccountCreated")
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

	requiredParams := []string{"accountId"}
	apiHelper.CheckMissingParams(requestBody, requiredParams)

	accountId := valueObject.NewAccountIdPanic(requestBody["accountId"])

	var passPtr *valueObject.Password
	if requestBody["password"] != nil {
		password, err := valueObject.NewPassword(requestBody["password"])
		if err != nil {
			return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
		}
		passPtr = &password
	}

	var shouldUpdateApiKeyPtr *bool
	if requestBody["shouldUpdateApiKey"] != nil {
		shouldUpdateApiKey := requestBody["shouldUpdateApiKey"].(bool)
		shouldUpdateApiKeyPtr = &shouldUpdateApiKey
	}

	var quotaPtr *valueObject.AccountQuota
	if requestBody["quota"] != nil {
		quota, err := controller.accQuotaFactory(requestBody["quota"], false)
		if err != nil {
			return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
		}
		quotaPtr = &quota
	}

	ipAddress, err := valueObject.NewIpAddress(c.RealIP())
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
	}

	updateDto := dto.NewUpdateAccount(
		accountId, passPtr, shouldUpdateApiKeyPtr, quotaPtr, &ipAddress,
	)

	accountQueryRepo := infra.NewAccountQueryRepo(controller.persistentDbSvc)
	accountCmdRepo := infra.NewAccountCmdRepo(controller.persistentDbSvc)
	securityCmdRepo := infra.NewSecurityCmdRepo(controller.persistentDbSvc)

	if updateDto.ShouldUpdateApiKey != nil && *updateDto.ShouldUpdateApiKey {
		newKey, err := useCase.UpdateAccountApiKey(
			accountQueryRepo, accountCmdRepo, securityCmdRepo, updateDto,
		)
		if err != nil {
			return apiHelper.ResponseWrapper(
				c, http.StatusInternalServerError, err.Error(),
			)
		}

		return apiHelper.ResponseWrapper(c, http.StatusOK, newKey)
	}

	err = useCase.UpdateAccount(
		accountQueryRepo, accountCmdRepo, securityCmdRepo, updateDto,
	)
	if err != nil {
		return apiHelper.ResponseWrapper(
			c, http.StatusInternalServerError, err.Error(),
		)
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, "AccountUpdated")
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
	accountId := valueObject.NewAccountIdPanic(c.Param("accountId"))

	accountQueryRepo := infra.NewAccountQueryRepo(controller.persistentDbSvc)
	accountCmdRepo := infra.NewAccountCmdRepo(controller.persistentDbSvc)
	containerQueryRepo := infra.NewContainerQueryRepo(controller.persistentDbSvc)
	securityCmdRepo := infra.NewSecurityCmdRepo(controller.persistentDbSvc)

	ipAddress, err := valueObject.NewIpAddress(c.RealIP())
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
	}

	err = useCase.DeleteAccount(
		accountQueryRepo, accountCmdRepo, containerQueryRepo,
		securityCmdRepo, accountId, &ipAddress,
	)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, "AccountDeleted")
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
