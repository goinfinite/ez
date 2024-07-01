package service

import (
	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/useCase"
	"github.com/speedianet/control/src/domain/valueObject"
	"github.com/speedianet/control/src/infra"
	"github.com/speedianet/control/src/infra/db"
	serviceHelper "github.com/speedianet/control/src/presentation/service/helper"
)

type AccountService struct {
	persistentDbSvc *db.PersistentDatabaseService
}

func NewAccountService(
	persistentDbSvc *db.PersistentDatabaseService,
) *AccountService {
	return &AccountService{persistentDbSvc: persistentDbSvc}
}

func (service *AccountService) Read() ServiceOutput {
	accountQueryRepo := infra.NewAccountQueryRepo(service.persistentDbSvc)
	accsList, err := useCase.GetAccounts(accountQueryRepo)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, accsList)
}

func (service *AccountService) Create(input map[string]interface{}) ServiceOutput {
	requiredParams := []string{"username", "password"}
	err := serviceHelper.RequiredParamsInspector(input, requiredParams)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	username, err := valueObject.NewUsername(input["username"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	password, err := valueObject.NewPassword(input["password"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	var quotaPtr *valueObject.AccountQuota
	if _, exists := input["quota"]; exists {
		accountQuota, assertOk := input["quota"].(valueObject.AccountQuota)
		if !assertOk {
			return NewServiceOutput(UserError, "InvalidQuota")
		}
		quotaPtr = &accountQuota
	}

	var ipAddressPtr *valueObject.IpAddress
	if _, exists := input["ipAddress"]; exists {
		ipAddress, err := valueObject.NewIpAddress(input["ipAddress"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		ipAddressPtr = &ipAddress
	}

	createDto := dto.NewCreateAccount(username, password, quotaPtr, ipAddressPtr)

	accountQueryRepo := infra.NewAccountQueryRepo(service.persistentDbSvc)
	accountCmdRepo := infra.NewAccountCmdRepo(service.persistentDbSvc)
	securityCmdRepo := infra.NewSecurityCmdRepo(service.persistentDbSvc)

	err = useCase.CreateAccount(
		accountQueryRepo, accountCmdRepo, securityCmdRepo, createDto,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Created, "AccountCreated")
}

func (service *AccountService) Update(input map[string]interface{}) ServiceOutput {
	requiredParams := []string{"accountId"}
	err := serviceHelper.RequiredParamsInspector(input, requiredParams)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	accountId, err := valueObject.NewAccountId(input["accountId"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	var passwordPtr *valueObject.Password
	if _, exists := input["password"]; exists {
		password, err := valueObject.NewPassword(input["password"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		passwordPtr = &password
	}

	var shouldUpdateApiKeyPtr *bool
	if _, exists := input["shouldUpdateApiKey"]; exists {
		shouldUpdateApiKey, assertOk := input["shouldUpdateApiKey"].(bool)
		if !assertOk {
			return NewServiceOutput(UserError, "InvalidShouldUpdateApiKey")
		}
		shouldUpdateApiKeyPtr = &shouldUpdateApiKey
	}

	var quotaPtr *valueObject.AccountQuota
	if _, exists := input["quota"]; exists {
		accountQuota, assertOk := input["quota"].(valueObject.AccountQuota)
		if !assertOk {
			return NewServiceOutput(UserError, "InvalidQuota")
		}
		quotaPtr = &accountQuota
	}

	var ipAddressPtr *valueObject.IpAddress
	if _, exists := input["ipAddress"]; exists {
		ipAddress, err := valueObject.NewIpAddress(input["ipAddress"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		ipAddressPtr = &ipAddress
	}

	updateDto := dto.NewUpdateAccount(
		accountId, passwordPtr, shouldUpdateApiKeyPtr, quotaPtr, ipAddressPtr,
	)

	accountQueryRepo := infra.NewAccountQueryRepo(service.persistentDbSvc)
	accountCmdRepo := infra.NewAccountCmdRepo(service.persistentDbSvc)
	securityCmdRepo := infra.NewSecurityCmdRepo(service.persistentDbSvc)

	if updateDto.ShouldUpdateApiKey != nil && *updateDto.ShouldUpdateApiKey {
		newKey, err := useCase.UpdateAccountApiKey(
			accountQueryRepo, accountCmdRepo, securityCmdRepo, updateDto,
		)
		if err != nil {
			return NewServiceOutput(InfraError, err.Error())
		}
		return NewServiceOutput(Success, newKey)
	}

	err = useCase.UpdateAccount(
		accountQueryRepo, accountCmdRepo, securityCmdRepo, updateDto,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, "AccountUpdated")
}

func (service *AccountService) Delete(input map[string]interface{}) ServiceOutput {
	requiredParams := []string{"accountId"}
	err := serviceHelper.RequiredParamsInspector(input, requiredParams)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	accountId, err := valueObject.NewAccountId(input["accountId"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	accountQueryRepo := infra.NewAccountQueryRepo(service.persistentDbSvc)
	accountCmdRepo := infra.NewAccountCmdRepo(service.persistentDbSvc)
	containerQueryRepo := infra.NewContainerQueryRepo(service.persistentDbSvc)
	securityCmdRepo := infra.NewSecurityCmdRepo(service.persistentDbSvc)

	var ipAddressPtr *valueObject.IpAddress
	if _, exists := input["ipAddress"]; exists {
		ipAddress, err := valueObject.NewIpAddress(input["ipAddress"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		ipAddressPtr = &ipAddress
	}

	err = useCase.DeleteAccount(
		accountQueryRepo, accountCmdRepo, containerQueryRepo,
		securityCmdRepo, accountId, ipAddressPtr,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, "AccountDeleted")
}
