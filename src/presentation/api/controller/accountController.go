package apiController

import (
	"errors"
	"net/http"
	"time"

	"github.com/goinfinite/fleet/src/domain/dto"
	"github.com/goinfinite/fleet/src/domain/useCase"
	"github.com/goinfinite/fleet/src/domain/valueObject"
	"github.com/goinfinite/fleet/src/infra"
	"github.com/goinfinite/fleet/src/infra/db"
	apiHelper "github.com/goinfinite/fleet/src/presentation/api/helper"
	"github.com/labstack/echo/v4"
)

// GetAccounts	 godoc
// @Summary      GetAccounts
// @Description  List accs.
// @Tags         account
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200 {array} entity.Account
// @Router       /account/ [get]
func GetAccountsController(c echo.Context) error {
	dbSvc := c.Get("dbSvc").(*db.DatabaseService)
	accsQueryRepo := infra.NewAccQueryRepo(dbSvc)
	accsList, err := useCase.GetAccounts(accsQueryRepo)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, accsList)
}

func accQuotaFactory(
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

	return valueObject.NewAccountQuota(
		cpuCores,
		memoryBytes,
		diskBytes,
		inodes,
	), nil
}

// AddAccount	 godoc
// @Summary      AddNewAccount
// @Description  Add a new account.
// @Tags         account
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        addAccountDto 	  body    dto.AddAccount  true  "NewAccount"
// @Success      201 {object} object{} "AccountCreated"
// @Router       /account/ [post]
func AddAccountController(c echo.Context) error {
	requiredParams := []string{"username", "password"}
	requestBody, _ := apiHelper.GetRequestBody(c)

	apiHelper.CheckMissingParams(requestBody, requiredParams)

	var quotaPtr *valueObject.AccountQuota
	if requestBody["quota"] != nil {
		quota, err := accQuotaFactory(requestBody["quota"], true)
		if err != nil {
			return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
		}
		quotaPtr = &quota
	}

	addAccountDto := dto.NewAddAccount(
		valueObject.NewUsernamePanic(requestBody["username"].(string)),
		valueObject.NewPasswordPanic(requestBody["password"].(string)),
		quotaPtr,
	)

	dbSvc := c.Get("dbSvc").(*db.DatabaseService)
	accQueryRepo := infra.NewAccQueryRepo(dbSvc)
	accCmdRepo := infra.NewAccCmdRepo(dbSvc)

	err := useCase.AddAccount(
		accQueryRepo,
		accCmdRepo,
		addAccountDto,
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
// @Param        updateAccountDto 	  body dto.UpdateAccount  true  "UpdateAccount (Only accountId is required.)"
// @Success      200 {object} object{} "AccountUpdated message or NewKeyString"
// @Router       /account/ [put]
func UpdateAccountController(c echo.Context) error {
	requiredParams := []string{"accountId"}
	requestBody, _ := apiHelper.GetRequestBody(c)

	apiHelper.CheckMissingParams(requestBody, requiredParams)

	accountId := valueObject.NewAccountIdPanic(requestBody["accountId"])

	var passPtr *valueObject.Password
	if requestBody["password"] != nil {
		password := valueObject.NewPasswordPanic(requestBody["password"].(string))
		passPtr = &password
	}

	var shouldUpdateApiKeyPtr *bool
	if requestBody["shouldUpdateApiKey"] != nil {
		shouldUpdateApiKey := requestBody["shouldUpdateApiKey"].(bool)
		shouldUpdateApiKeyPtr = &shouldUpdateApiKey
	}

	var quotaPtr *valueObject.AccountQuota
	if requestBody["quota"] != nil {
		quota, err := accQuotaFactory(requestBody["quota"], false)
		if err != nil {
			return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
		}
		quotaPtr = &quota
	}

	updateAccountDto := dto.NewUpdateAccount(
		accountId,
		passPtr,
		shouldUpdateApiKeyPtr,
		quotaPtr,
	)

	dbSvc := c.Get("dbSvc").(*db.DatabaseService)
	accQueryRepo := infra.NewAccQueryRepo(dbSvc)
	accCmdRepo := infra.NewAccCmdRepo(dbSvc)

	if updateAccountDto.ShouldUpdateApiKey != nil && *updateAccountDto.ShouldUpdateApiKey {
		newKey, err := useCase.UpdateAccountApiKey(
			accQueryRepo,
			accCmdRepo,
			updateAccountDto,
		)
		if err != nil {
			return apiHelper.ResponseWrapper(
				c, http.StatusInternalServerError, err.Error(),
			)
		}

		return apiHelper.ResponseWrapper(c, http.StatusOK, newKey)
	}

	err := useCase.UpdateAccount(
		accQueryRepo,
		accCmdRepo,
		updateAccountDto,
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
// @Router       /account/{accountId}/ [delete]
func DeleteAccountController(c echo.Context) error {
	accountId := valueObject.NewAccountIdPanic(c.Param("accountId"))

	dbSvc := c.Get("dbSvc").(*db.DatabaseService)
	accQueryRepo := infra.NewAccQueryRepo(dbSvc)
	accCmdRepo := infra.NewAccCmdRepo(dbSvc)
	containerQueryRepo := infra.NewContainerQueryRepo(dbSvc)

	err := useCase.DeleteAccount(
		accQueryRepo,
		accCmdRepo,
		accountId,
		containerQueryRepo,
	)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, "AccountDeleted")
}

func AutoUpdateAccountsQuotaUsageController(dbSvc *db.DatabaseService) {
	taskInterval := time.Duration(15) * time.Minute
	timer := time.NewTicker(taskInterval)
	defer timer.Stop()

	accQueryRepo := infra.NewAccQueryRepo(dbSvc)
	accCmdRepo := infra.NewAccCmdRepo(dbSvc)
	for range timer.C {
		useCase.AutoUpdateAccountsQuotaUsage(accQueryRepo, accCmdRepo)
	}
}
