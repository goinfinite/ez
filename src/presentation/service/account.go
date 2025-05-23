package service

import (
	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/useCase"
	"github.com/goinfinite/ez/src/domain/valueObject"
	voHelper "github.com/goinfinite/ez/src/domain/valueObject/helper"
	"github.com/goinfinite/ez/src/infra"
	"github.com/goinfinite/ez/src/infra/db"
	serviceHelper "github.com/goinfinite/ez/src/presentation/service/helper"
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

func (service *AccountService) ReadAccountsRequestFactory(
	serviceInput map[string]any,
) (readRequestDto dto.ReadAccountsRequest, err error) {
	if serviceInput["accountId"] != nil {
		serviceInput["id"] = serviceInput["accountId"]
	}

	if serviceInput["accountUsername"] != nil {
		serviceInput["username"] = serviceInput["accountUsername"]
	}

	var accountIdPtr *valueObject.AccountId
	if serviceInput["id"] != nil {
		accountId, err := valueObject.NewAccountId(serviceInput["id"])
		if err != nil {
			return readRequestDto, err
		}
		accountIdPtr = &accountId
	}

	var accountUsernamePtr *valueObject.UnixUsername
	if serviceInput["username"] != nil {
		accountUsername, err := valueObject.NewUnixUsername(serviceInput["username"])
		if err != nil {
			return readRequestDto, err
		}
		accountUsernamePtr = &accountUsername
	}

	requestPagination, err := serviceHelper.PaginationParser(
		serviceInput, useCase.AccountsDefaultPagination,
	)
	if err != nil {
		return readRequestDto, err
	}

	return dto.ReadAccountsRequest{
		Pagination:      requestPagination,
		AccountId:       accountIdPtr,
		AccountUsername: accountUsernamePtr,
	}, nil
}

func (service *AccountService) Read(input map[string]any) ServiceOutput {
	readAccountsRequestDto, err := service.ReadAccountsRequestFactory(input)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}
	accountsResponseDto, err := useCase.ReadAccounts(
		service.accountQueryRepo, readAccountsRequestDto,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, accountsResponseDto)
}

func (service *AccountService) Create(input map[string]interface{}) ServiceOutput {
	requiredParams := []string{"username", "password"}
	err := serviceHelper.RequiredParamsInspector(input, requiredParams)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	username, err := valueObject.NewUnixUsername(input["username"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	password, err := valueObject.NewPassword(input["password"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	var quotaPtr *valueObject.AccountQuota
	if input["quota"] != nil {
		accountQuota, assertOk := input["quota"].(valueObject.AccountQuota)
		if !assertOk {
			return NewServiceOutput(UserError, "InvalidQuota")
		}
		quotaPtr = &accountQuota
	}

	operatorAccountId := LocalOperatorAccountId
	if input["operatorAccountId"] != nil {
		operatorAccountId, err = valueObject.NewAccountId(input["operatorAccountId"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	operatorIpAddress := LocalOperatorIpAddress
	if input["operatorIpAddress"] != nil {
		operatorIpAddress, err = valueObject.NewIpAddress(input["operatorIpAddress"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	createDto := dto.NewCreateAccount(
		username, password, quotaPtr, operatorAccountId, operatorIpAddress,
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
	if input["id"] != nil {
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
	if input["password"] != nil {
		password, err := valueObject.NewPassword(input["password"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		passwordPtr = &password
	}

	var shouldUpdateApiKeyPtr *bool
	if input["shouldUpdateApiKey"] != nil {
		shouldUpdateApiKey, err := voHelper.InterfaceToBool(input["shouldUpdateApiKey"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		shouldUpdateApiKeyPtr = &shouldUpdateApiKey
	}

	var quotaPtr *valueObject.AccountQuota
	if input["quota"] != nil {
		accountQuota, assertOk := input["quota"].(valueObject.AccountQuota)
		if !assertOk {
			return NewServiceOutput(UserError, "InvalidQuota")
		}
		quotaPtr = &accountQuota
	}

	operatorAccountId := LocalOperatorAccountId
	if input["operatorAccountId"] != nil {
		operatorAccountId, err = valueObject.NewAccountId(input["operatorAccountId"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	operatorIpAddress := LocalOperatorIpAddress
	if input["operatorIpAddress"] != nil {
		operatorIpAddress, err = valueObject.NewIpAddress(input["operatorIpAddress"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	updateDto := dto.NewUpdateAccount(
		accountId, passwordPtr, shouldUpdateApiKeyPtr, quotaPtr,
		operatorAccountId, operatorIpAddress,
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

func (service *AccountService) RefreshQuotas() ServiceOutput {
	err := useCase.RefreshAccountQuotas(
		service.accountQueryRepo, service.accountCmdRepo,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, "AccountQuotasRefreshed")
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
	if input["operatorAccountId"] != nil {
		operatorAccountId, err = valueObject.NewAccountId(input["operatorAccountId"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	operatorIpAddress := LocalOperatorIpAddress
	if input["operatorIpAddress"] != nil {
		operatorIpAddress, err = valueObject.NewIpAddress(input["operatorIpAddress"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	deleteDto := dto.NewDeleteAccount(accountId, operatorAccountId, operatorIpAddress)

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
