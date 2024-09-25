package service

import (
	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/useCase"
	"github.com/speedianet/control/src/domain/valueObject"
	"github.com/speedianet/control/src/infra"
	"github.com/speedianet/control/src/infra/db"
	serviceHelper "github.com/speedianet/control/src/presentation/service/helper"
)

type MappingService struct {
	persistentDbSvc       *db.PersistentDatabaseService
	activityRecordCmdRepo *infra.ActivityRecordCmdRepo
}

func NewMappingService(
	persistentDbSvc *db.PersistentDatabaseService,
	trailDbSvc *db.TrailDatabaseService,
) *MappingService {
	return &MappingService{
		persistentDbSvc:       persistentDbSvc,
		activityRecordCmdRepo: infra.NewActivityRecordCmdRepo(trailDbSvc),
	}
}

func (service *MappingService) Read() ServiceOutput {
	mappingQueryRepo := infra.NewMappingQueryRepo(service.persistentDbSvc)
	mappingsList, err := useCase.ReadMappings(mappingQueryRepo)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, mappingsList)
}

func (service *MappingService) Create(input map[string]interface{}) ServiceOutput {
	requiredParams := []string{"accountId", "publicPort", "containerIds"}

	err := serviceHelper.RequiredParamsInspector(input, requiredParams)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	accountId, err := valueObject.NewAccountId(input["accountId"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	var hostnamePtr *valueObject.Fqdn
	if input["hostname"] != nil {
		hostname, err := valueObject.NewFqdn(input["hostname"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		hostnamePtr = &hostname
	}

	publicPort, err := valueObject.NewNetworkPort(input["publicPort"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	protocol := valueObject.GuessNetworkProtocolByPort(publicPort)
	if input["protocol"] != nil {
		protocol, err = valueObject.NewNetworkProtocol(input["protocol"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	containerIds, assertOk := input["containerIds"].([]valueObject.ContainerId)
	if !assertOk {
		return NewServiceOutput(UserError, "InvalidContainerIds")
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

	createMappingDto := dto.NewCreateMapping(
		accountId, hostnamePtr, publicPort, protocol, containerIds,
		operatorAccountId, operatorIpAddress,
	)

	mappingQueryRepo := infra.NewMappingQueryRepo(service.persistentDbSvc)
	mappingCmdRepo := infra.NewMappingCmdRepo(service.persistentDbSvc)
	containerQueryRepo := infra.NewContainerQueryRepo(service.persistentDbSvc)

	err = useCase.CreateMapping(
		mappingQueryRepo, mappingCmdRepo, containerQueryRepo,
		service.activityRecordCmdRepo, createMappingDto,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Created, "MappingCreated")
}

func (service *MappingService) Delete(input map[string]interface{}) ServiceOutput {
	requiredParams := []string{"mappingId"}

	err := serviceHelper.RequiredParamsInspector(input, requiredParams)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	accountId, err := valueObject.NewAccountId(input["accountId"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	mappingId, err := valueObject.NewMappingId(input["mappingId"])
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

	deleteDto := dto.NewDeleteMapping(
		accountId, mappingId, operatorAccountId, operatorIpAddress,
	)

	mappingQueryRepo := infra.NewMappingQueryRepo(service.persistentDbSvc)
	mappingCmdRepo := infra.NewMappingCmdRepo(service.persistentDbSvc)

	err = useCase.DeleteMapping(
		mappingQueryRepo, mappingCmdRepo, service.activityRecordCmdRepo,
		deleteDto,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, "MappingDeleted")
}

func (service *MappingService) CreateTarget(input map[string]interface{}) ServiceOutput {
	requiredParams := []string{"mappingId", "containerId"}

	err := serviceHelper.RequiredParamsInspector(input, requiredParams)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	accountId, err := valueObject.NewAccountId(input["accountId"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	mappingId, err := valueObject.NewMappingId(input["mappingId"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	containerId, err := valueObject.NewContainerId(input["containerId"])
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

	createTargetDto := dto.NewCreateMappingTarget(
		accountId, mappingId, containerId, operatorAccountId, operatorIpAddress,
	)

	mappingQueryRepo := infra.NewMappingQueryRepo(service.persistentDbSvc)
	mappingCmdRepo := infra.NewMappingCmdRepo(service.persistentDbSvc)
	containerQueryRepo := infra.NewContainerQueryRepo(service.persistentDbSvc)

	err = useCase.CreateMappingTarget(
		mappingQueryRepo, mappingCmdRepo, containerQueryRepo,
		service.activityRecordCmdRepo, createTargetDto,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Created, "MappingTargetCreated")
}

func (service *MappingService) DeleteTarget(input map[string]interface{}) ServiceOutput {
	requiredParams := []string{"mappingId", "targetId"}

	err := serviceHelper.RequiredParamsInspector(input, requiredParams)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	accountId, err := valueObject.NewAccountId(input["accountId"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	mappingId, err := valueObject.NewMappingId(input["mappingId"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	targetId, err := valueObject.NewMappingTargetId(input["targetId"])
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

	deleteDto := dto.NewDeleteMappingTarget(
		accountId, mappingId, targetId, operatorAccountId, operatorIpAddress,
	)

	mappingQueryRepo := infra.NewMappingQueryRepo(service.persistentDbSvc)
	mappingCmdRepo := infra.NewMappingCmdRepo(service.persistentDbSvc)

	err = useCase.DeleteMappingTarget(
		mappingQueryRepo, mappingCmdRepo, service.activityRecordCmdRepo,
		deleteDto,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, "MappingTargetDeleted")
}
