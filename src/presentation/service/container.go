package service

import (
	"strconv"

	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/useCase"
	"github.com/speedianet/control/src/domain/valueObject"
	"github.com/speedianet/control/src/infra"
	"github.com/speedianet/control/src/infra/db"
	serviceHelper "github.com/speedianet/control/src/presentation/service/helper"
)

type ContainerService struct {
	persistentDbSvc *db.PersistentDatabaseService
}

func NewContainerService(
	persistentDbSvc *db.PersistentDatabaseService,
) *ContainerService {
	return &ContainerService{
		persistentDbSvc: persistentDbSvc,
	}
}

func (service *ContainerService) Create(input map[string]interface{}) ServiceOutput {
	requiredParams := []string{"accountId", "hostname"}

	err := serviceHelper.RequiredParamsInspector(input, requiredParams)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	accId, err := valueObject.NewAccountId(input["accountId"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	hostname, err := valueObject.NewFqdn(input["hostname"].(string))
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	imgAddrStr, assertOk := input["imageAddress"].(string)
	if !assertOk {
		imgAddrStr, assertOk = input["imgAddr"].(string)
		if !assertOk {
			return NewServiceOutput(UserError, "InvalidImageAddress")
		}
	}
	imgAddr, err := valueObject.NewContainerImageAddress(imgAddrStr)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	portBindings := []valueObject.PortBinding{}
	if input["portBindings"] != nil {
		portBindings, assertOk = input["portBindings"].([]valueObject.PortBinding)
		if !assertOk {
			return NewServiceOutput(UserError, "InvalidPortBindings")
		}
	}

	var restartPolicyPtr *valueObject.ContainerRestartPolicy
	if input["restartPolicy"] != nil {
		restartPolicy, err := valueObject.NewContainerRestartPolicy(
			input["restartPolicy"].(string),
		)
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		restartPolicyPtr = &restartPolicy
	}

	var entrypointPtr *valueObject.ContainerEntrypoint
	if input["entrypoint"] != nil {
		entrypoint, err := valueObject.NewContainerEntrypoint(
			input["entrypoint"].(string),
		)
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		entrypointPtr = &entrypoint
	}

	var profileIdPtr *valueObject.ContainerProfileId
	if input["profileId"] != nil {
		profileId, err := valueObject.NewContainerProfileId(input["profileId"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		profileIdPtr = &profileId
	}

	envs := []valueObject.ContainerEnv{}
	if input["envs"] != nil {
		envs, assertOk = input["envs"].([]valueObject.ContainerEnv)
		if !assertOk {
			return NewServiceOutput(UserError, "InvalidEnvs")
		}
	}

	autoCreateMappings := true
	if input["autoCreateMappings"] != nil {
		var assertOk bool
		autoCreateMappings, assertOk = input["autoCreateMappings"].(bool)
		if !assertOk {
			var err error
			autoCreateMappings, err = strconv.ParseBool(
				input["autoCreateMappings"].(string),
			)
			if err != nil {
				return NewServiceOutput(UserError, err.Error())
			}
		}
	}

	addContainerDto := dto.NewAddContainer(
		accId,
		hostname,
		imgAddr,
		portBindings,
		restartPolicyPtr,
		entrypointPtr,
		profileIdPtr,
		envs,
		autoCreateMappings,
	)

	containerQueryRepo := infra.NewContainerQueryRepo(service.persistentDbSvc)
	containerCmdRepo := infra.NewContainerCmdRepo(service.persistentDbSvc)
	accQueryRepo := infra.NewAccQueryRepo(service.persistentDbSvc)
	accCmdRepo := infra.NewAccCmdRepo(service.persistentDbSvc)
	containerProfileQueryRepo := infra.NewContainerProfileQueryRepo(service.persistentDbSvc)
	mappingQueryRepo := infra.NewMappingQueryRepo(service.persistentDbSvc)
	mappingCmdRepo := infra.NewMappingCmdRepo(service.persistentDbSvc)

	err = useCase.AddContainer(
		containerQueryRepo,
		containerCmdRepo,
		accQueryRepo,
		accCmdRepo,
		containerProfileQueryRepo,
		mappingQueryRepo,
		mappingCmdRepo,
		addContainerDto,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Created, "ContainerCreated")
}
