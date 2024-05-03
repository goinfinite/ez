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
	persistentDbSvc *db.PersistentDatabaseService
}

func NewMappingService(
	persistentDbSvc *db.PersistentDatabaseService,
) *MappingService {
	return &MappingService{
		persistentDbSvc: persistentDbSvc,
	}
}

func (service *MappingService) Read() ServiceOutput {
	mappingQueryRepo := infra.NewMappingQueryRepo(service.persistentDbSvc)
	mappingsList, err := useCase.GetMappings(mappingQueryRepo)
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
	if _, exists := input["hostname"]; exists {
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
	if _, exists := input["protocol"]; exists {
		protocol, err = valueObject.NewNetworkProtocol(input["protocol"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	var pathPtr *valueObject.MappingPath
	if _, exists := input["path"]; exists {
		path, err := valueObject.NewMappingPath(input["path"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		pathPtr = &path
	}

	var matchPatternPtr *valueObject.MappingMatchPattern
	if _, exists := input["matchPattern"]; exists {
		matchPattern, err := valueObject.NewMappingMatchPattern(input["matchPattern"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		matchPatternPtr = &matchPattern
	}

	containerIds, assertOk := input["containerIds"].([]valueObject.ContainerId)
	if !assertOk {
		return NewServiceOutput(UserError, "InvalidContainerIds")
	}

	createMappingDto := dto.NewCreateMapping(
		accountId,
		hostnamePtr,
		publicPort,
		protocol,
		pathPtr,
		matchPatternPtr,
		containerIds,
	)

	mappingQueryRepo := infra.NewMappingQueryRepo(service.persistentDbSvc)
	mappingCmdRepo := infra.NewMappingCmdRepo(service.persistentDbSvc)
	containerQueryRepo := infra.NewContainerQueryRepo(service.persistentDbSvc)

	err = useCase.CreateMapping(
		mappingQueryRepo,
		mappingCmdRepo,
		containerQueryRepo,
		createMappingDto,
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

	mappingId, err := valueObject.NewMappingId(input["mappingId"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	mappingQueryRepo := infra.NewMappingQueryRepo(service.persistentDbSvc)
	mappingCmdRepo := infra.NewMappingCmdRepo(service.persistentDbSvc)

	err = useCase.DeleteMapping(
		mappingQueryRepo,
		mappingCmdRepo,
		mappingId,
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

	mappingId, err := valueObject.NewMappingId(input["mappingId"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	containerId, err := valueObject.NewContainerId(input["containerId"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	createTargetDto := dto.NewCreateMappingTarget(
		mappingId,
		containerId,
	)

	mappingQueryRepo := infra.NewMappingQueryRepo(service.persistentDbSvc)
	mappingCmdRepo := infra.NewMappingCmdRepo(service.persistentDbSvc)
	containerQueryRepo := infra.NewContainerQueryRepo(service.persistentDbSvc)

	err = useCase.CreateMappingTarget(
		mappingQueryRepo,
		mappingCmdRepo,
		containerQueryRepo,
		createTargetDto,
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

	mappingId, err := valueObject.NewMappingId(input["mappingId"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	targetId, err := valueObject.NewMappingTargetId(input["targetId"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	mappingQueryRepo := infra.NewMappingQueryRepo(service.persistentDbSvc)
	mappingCmdRepo := infra.NewMappingCmdRepo(service.persistentDbSvc)

	err = useCase.DeleteMappingTarget(
		mappingQueryRepo,
		mappingCmdRepo,
		mappingId,
		targetId,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, "MappingTargetDeleted")
}
