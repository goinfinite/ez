package service

import (
	"strings"

	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/useCase"
	"github.com/speedianet/control/src/domain/valueObject"
	"github.com/speedianet/control/src/infra"
	"github.com/speedianet/control/src/infra/db"
	infraEnvs "github.com/speedianet/control/src/infra/envs"
	serviceHelper "github.com/speedianet/control/src/presentation/service/helper"
)

type ContainerImageService struct {
	persistentDbSvc         *db.PersistentDatabaseService
	containerImageQueryRepo *infra.ContainerImageQueryRepo
	activityRecordCmdRepo   *infra.ActivityRecordCmdRepo
}

func NewContainerImageService(
	persistentDbSvc *db.PersistentDatabaseService,
	trailDbSvc *db.TrailDatabaseService,
) *ContainerImageService {
	return &ContainerImageService{
		persistentDbSvc:         persistentDbSvc,
		containerImageQueryRepo: infra.NewContainerImageQueryRepo(persistentDbSvc),
		activityRecordCmdRepo:   infra.NewActivityRecordCmdRepo(trailDbSvc),
	}
}

func (service *ContainerImageService) Read() ServiceOutput {
	imagesList, err := useCase.ReadContainerImages(service.containerImageQueryRepo)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, imagesList)
}

func (service *ContainerImageService) CreateSnapshot(
	input map[string]interface{},
	shouldSchedule bool,
) ServiceOutput {
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

	if shouldSchedule {
		cliCmd := infraEnvs.SpeediaControlBinary + " container image create-snapshot"
		createParams := []string{
			"--account-id", accountId.String(),
			"--container-id", containerId.String(),
		}

		cliCmd = cliCmd + " " + strings.Join(createParams, " ")

		scheduledTaskCmdRepo := infra.NewScheduledTaskCmdRepo(service.persistentDbSvc)
		taskName, _ := valueObject.NewScheduledTaskName("CreateContainerSnapshotImage")
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

		return NewServiceOutput(Created, "ContainerSnapshotImageCreationScheduled")
	}

	operatorAccountId := LocalOperatorAccountId
	if input["operatorAccountId"] != nil {
		operatorAccountId, err = valueObject.NewAccountId(input["operatorAccountId"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	ipAddress := LocalOperatorIpAddress
	if input["ipAddress"] != nil {
		ipAddress, err = valueObject.NewIpAddress(input["ipAddress"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	createSnapshotImageDto := dto.NewCreateContainerSnapshotImage(
		accountId, containerId, operatorAccountId, ipAddress,
	)

	containerImageCmdRepo := infra.NewContainerImageCmdRepo(service.persistentDbSvc)

	err = useCase.CreateContainerSnapshotImage(
		containerImageCmdRepo, service.activityRecordCmdRepo, createSnapshotImageDto,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Created, "ContainerSnapshotImageCreated")
}

func (service *ContainerImageService) Export(
	input map[string]interface{},
) ServiceOutput {
	requiredParams := []string{"accountId", "imageId"}

	err := serviceHelper.RequiredParamsInspector(input, requiredParams)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	accountId, err := valueObject.NewAccountId(input["accountId"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	imageId, err := valueObject.NewContainerImageId(input["imageId"])
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

	ipAddress := LocalOperatorIpAddress
	if input["ipAddress"] != nil {
		ipAddress, err = valueObject.NewIpAddress(input["ipAddress"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	exportDto := dto.NewExportContainerImage(accountId, imageId, operatorAccountId, ipAddress)

	containerImageCmdRepo := infra.NewContainerImageCmdRepo(service.persistentDbSvc)

	downloadUrl, err := useCase.ExportContainerImage(
		service.containerImageQueryRepo, containerImageCmdRepo,
		service.activityRecordCmdRepo, exportDto,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, downloadUrl.String())
}

func (service *ContainerImageService) Delete(
	input map[string]interface{},
) ServiceOutput {
	requiredParams := []string{"accountId", "imageId"}
	err := serviceHelper.RequiredParamsInspector(input, requiredParams)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	accountId, err := valueObject.NewAccountId(input["accountId"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	imageId, err := valueObject.NewContainerImageId(input["imageId"])
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

	ipAddress := LocalOperatorIpAddress
	if input["ipAddress"] != nil {
		ipAddress, err = valueObject.NewIpAddress(input["ipAddress"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	deleteDto := dto.NewDeleteContainerImage(accountId, imageId, operatorAccountId, ipAddress)

	containerImageCmdRepo := infra.NewContainerImageCmdRepo(service.persistentDbSvc)

	err = useCase.DeleteContainerImage(
		service.containerImageQueryRepo, containerImageCmdRepo,
		service.activityRecordCmdRepo, deleteDto,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, "ContainerImageDeleted")
}
