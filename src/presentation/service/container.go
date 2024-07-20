package service

import (
	"errors"
	"strings"

	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/useCase"
	"github.com/speedianet/control/src/domain/valueObject"
	"github.com/speedianet/control/src/infra"
	"github.com/speedianet/control/src/infra/db"
	infraHelper "github.com/speedianet/control/src/infra/helper"
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

func (service *ContainerService) Read() ServiceOutput {
	containerQueryRepo := infra.NewContainerQueryRepo(service.persistentDbSvc)
	containersList, err := useCase.ReadContainers(containerQueryRepo)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, containersList)
}

func (service *ContainerService) ReadWithMetrics() ServiceOutput {
	containerQueryRepo := infra.NewContainerQueryRepo(service.persistentDbSvc)
	containersList, err := useCase.ReadContainersWithMetrics(containerQueryRepo)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, containersList)
}

func (service *ContainerService) AutoLogin(input map[string]interface{}) ServiceOutput {
	requiredParams := []string{"containerId", "ipAddress"}

	err := serviceHelper.RequiredParamsInspector(input, requiredParams)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	containerId, err := valueObject.NewContainerId(input["containerId"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	ipAddress, err := valueObject.NewIpAddress(input["ipAddress"].(string))
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	autoLoginDto := dto.NewContainerAutoLogin(containerId, ipAddress)

	containerQueryRepo := infra.NewContainerQueryRepo(service.persistentDbSvc)
	containerCmdRepo := infra.NewContainerCmdRepo(service.persistentDbSvc)
	accessToken, err := useCase.ContainerAutoLogin(
		containerQueryRepo, containerCmdRepo, autoLoginDto,
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
		input["imageAddress"] = "speedianet/os"
	}
	imgAddr, err := valueObject.NewContainerImageAddress(input["imageAddress"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	portBindings := []valueObject.PortBinding{}
	if _, exists := input["portBindings"]; exists {
		var assertOk bool
		portBindings, assertOk = input["portBindings"].([]valueObject.PortBinding)
		if !assertOk {
			return NewServiceOutput(UserError, "InvalidPortBindings")
		}
	}

	var restartPolicyPtr *valueObject.ContainerRestartPolicy
	if _, exists := input["restartPolicy"]; exists {
		restartPolicy, err := valueObject.NewContainerRestartPolicy(input["restartPolicy"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		restartPolicyPtr = &restartPolicy
	}

	var entrypointPtr *valueObject.ContainerEntrypoint
	if _, exists := input["entrypoint"]; exists {
		entrypoint, err := valueObject.NewContainerEntrypoint(input["entrypoint"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		entrypointPtr = &entrypoint
	}

	var profileIdPtr *valueObject.ContainerProfileId
	if _, exists := input["profileId"]; exists {
		profileId, err := valueObject.NewContainerProfileId(input["profileId"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		profileIdPtr = &profileId
	}

	envs := []valueObject.ContainerEnv{}
	if _, exists := input["envs"]; exists {
		var assertOk bool
		envs, assertOk = input["envs"].([]valueObject.ContainerEnv)
		if !assertOk {
			return NewServiceOutput(UserError, "InvalidEnvs")
		}
	}

	var launchScriptPtr *valueObject.LaunchScript
	if _, exists := input["launchScript"]; exists {
		launchScript, assertOk := input["launchScript"].(valueObject.LaunchScript)
		if !assertOk {
			return NewServiceOutput(UserError, "InvalidLaunchScript")
		}
		launchScriptPtr = &launchScript
	}

	autoCreateMappings := true
	if _, exists := input["autoCreateMappings"]; exists {
		autoCreateMappings, err = serviceHelper.ParseBoolParam(input["autoCreateMappings"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	if shouldSchedule {
		cliCmd := "/var/speedia/control container create"
		createParams := []string{
			"--account-id", accountId.String(),
			"--hostname", hostname.String(),
			"--image-address", imgAddr.String(),
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

		cliCmd = cliCmd + " " + strings.Join(createParams, " ")

		scheduledTaskCmdRepo := infra.NewScheduledTaskCmdRepo(service.persistentDbSvc)
		taskName, _ := valueObject.NewScheduledTaskName("CreateContainer")
		taskCmd, _ := valueObject.NewUnixCommand(cliCmd)
		taskTag, _ := valueObject.NewScheduledTaskTag("container")
		taskTags := []valueObject.ScheduledTaskTag{taskTag}
		timeoutSeconds := uint(900)

		scheduledTaskCreateDto := dto.NewCreateScheduledTask(
			taskName, taskCmd, taskTags, &timeoutSeconds, nil,
		)

		err = useCase.CreateScheduledTask(scheduledTaskCmdRepo, scheduledTaskCreateDto)
		if err != nil {
			return NewServiceOutput(InfraError, err.Error())
		}

		return NewServiceOutput(Created, "ContainerCreationScheduled")
	}

	createContainerDto := dto.NewCreateContainer(
		accountId, hostname, imgAddr, portBindings, restartPolicyPtr, entrypointPtr,
		profileIdPtr, envs, launchScriptPtr, autoCreateMappings,
	)

	containerQueryRepo := infra.NewContainerQueryRepo(service.persistentDbSvc)
	containerCmdRepo := infra.NewContainerCmdRepo(service.persistentDbSvc)
	accountQueryRepo := infra.NewAccountQueryRepo(service.persistentDbSvc)
	accountCmdRepo := infra.NewAccountCmdRepo(service.persistentDbSvc)
	containerProfileQueryRepo := infra.NewContainerProfileQueryRepo(service.persistentDbSvc)
	mappingQueryRepo := infra.NewMappingQueryRepo(service.persistentDbSvc)
	mappingCmdRepo := infra.NewMappingCmdRepo(service.persistentDbSvc)
	containerProxyCmdRepo := infra.NewContainerProxyCmdRepo(service.persistentDbSvc)

	err = useCase.CreateContainer(
		containerQueryRepo, containerCmdRepo, accountQueryRepo, accountCmdRepo,
		containerProfileQueryRepo, mappingQueryRepo, mappingCmdRepo,
		containerProxyCmdRepo, createContainerDto,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Created, "ContainerCreated")
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
	if _, exists := input["status"]; exists {
		containerStatus, err := serviceHelper.ParseBoolParam(input["status"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		containerStatusPtr = &containerStatus
	}

	var profileIdPtr *valueObject.ContainerProfileId
	if _, exists := input["profileId"]; exists {
		profileId, err := valueObject.NewContainerProfileId(input["profileId"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		profileIdPtr = &profileId
	}

	updateContainerDto := dto.NewUpdateContainer(
		accountId,
		containerId,
		containerStatusPtr,
		profileIdPtr,
	)

	containerQueryRepo := infra.NewContainerQueryRepo(service.persistentDbSvc)
	containerCmdRepo := infra.NewContainerCmdRepo(service.persistentDbSvc)
	accountQueryRepo := infra.NewAccountQueryRepo(service.persistentDbSvc)
	accountCmdRepo := infra.NewAccountCmdRepo(service.persistentDbSvc)
	containerProfileQueryRepo := infra.NewContainerProfileQueryRepo(service.persistentDbSvc)

	err = useCase.UpdateContainer(
		containerQueryRepo,
		containerCmdRepo,
		accountQueryRepo,
		accountCmdRepo,
		containerProfileQueryRepo,
		updateContainerDto,
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

	containerQueryRepo := infra.NewContainerQueryRepo(service.persistentDbSvc)
	containerCmdRepo := infra.NewContainerCmdRepo(service.persistentDbSvc)
	accountCmdRepo := infra.NewAccountCmdRepo(service.persistentDbSvc)
	mappingQueryRepo := infra.NewMappingQueryRepo(service.persistentDbSvc)
	mappingCmdRepo := infra.NewMappingCmdRepo(service.persistentDbSvc)
	containerProxyCmdRepo := infra.NewContainerProxyCmdRepo(service.persistentDbSvc)

	err = useCase.DeleteContainer(
		containerQueryRepo,
		containerCmdRepo,
		accountCmdRepo,
		mappingQueryRepo,
		mappingCmdRepo,
		containerProxyCmdRepo,
		accountId,
		containerId,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, "ContainerDeleted")
}
