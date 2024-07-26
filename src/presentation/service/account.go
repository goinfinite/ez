package service

import (
	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/useCase"
	"github.com/speedianet/control/src/domain/valueObject"
	"github.com/speedianet/control/src/infra"
	"github.com/speedianet/control/src/infra/db"
	serviceHelper "github.com/speedianet/control/src/presentation/service/helper"
)

var LocalOperatorAccountId, _ = valueObject.NewAccountId(0)
var LocalOperatorIpAddress = valueObject.NewLocalhostIpAddress()

type AccountService struct {
	persistentDbSvc       *db.PersistentDatabaseService
	accountQueryRepo      *infra.AccountQueryRepo
	accountCmdRepo        *infra.AccountCmdRepo
	activityRecordCmdRepo *infra.ActivityRecordCmdRepo
}

func NewAccountService(
	persistentDbSvc *db.PersistentDatabaseService,
	trailDbSvc *db.TrailDatabaseService,
) *AccountService {
	return &AccountService{
		persistentDbSvc:       persistentDbSvc,
		accountQueryRepo:      infra.NewAccountQueryRepo(persistentDbSvc),
		accountCmdRepo:        infra.NewAccountCmdRepo(persistentDbSvc),
		activityRecordCmdRepo: infra.NewActivityRecordCmdRepo(trailDbSvc),
	}
}

func (service *AccountService) Read() ServiceOutput {
	accountsList, err := useCase.ReadAccounts(service.accountQueryRepo)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, accountsList)
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

	operatorAccountId := LocalOperatorAccountId
	if _, exists := input["operatorAccountId"]; exists {
		operatorAccountId, err = valueObject.NewAccountId(input["operatorAccountId"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	ipAddress := LocalOperatorIpAddress
	if _, exists := input["ipAddress"]; exists {
		ipAddress, err = valueObject.NewIpAddress(input["ipAddress"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	createDto := dto.NewCreateAccount(
		username, password, quotaPtr, operatorAccountId, ipAddress,
	)

	err = useCase.CreateAccount(
		service.accountQueryRepo, service.accountCmdRepo,
		service.activityRecordCmdRepo, createDto,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Created, "AccountCreated")
}

func (service *AccountService) Update(input map[string]interface{}) ServiceOutput {
	if _, exists := input["id"]; exists {
		input["accountId"] = input["id"]
	}

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

	operatorAccountId := LocalOperatorAccountId
	if _, exists := input["operatorAccountId"]; exists {
		operatorAccountId, err = valueObject.NewAccountId(input["operatorAccountId"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	ipAddress := LocalOperatorIpAddress
	if _, exists := input["ipAddress"]; exists {
		ipAddress, err = valueObject.NewIpAddress(input["ipAddress"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	updateDto := dto.NewUpdateAccount(
		accountId, passwordPtr, shouldUpdateApiKeyPtr, quotaPtr,
		operatorAccountId, ipAddress,
	)

	if updateDto.ShouldUpdateApiKey != nil && *updateDto.ShouldUpdateApiKey {
		newKey, err := useCase.UpdateAccountApiKey(
			service.accountQueryRepo, service.accountCmdRepo,
			service.activityRecordCmdRepo, updateDto,
		)
		if err != nil {
			return NewServiceOutput(InfraError, err.Error())
		}
		return NewServiceOutput(Success, newKey)
	}

	err = useCase.UpdateAccount(
		service.accountQueryRepo, service.accountCmdRepo,
		service.activityRecordCmdRepo, updateDto,
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

	operatorAccountId := LocalOperatorAccountId
	if _, exists := input["operatorAccountId"]; exists {
		operatorAccountId, err = valueObject.NewAccountId(input["operatorAccountId"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	ipAddress := LocalOperatorIpAddress
	if _, exists := input["ipAddress"]; exists {
		ipAddress, err = valueObject.NewIpAddress(input["ipAddress"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	deleteDto := dto.NewDeleteAccount(accountId, operatorAccountId, ipAddress)

	containerQueryRepo := infra.NewContainerQueryRepo(service.persistentDbSvc)

	err = useCase.DeleteAccount(
		service.accountQueryRepo, service.accountCmdRepo, containerQueryRepo,
		service.activityRecordCmdRepo, deleteDto,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, "AccountDeleted")
}
