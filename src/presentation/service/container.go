package service

import (
	"errors"
	"strings"

	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/useCase"
	"github.com/goinfinite/ez/src/domain/valueObject"
	voHelper "github.com/goinfinite/ez/src/domain/valueObject/helper"
	"github.com/goinfinite/ez/src/infra"
	"github.com/goinfinite/ez/src/infra/db"
	infraEnvs "github.com/goinfinite/ez/src/infra/envs"
	infraHelper "github.com/goinfinite/ez/src/infra/helper"
	serviceHelper "github.com/goinfinite/ez/src/presentation/service/helper"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type ContainerService struct {
	persistentDbSvc       *db.PersistentDatabaseService
	containerQueryRepo    *infra.ContainerQueryRepo
	activityRecordCmdRepo *infra.ActivityRecordCmdRepo
}

func NewContainerService(
	persistentDbSvc *db.PersistentDatabaseService,
	trailDbSvc *db.TrailDatabaseService,
) *ContainerService {
	return &ContainerService{
		persistentDbSvc:       persistentDbSvc,
		containerQueryRepo:    infra.NewContainerQueryRepo(persistentDbSvc),
		activityRecordCmdRepo: infra.NewActivityRecordCmdRepo(trailDbSvc),
	}
}

func (service *ContainerService) Read(input map[string]interface{}) ServiceOutput {
	var containerId []valueObject.ContainerId
	var assertOk bool
	if input["containerId"] != nil {
		containerId, assertOk = input["containerId"].([]valueObject.ContainerId)
		if !assertOk {
			return NewServiceOutput(UserError, "InvalidContainerId")
		}
	}

	var containerAccountId []valueObject.AccountId
	if input["containerAccountId"] != nil {
		containerAccountId, assertOk = input["containerAccountId"].([]valueObject.AccountId)
		if !assertOk {
			return NewServiceOutput(UserError, "InvalidContainerAccountId")
		}
	}

	var containerHostnamePtr *valueObject.Fqdn
	if input["containerHostname"] != nil {
		containerHostname, err := valueObject.NewFqdn(input["containerHostname"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		containerHostnamePtr = &containerHostname
	}

	var containerStatusPtr *bool
	if input["containerStatus"] != nil {
		containerStatus, err := voHelper.InterfaceToBool(input["containerStatus"])
		if err != nil {
			return NewServiceOutput(UserError, "InvalidContainerStatus")
		}
		containerStatusPtr = &containerStatus
	}

	var containerImageIdPtr *valueObject.ContainerImageId
	if input["containerImageId"] != nil {
		containerImageId, err := valueObject.NewContainerImageId(input["containerImageId"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		containerImageIdPtr = &containerImageId
	}

	var containerImageAddressPtr *valueObject.ContainerImageAddress
	if input["containerImageAddress"] != nil {
		containerImageAddress, err := valueObject.NewContainerImageAddress(
			input["containerImageAddress"],
		)
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		containerImageAddressPtr = &containerImageAddress
	}

	var containerImageHashPtr *valueObject.Hash
	if input["containerImageHash"] != nil {
		containerImageHash, err := valueObject.NewHash(input["containerImageHash"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		containerImageHashPtr = &containerImageHash
	}

	containerPortBindings := []valueObject.PortBinding{}
	if input["containerPortBindings"] != nil {
		containerPortBindings, assertOk = input["containerPortBindings"].([]valueObject.PortBinding)
		if !assertOk {
			return NewServiceOutput(UserError, "InvalidContainerPortBindings")
		}
	}

	var containerRestartPolicyPtr *valueObject.ContainerRestartPolicy
	if input["containerRestartPolicy"] != nil {
		containerRestartPolicy, err := valueObject.NewContainerRestartPolicy(
			input["containerRestartPolicy"],
		)
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		containerRestartPolicyPtr = &containerRestartPolicy
	}

	var containerProfileIdPtr *valueObject.ContainerProfileId
	if input["containerProfileId"] != nil {
		containerProfileId, err := valueObject.NewContainerProfileId(input["containerProfileId"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		containerProfileIdPtr = &containerProfileId
	}

	containerEnv := []valueObject.ContainerEnv{}
	if input["containerEnv"] != nil {
		var assertOk bool
		containerEnv, assertOk = input["containerEnv"].([]valueObject.ContainerEnv)
		if !assertOk {
			return NewServiceOutput(UserError, "InvalidContainerEnv")
		}
	}

	var createdBeforeAtPtr, createdAfterAtPtr *valueObject.UnixTime
	var startedBeforeAtPtr, startedAfterAtPtr *valueObject.UnixTime
	var stoppedBeforeAtPtr, stoppedAfterAtPtr *valueObject.UnixTime

	timeParamNames := []string{
		"createdBeforeAt", "createdAfterAt",
		"startedBeforeAt", "startedAfterAt",
		"stoppedBeforeAt", "stoppedAfterAt",
	}
	for _, timeParamName := range timeParamNames {
		if input[timeParamName] == nil {
			continue
		}

		timeParam, err := valueObject.NewUnixTime(input[timeParamName])
		if err != nil {
			capitalParamName := cases.Title(language.English).String(timeParamName)
			return NewServiceOutput(UserError, "Invalid"+capitalParamName)
		}

		switch timeParamName {
		case "createdBeforeAt":
			createdBeforeAtPtr = &timeParam
		case "createdAfterAt":
			createdAfterAtPtr = &timeParam
		case "startedBeforeAt":
			startedBeforeAtPtr = &timeParam
		case "startedAfterAt":
			startedAfterAtPtr = &timeParam
		case "stoppedBeforeAt":
			stoppedBeforeAtPtr = &timeParam
		case "stoppedAfterAt":
			stoppedAfterAtPtr = &timeParam
		}
	}

	var withMetricsPtr *bool
	if input["withMetrics"] != nil {
		withMetrics, err := voHelper.InterfaceToBool(input["withMetrics"])
		if err != nil {
			return NewServiceOutput(UserError, "InvalidWithMetrics")
		}
		withMetricsPtr = &withMetrics
	}

	paginationDto := useCase.ContainersDefaultPagination
	if input["pageNumber"] != nil {
		pageNumber, err := voHelper.InterfaceToUint32(input["pageNumber"])
		if err != nil {
			return NewServiceOutput(UserError, "InvalidPageNumber")
		}
		paginationDto.PageNumber = pageNumber
	}

	if input["itemsPerPage"] != nil {
		itemsPerPage, err := voHelper.InterfaceToUint16(input["itemsPerPage"])
		if err != nil {
			return NewServiceOutput(UserError, "InvalidItemsPerPage")
		}
		paginationDto.ItemsPerPage = itemsPerPage
	}

	if input["sortBy"] != nil {
		sortBy, err := valueObject.NewPaginationSortBy(input["sortBy"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		paginationDto.SortBy = &sortBy
	}

	if input["sortDirection"] != nil {
		sortDirection, err := valueObject.NewPaginationSortDirection(input["sortDirection"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		paginationDto.SortDirection = &sortDirection
	}

	if input["lastSeenId"] != nil {
		lastSeenId, err := valueObject.NewPaginationLastSeenId(input["lastSeenId"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		paginationDto.LastSeenId = &lastSeenId
	}

	readDto := dto.ReadContainersRequest{
		Pagination:             paginationDto,
		ContainerId:            containerId,
		ContainerAccountId:     containerAccountId,
		ContainerHostname:      containerHostnamePtr,
		ContainerStatus:        containerStatusPtr,
		ContainerImageId:       containerImageIdPtr,
		ContainerImageAddress:  containerImageAddressPtr,
		ContainerImageHash:     containerImageHashPtr,
		ContainerPortBindings:  containerPortBindings,
		ContainerRestartPolicy: containerRestartPolicyPtr,
		ContainerProfileId:     containerProfileIdPtr,
		ContainerEnv:           containerEnv,
		CreatedBeforeAt:        createdBeforeAtPtr,
		CreatedAfterAt:         createdAfterAtPtr,
		StartedBeforeAt:        startedBeforeAtPtr,
		StartedAfterAt:         startedAfterAtPtr,
		StoppedBeforeAt:        stoppedBeforeAtPtr,
		StoppedAfterAt:         stoppedAfterAtPtr,
		WithMetrics:            withMetricsPtr,
	}

	responseDto, err := useCase.ReadContainers(service.containerQueryRepo, readDto)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, responseDto)
}

func (service *ContainerService) CreateContainerSessionToken(
	input map[string]interface{},
) ServiceOutput {
	requiredParams := []string{"accountId", "containerId", "sessionIpAddress"}

	err := serviceHelper.RequiredParamsInspector(input, requiredParams)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	accountId, err := valueObject.NewAccountId(input["accountId"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	containerId, err := valueObject.NewContainerId(input["containerId"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	sessionIpAddress, err := valueObject.NewIpAddress(input["sessionIpAddress"])
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

	createDto := dto.NewCreateContainerSessionToken(
		accountId, containerId, sessionIpAddress,
		operatorAccountId, operatorIpAddress,
	)

	containerCmdRepo := infra.NewContainerCmdRepo(service.persistentDbSvc)
	accessToken, err := useCase.CreateContainerSessionToken(
		service.containerQueryRepo, containerCmdRepo,
		service.activityRecordCmdRepo, createDto,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, accessToken)

}

func (service *ContainerService) Create(
	input map[string]interface{},
	shouldSchedule bool,
) ServiceOutput {
	requiredParams := []string{"accountId", "hostname"}

	err := serviceHelper.RequiredParamsInspector(input, requiredParams)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	accountId, err := valueObject.NewAccountId(input["accountId"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	hostname, err := valueObject.NewFqdn(input["hostname"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	if _, exists := input["imageAddress"]; !exists {
		input["imageAddress"] = "goinfinite/os"
	}
	imgAddr, err := valueObject.NewContainerImageAddress(input["imageAddress"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	var imageIdPtr *valueObject.ContainerImageId
	if input["imageId"] != nil {
		imageId, err := valueObject.NewContainerImageId(input["imageId"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		imageIdPtr = &imageId
	}

	portBindings := []valueObject.PortBinding{}
	if input["portBindings"] != nil {
		var assertOk bool
		portBindings, assertOk = input["portBindings"].([]valueObject.PortBinding)
		if !assertOk {
			return NewServiceOutput(UserError, "InvalidPortBindings")
		}
	}

	var restartPolicyPtr *valueObject.ContainerRestartPolicy
	if input["restartPolicy"] != nil {
		restartPolicy, err := valueObject.NewContainerRestartPolicy(input["restartPolicy"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		restartPolicyPtr = &restartPolicy
	}

	var entrypointPtr *valueObject.ContainerEntrypoint
	if input["entrypoint"] != nil {
		entrypoint, err := valueObject.NewContainerEntrypoint(input["entrypoint"])
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
		var assertOk bool
		envs, assertOk = input["envs"].([]valueObject.ContainerEnv)
		if !assertOk {
			return NewServiceOutput(UserError, "InvalidEnvs")
		}
	}

	var launchScriptPtr *valueObject.LaunchScript
	if input["launchScript"] != nil {
		launchScript, assertOk := input["launchScript"].(valueObject.LaunchScript)
		if !assertOk {
			return NewServiceOutput(UserError, "InvalidLaunchScript")
		}
		launchScriptPtr = &launchScript
	}

	autoCreateMappings := true
	if input["autoCreateMappings"] != nil {
		autoCreateMappings, err = voHelper.InterfaceToBool(input["autoCreateMappings"])
		if err != nil {
			return NewServiceOutput(UserError, "InvalidAutoCreateMappings")
		}
	}

	useImageExposedPorts := false
	if input["useImageExposedPorts"] != nil {
		useImageExposedPorts, err = voHelper.InterfaceToBool(input["useImageExposedPorts"])
		if err != nil {
			return NewServiceOutput(UserError, "InvalidUseImageExposedPorts")
		}
	}

	var existingContainerIdPtr *valueObject.ContainerId
	if input["existingContainerId"] != nil {
		existingContainerId, err := valueObject.NewContainerId(input["existingContainerId"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		existingContainerIdPtr = &existingContainerId
	}

	if shouldSchedule {
		cliCmd := infraEnvs.InfiniteEzBinary + " container create"
		createParams := []string{
			"--account-id", accountId.String(),
			"--hostname", hostname.String(),
		}
		if imageIdPtr == nil {
			createParams = append(createParams, "--image-address")
			createParams = append(createParams, imgAddr.String())
		} else {
			createParams = append(createParams, "--image-id")
			createParams = append(createParams, imageIdPtr.String())
		}

		if len(portBindings) > 0 {
			for _, portBinding := range portBindings {
				createParams = append(createParams, "--port-bindings")
				createParams = append(createParams, portBinding.String())
			}
		}
		if restartPolicyPtr != nil {
			createParams = append(createParams, "--restart-policy")
			createParams = append(createParams, restartPolicyPtr.String())
		}
		if entrypointPtr != nil {
			createParams = append(createParams, "--entrypoint")
			createParams = append(createParams, entrypointPtr.String())
		}
		if profileIdPtr != nil {
			createParams = append(createParams, "--profile-id")
			createParams = append(createParams, profileIdPtr.String())
		}
		if len(envs) > 0 {
			for _, env := range envs {
				createParams = append(createParams, "--envs")
				createParams = append(createParams, env.String())
			}
		}

		if launchScriptPtr != nil {
			launchScriptTempFilePath := "/tmp/ls-" + hostname.String() + ".sh"
			err = infraHelper.UpdateFile(
				launchScriptTempFilePath, launchScriptPtr.String(), true,
			)
			if err != nil {
				return NewServiceOutput(
					InfraError, errors.New("SaveLaunchScriptError").Error(),
				)
			}
			createParams = append(createParams, "--launch-script-path")
			createParams = append(createParams, launchScriptTempFilePath)
		}

		if !autoCreateMappings {
			createParams = append(createParams, "--auto-create-mappings")
			createParams = append(createParams, "false")
		}

		if useImageExposedPorts {
			createParams = append(createParams, "--use-image-exposed-ports")
			createParams = append(createParams, "true")
		}

		if existingContainerIdPtr != nil {
			createParams = append(createParams, "--existing-container-id")
			createParams = append(createParams, existingContainerIdPtr.String())
		}

		cliCmd = cliCmd + " " + strings.Join(createParams, " ")

		scheduledTaskCmdRepo := infra.NewScheduledTaskCmdRepo(service.persistentDbSvc)
		taskName, _ := valueObject.NewScheduledTaskName("CreateContainer")
		taskCmd, _ := valueObject.NewUnixCommand(cliCmd)
		taskTag, _ := valueObject.NewScheduledTaskTag("container")
		taskTags := []valueObject.ScheduledTaskTag{taskTag}
		taskTimeoutSecs := valueObject.TimeDuration(uint64(1800))

		scheduledTaskCreateDto := dto.NewCreateScheduledTask(
			taskName, taskCmd, taskTags, &taskTimeoutSecs, nil,
		)

		err = useCase.CreateScheduledTask(scheduledTaskCmdRepo, scheduledTaskCreateDto)
		if err != nil {
			return NewServiceOutput(InfraError, err.Error())
		}

		return NewServiceOutput(Created, "ContainerCreationScheduled")
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

	createContainerDto := dto.NewCreateContainer(
		accountId, hostname, imgAddr, imageIdPtr, portBindings, restartPolicyPtr, entrypointPtr,
		profileIdPtr, envs, launchScriptPtr, autoCreateMappings, useImageExposedPorts,
		existingContainerIdPtr, operatorAccountId, operatorIpAddress,
	)

	containerCmdRepo := infra.NewContainerCmdRepo(service.persistentDbSvc)
	containerImageQueryRepo := infra.NewContainerImageQueryRepo(service.persistentDbSvc)
	containerImageCmdRepo := infra.NewContainerImageCmdRepo(service.persistentDbSvc)
	accountQueryRepo := infra.NewAccountQueryRepo(service.persistentDbSvc)
	accountCmdRepo := infra.NewAccountCmdRepo(service.persistentDbSvc)
	containerProfileQueryRepo := infra.NewContainerProfileQueryRepo(service.persistentDbSvc)
	mappingQueryRepo := infra.NewMappingQueryRepo(service.persistentDbSvc)
	mappingCmdRepo := infra.NewMappingCmdRepo(service.persistentDbSvc)
	containerProxyCmdRepo := infra.NewContainerProxyCmdRepo(service.persistentDbSvc)

	containerId, err := useCase.CreateContainer(
		service.containerQueryRepo, containerCmdRepo, containerImageQueryRepo,
		containerImageCmdRepo, accountQueryRepo, accountCmdRepo,
		containerProfileQueryRepo, mappingQueryRepo, mappingCmdRepo,
		containerProxyCmdRepo, service.activityRecordCmdRepo, createContainerDto,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Created, containerId)
}

func (service *ContainerService) Update(input map[string]interface{}) ServiceOutput {
	requiredParams := []string{"accountId", "containerId"}

	err := serviceHelper.RequiredParamsInspector(input, requiredParams)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	accountId, err := valueObject.NewAccountId(input["accountId"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	containerId, err := valueObject.NewContainerId(input["containerId"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	var containerStatusPtr *bool
	if input["status"] != nil {
		containerStatus, err := voHelper.InterfaceToBool(input["status"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		containerStatusPtr = &containerStatus
	}

	var profileIdPtr *valueObject.ContainerProfileId
	if input["profileId"] != nil {
		profileId, err := valueObject.NewContainerProfileId(input["profileId"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		profileIdPtr = &profileId
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

	updateContainerDto := dto.NewUpdateContainer(
		accountId, containerId, containerStatusPtr, profileIdPtr,
		operatorAccountId, operatorIpAddress,
	)

	containerCmdRepo := infra.NewContainerCmdRepo(service.persistentDbSvc)
	accountQueryRepo := infra.NewAccountQueryRepo(service.persistentDbSvc)
	accountCmdRepo := infra.NewAccountCmdRepo(service.persistentDbSvc)
	containerProfileQueryRepo := infra.NewContainerProfileQueryRepo(service.persistentDbSvc)

	err = useCase.UpdateContainer(
		service.containerQueryRepo, containerCmdRepo, accountQueryRepo, accountCmdRepo,
		containerProfileQueryRepo, service.activityRecordCmdRepo, updateContainerDto,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, "ContainerUpdated")
}

func (service *ContainerService) Delete(input map[string]interface{}) ServiceOutput {
	accountId, err := valueObject.NewAccountId(input["accountId"])
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

	deleteDto := dto.NewDeleteContainer(
		accountId, containerId, operatorAccountId, operatorIpAddress,
	)

	containerCmdRepo := infra.NewContainerCmdRepo(service.persistentDbSvc)
	accountCmdRepo := infra.NewAccountCmdRepo(service.persistentDbSvc)
	mappingCmdRepo := infra.NewMappingCmdRepo(service.persistentDbSvc)
	containerProxyCmdRepo := infra.NewContainerProxyCmdRepo(service.persistentDbSvc)

	err = useCase.DeleteContainer(
		service.containerQueryRepo, containerCmdRepo, accountCmdRepo,
		mappingCmdRepo, containerProxyCmdRepo, service.activityRecordCmdRepo,
		deleteDto,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, "ContainerDeleted")
}
