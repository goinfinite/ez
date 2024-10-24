package service

import (
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
	requiredParams := []string{"containerId"}

	err := serviceHelper.RequiredParamsInspector(input, requiredParams)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	containerId, err := valueObject.NewContainerId(input["containerId"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	var shouldCreateArchivePtr *bool
	if input["shouldCreateArchive"] != nil {
		shouldCreateArchive, err := voHelper.InterfaceToBool(input["shouldCreateArchive"])
		if err != nil {
			return NewServiceOutput(UserError, "InvalidShouldCreateArchive")
		}
		shouldCreateArchivePtr = &shouldCreateArchive
	}

	var archiveCompressionFormatPtr *valueObject.CompressionFormat
	if input["archiveCompressionFormat"] != nil {
		compressionFormat, err := valueObject.NewCompressionFormat(
			input["archiveCompressionFormat"],
		)
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		archiveCompressionFormatPtr = &compressionFormat
	}

	var shouldDiscardImagePtr *bool
	if input["shouldDiscardImage"] != nil {
		shouldDiscardImage, err := voHelper.InterfaceToBool(input["shouldDiscardImage"])
		if err != nil {
			return NewServiceOutput(UserError, "InvalidShouldDiscardImage")
		}
		shouldDiscardImagePtr = &shouldDiscardImage
	}

	if shouldSchedule {
		cliCmd := infraEnvs.InfiniteEzBinary + " container image create-snapshot"
		createParams := []string{
			"--container-id", containerId.String(),
		}
		if shouldCreateArchivePtr != nil && *shouldCreateArchivePtr {
			createParams = append(createParams, "--should-create-archive", "true")
		}

		if archiveCompressionFormatPtr != nil {
			createParams = append(
				createParams, "--archive-compression-format", archiveCompressionFormatPtr.String(),
			)
		}

		if shouldDiscardImagePtr != nil && *shouldDiscardImagePtr {
			createParams = append(createParams, "--should-discard-image", "true")
		}

		cliCmd = cliCmd + " " + strings.Join(createParams, " ")

		scheduledTaskCmdRepo := infra.NewScheduledTaskCmdRepo(service.persistentDbSvc)
		taskName, _ := valueObject.NewScheduledTaskName("CreateContainerSnapshotImage")
		taskCmd, _ := valueObject.NewUnixCommand(cliCmd)
		taskTag, _ := valueObject.NewScheduledTaskTag("containerImage")
		taskTags := []valueObject.ScheduledTaskTag{taskTag}
		timeoutSeconds := uint16(1800)

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

	operatorIpAddress := LocalOperatorIpAddress
	if input["operatorIpAddress"] != nil {
		operatorIpAddress, err = valueObject.NewIpAddress(input["operatorIpAddress"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	createSnapshotImageDto := dto.NewCreateContainerSnapshotImage(
		containerId, shouldCreateArchivePtr, archiveCompressionFormatPtr,
		shouldDiscardImagePtr, operatorAccountId, operatorIpAddress,
	)

	containerImageCmdRepo := infra.NewContainerImageCmdRepo(service.persistentDbSvc)
	containerQueryRepo := infra.NewContainerQueryRepo(service.persistentDbSvc)
	accountQueryRepo := infra.NewAccountQueryRepo(service.persistentDbSvc)

	err = useCase.CreateContainerSnapshotImage(
		service.containerImageQueryRepo, containerImageCmdRepo, containerQueryRepo,
		accountQueryRepo, service.activityRecordCmdRepo, createSnapshotImageDto,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Created, "ContainerSnapshotImageCreated")
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

	operatorIpAddress := LocalOperatorIpAddress
	if input["operatorIpAddress"] != nil {
		operatorIpAddress, err = valueObject.NewIpAddress(input["operatorIpAddress"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	deleteDto := dto.NewDeleteContainerImage(
		accountId, imageId, operatorAccountId, operatorIpAddress,
	)

	containerImageCmdRepo := infra.NewContainerImageCmdRepo(service.persistentDbSvc)
	containerQueryRepo := infra.NewContainerQueryRepo(service.persistentDbSvc)

	err = useCase.DeleteContainerImage(
		service.containerImageQueryRepo, containerImageCmdRepo, containerQueryRepo,
		service.activityRecordCmdRepo, deleteDto,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, "ContainerImageDeleted")
}

func (service *ContainerImageService) ReadArchiveFiles(
	requestHostname *string,
) ServiceOutput {
	archiveFilesList, err := useCase.ReadContainerImageArchiveFiles(service.containerImageQueryRepo)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	if requestHostname != nil {
		serverHostname, err := infraHelper.ReadServerHostname()
		if err != nil {
			return NewServiceOutput(InfraError, err.Error())
		}
		serverHostnameStr := serverHostname.String()

		for archiveFileIndex, archiveFile := range archiveFilesList {
			rawUpdatedUrl := strings.Replace(
				archiveFile.DownloadUrl.String(), serverHostnameStr, *requestHostname, 1,
			)

			updatedUrl, err := valueObject.NewUrl(rawUpdatedUrl)
			if err != nil {
				return NewServiceOutput(InfraError, err.Error())
			}

			archiveFilesList[archiveFileIndex].DownloadUrl = updatedUrl
		}
	}

	return NewServiceOutput(Success, archiveFilesList)
}

func (service *ContainerImageService) CreateArchiveFile(
	input map[string]interface{},
	shouldSchedule bool,
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

	var compressionFormatPtr *valueObject.CompressionFormat
	if input["compressionFormat"] != nil {
		compressionFormat, err := valueObject.NewCompressionFormat(
			input["compressionFormat"],
		)
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		compressionFormatPtr = &compressionFormat
	}

	if shouldSchedule {
		cliCmd := infraEnvs.InfiniteEzBinary + " container image archive create"
		createParams := []string{
			"--account-id", accountId.String(),
			"--image-id", imageId.String(),
		}

		if compressionFormatPtr != nil {
			createParams = append(
				createParams, "--compression-format", compressionFormatPtr.String(),
			)
		}

		cliCmd = cliCmd + " " + strings.Join(createParams, " ")

		scheduledTaskCmdRepo := infra.NewScheduledTaskCmdRepo(service.persistentDbSvc)
		taskName, _ := valueObject.NewScheduledTaskName("CreateContainerImageArchiveFile")
		taskCmd, _ := valueObject.NewUnixCommand(cliCmd)
		taskTag, _ := valueObject.NewScheduledTaskTag("containerImageArchive")
		taskTags := []valueObject.ScheduledTaskTag{taskTag}
		timeoutSeconds := uint16(1800)

		scheduledTaskCreateDto := dto.NewCreateScheduledTask(
			taskName, taskCmd, taskTags, &timeoutSeconds, nil,
		)

		err = useCase.CreateScheduledTask(scheduledTaskCmdRepo, scheduledTaskCreateDto)
		if err != nil {
			return NewServiceOutput(InfraError, err.Error())
		}

		return NewServiceOutput(Created, "ContainerImageArchiveFileCreationScheduled")
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

	createDto := dto.NewCreateContainerImageArchiveFile(
		accountId, imageId, compressionFormatPtr, operatorAccountId, operatorIpAddress,
	)

	containerImageCmdRepo := infra.NewContainerImageCmdRepo(service.persistentDbSvc)
	accountQueryRepo := infra.NewAccountQueryRepo(service.persistentDbSvc)

	archiveFile, err := useCase.CreateContainerImageArchiveFile(
		service.containerImageQueryRepo, containerImageCmdRepo, accountQueryRepo,
		service.activityRecordCmdRepo, createDto,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, archiveFile)
}

func (service *ContainerImageService) DeleteArchiveFile(
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

	operatorIpAddress := LocalOperatorIpAddress
	if input["operatorIpAddress"] != nil {
		operatorIpAddress, err = valueObject.NewIpAddress(input["operatorIpAddress"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	deleteDto := dto.NewDeleteContainerImageArchiveFile(
		accountId, imageId, operatorAccountId, operatorIpAddress,
	)

	containerImageCmdRepo := infra.NewContainerImageCmdRepo(service.persistentDbSvc)

	err = useCase.DeleteContainerImageArchiveFile(
		service.containerImageQueryRepo, containerImageCmdRepo,
		service.activityRecordCmdRepo, deleteDto,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, "ContainerImageArchiveFileDeleted")
}
