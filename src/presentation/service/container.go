package service

import (
	"errors"
	"strings"

	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/useCase"
	"github.com/speedianet/control/src/domain/valueObject"
	voHelper "github.com/speedianet/control/src/domain/valueObject/helper"
	"github.com/speedianet/control/src/infra"
	"github.com/speedianet/control/src/infra/db"
	infraEnvs "github.com/speedianet/control/src/infra/envs"
	infraHelper "github.com/speedianet/control/src/infra/helper"
	serviceHelper "github.com/speedianet/control/src/presentation/service/helper"
)

type ContainerService struct {
	persistentDbSvc    *db.PersistentDatabaseService
	containerQueryRepo *infra.ContainerQueryRepo
}

func NewContainerService(
	persistentDbSvc *db.PersistentDatabaseService,
) *ContainerService {
	return &ContainerService{
		persistentDbSvc:    persistentDbSvc,
		containerQueryRepo: infra.NewContainerQueryRepo(persistentDbSvc),
	}
}

func (service *ContainerService) Read() ServiceOutput {
	containersList, err := useCase.ReadContainers(service.containerQueryRepo)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, containersList)
}

func (service *ContainerService) ReadWithMetrics() ServiceOutput {
	containersList, err := useCase.ReadContainersWithMetrics(service.containerQueryRepo)
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

	containerCmdRepo := infra.NewContainerCmdRepo(service.persistentDbSvc)
	accessToken, err := useCase.ContainerAutoLogin(
		service.containerQueryRepo, containerCmdRepo, autoLoginDto,
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
			return NewServiceOutput(UserError, err.Error())
		}
	}

	if shouldSchedule {
		cliCmd := infraEnvs.SpeediaControlBinary + " container create"
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

	containerCmdRepo := infra.NewContainerCmdRepo(service.persistentDbSvc)
	accountQueryRepo := infra.NewAccountQueryRepo(service.persistentDbSvc)
	accountCmdRepo := infra.NewAccountCmdRepo(service.persistentDbSvc)
	containerProfileQueryRepo := infra.NewContainerProfileQueryRepo(service.persistentDbSvc)
	mappingQueryRepo := infra.NewMappingQueryRepo(service.persistentDbSvc)
	mappingCmdRepo := infra.NewMappingCmdRepo(service.persistentDbSvc)
	containerProxyCmdRepo := infra.NewContainerProxyCmdRepo(service.persistentDbSvc)

	err = useCase.CreateContainer(
		service.containerQueryRepo, containerCmdRepo, accountQueryRepo, accountCmdRepo,
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

	updateContainerDto := dto.NewUpdateContainer(
		accountId,
		containerId,
		containerStatusPtr,
		profileIdPtr,
	)

	containerCmdRepo := infra.NewContainerCmdRepo(service.persistentDbSvc)
	accountQueryRepo := infra.NewAccountQueryRepo(service.persistentDbSvc)
	accountCmdRepo := infra.NewAccountCmdRepo(service.persistentDbSvc)
	containerProfileQueryRepo := infra.NewContainerProfileQueryRepo(service.persistentDbSvc)

	err = useCase.UpdateContainer(
		service.containerQueryRepo, containerCmdRepo, accountQueryRepo, accountCmdRepo,
		containerProfileQueryRepo, updateContainerDto,
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

	containerCmdRepo := infra.NewContainerCmdRepo(service.persistentDbSvc)
	accountCmdRepo := infra.NewAccountCmdRepo(service.persistentDbSvc)
	mappingQueryRepo := infra.NewMappingQueryRepo(service.persistentDbSvc)
	mappingCmdRepo := infra.NewMappingCmdRepo(service.persistentDbSvc)
	containerProxyCmdRepo := infra.NewContainerProxyCmdRepo(service.persistentDbSvc)

	err = useCase.DeleteContainer(
		service.containerQueryRepo, containerCmdRepo, accountCmdRepo, mappingQueryRepo,
		mappingCmdRepo, containerProxyCmdRepo, accountId, containerId,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, "ContainerDeleted")
}
