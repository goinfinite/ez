package apiController

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/speedianet/sfm/src/domain/dto"
	"github.com/speedianet/sfm/src/domain/useCase"
	"github.com/speedianet/sfm/src/domain/valueObject"
	"github.com/speedianet/sfm/src/infra"
	apiHelper "github.com/speedianet/sfm/src/presentation/api/helper"
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
	accsQueryRepo := infra.AccQueryRepo{}
	accsList, err := useCase.GetAccounts(accsQueryRepo)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, accsList)
}

func accQuotaFactory(quota interface{}) (valueObject.AccountQuota, error) {
	quotaMap, quotaMapOk := quota.(map[string]interface{})
	if !quotaMapOk {
		return valueObject.AccountQuota{}, errors.New("InvalidQuotaStructure")
	}

	accQuotaDefaults := valueObject.NewAccountQuotaWithDefaultValues()

	cpuCores := accQuotaDefaults.CpuCores
	if quotaMap["cpuCores"] != nil {
		cpuCores = valueObject.NewCpuCoresCountPanic(quotaMap["cpuCores"])
	}

	memoryBytes := accQuotaDefaults.MemoryBytes
	if quotaMap["memoryBytes"] != nil {
		memoryBytes = valueObject.NewBytePanic(quotaMap["memoryBytes"])
	}

	diskBytes := accQuotaDefaults.DiskBytes
	if quotaMap["diskBytes"] != nil {
		diskBytes = valueObject.NewBytePanic(quotaMap["diskBytes"])
	}

	inodes := accQuotaDefaults.Inodes
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
		quota, err := accQuotaFactory(requestBody["quota"])
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

	accQueryRepo := infra.AccQueryRepo{}
	accCmdRepo := infra.AccCmdRepo{}

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

	accQueryRepo := infra.AccQueryRepo{}
	accCmdRepo := infra.AccCmdRepo{}

	err := useCase.DeleteAccount(
		accQueryRepo,
		accCmdRepo,
		accountId,
	)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, "AccountDeleted")
}

// UpdateAccount godoc
// @Summary      UpdateAccount
// @Description  Update an account.
// @Tags         account
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        updateAccountDto 	  body dto.UpdateAccount  true  "UpdateAccount"
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

	updateAccountDto := dto.NewUpdateAccount(
		accountId,
		passPtr,
		shouldUpdateApiKeyPtr,
	)

	accQueryRepo := infra.AccQueryRepo{}
	accCmdRepo := infra.AccCmdRepo{}

	if updateAccountDto.Password != nil {
		useCase.UpdateAccountPassword(
			accQueryRepo,
			accCmdRepo,
			updateAccountDto,
		)
	}

	if updateAccountDto.ShouldUpdateApiKey != nil && *updateAccountDto.ShouldUpdateApiKey {
		newKey, err := useCase.UpdateAccountApiKey(
			accQueryRepo,
			accCmdRepo,
			updateAccountDto,
		)
		if err != nil {
			return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
		}

		return apiHelper.ResponseWrapper(c, http.StatusOK, newKey)
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, "AccountUpdated")
}
