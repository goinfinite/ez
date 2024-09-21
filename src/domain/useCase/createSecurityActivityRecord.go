package useCase

import (
	"log/slog"

	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/repository"
	"github.com/speedianet/control/src/domain/valueObject"
)

type CreateSecurityActivityRecord struct {
	activityRecordCmdRepo repository.ActivityRecordCmdRepo
	recordLevel           valueObject.ActivityRecordLevel
}

func NewCreateSecurityActivityRecord(
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
) *CreateSecurityActivityRecord {
	recordLevel, _ := valueObject.NewActivityRecordLevel("SEC")
	return &CreateSecurityActivityRecord{
		activityRecordCmdRepo: activityRecordCmdRepo,
		recordLevel:           recordLevel,
	}
}

func (uc *CreateSecurityActivityRecord) createActivityRecord(
	createDto dto.CreateActivityRecord,
) {
	err := uc.activityRecordCmdRepo.Create(createDto)
	if err != nil {
		slog.Debug(
			"CreateSecurityActivityRecordError",
			slog.Any("createDto", createDto),
			slog.Any("error", err),
		)
	}
}

func (uc *CreateSecurityActivityRecord) CreateSessionToken(
	recordCode valueObject.ActivityRecordCode,
	createDto dto.CreateSessionToken,
) {
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel:       uc.recordLevel,
		RecordCode:        recordCode,
		OperatorIpAddress: &createDto.OperatorIpAddress,
		RecordDetails:     createDto.Username,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) CreateContainerSessionToken(
	createDto dto.CreateContainerSessionToken,
) {
	recordCode, _ := valueObject.NewActivityRecordCode("ContainerSessionTokenGenerated")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel:       uc.recordLevel,
		RecordCode:        recordCode,
		OperatorAccountId: &createDto.OperatorAccountId,
		OperatorIpAddress: &createDto.OperatorIpAddress,
		AccountId:         &createDto.AccountId,
		ContainerId:       &createDto.ContainerId,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) CreateAccount(
	createDto dto.CreateAccount,
	accountId valueObject.AccountId,
) {
	recordCode, _ := valueObject.NewActivityRecordCode("AccountCreated")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel:       uc.recordLevel,
		RecordCode:        recordCode,
		OperatorAccountId: &createDto.OperatorAccountId,
		OperatorIpAddress: &createDto.OperatorIpAddress,
		AccountId:         &accountId,
		RecordDetails:     createDto.Username,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) UpdateAccount(
	updateDto dto.UpdateAccount,
) {
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel:       uc.recordLevel,
		OperatorAccountId: &updateDto.OperatorAccountId,
		OperatorIpAddress: &updateDto.OperatorIpAddress,
		AccountId:         &updateDto.AccountId,
		RecordDetails:     updateDto,
	}

	codeStr := "AccountUpdated"
	if updateDto.Password != nil {
		codeStr = "AccountPasswordUpdated"
		createRecordDto.RecordDetails = nil
	}
	recordCode, _ := valueObject.NewActivityRecordCode(codeStr)
	createRecordDto.RecordCode = recordCode

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) DeleteAccount(
	deleteDto dto.DeleteAccount,
) {
	recordCode, _ := valueObject.NewActivityRecordCode("AccountDeleted")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel:       uc.recordLevel,
		RecordCode:        recordCode,
		OperatorAccountId: &deleteDto.OperatorAccountId,
		OperatorIpAddress: &deleteDto.OperatorIpAddress,
		AccountId:         &deleteDto.AccountId,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) CreateContainer(
	createDto dto.CreateContainer,
	containerId valueObject.ContainerId,
) {
	recordCode, _ := valueObject.NewActivityRecordCode("ContainerCreated")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel:       uc.recordLevel,
		RecordCode:        recordCode,
		OperatorAccountId: &createDto.OperatorAccountId,
		OperatorIpAddress: &createDto.OperatorIpAddress,
		AccountId:         &createDto.AccountId,
		ContainerId:       &containerId,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) UpdateContainer(
	updateDto dto.UpdateContainer,
) {
	recordCode, _ := valueObject.NewActivityRecordCode("ContainerUpdated")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel:       uc.recordLevel,
		RecordCode:        recordCode,
		OperatorAccountId: &updateDto.OperatorAccountId,
		OperatorIpAddress: &updateDto.OperatorIpAddress,
		AccountId:         &updateDto.AccountId,
		ContainerId:       &updateDto.ContainerId,
		RecordDetails:     updateDto,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) DeleteContainer(
	deleteDto dto.DeleteContainer,
) {
	recordCode, _ := valueObject.NewActivityRecordCode("ContainerDeleted")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel:       uc.recordLevel,
		RecordCode:        recordCode,
		OperatorAccountId: &deleteDto.OperatorAccountId,
		OperatorIpAddress: &deleteDto.OperatorIpAddress,
		AccountId:         &deleteDto.AccountId,
		ContainerId:       &deleteDto.ContainerId,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) CreateContainerProfile(
	createDto dto.CreateContainerProfile,
	profileId valueObject.ContainerProfileId,
) {
	recordCode, _ := valueObject.NewActivityRecordCode("ContainerProfileCreated")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel:        uc.recordLevel,
		RecordCode:         recordCode,
		OperatorAccountId:  &createDto.OperatorAccountId,
		OperatorIpAddress:  &createDto.OperatorIpAddress,
		AccountId:          &createDto.AccountId,
		ContainerProfileId: &profileId,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) UpdateContainerProfile(
	updateDto dto.UpdateContainerProfile,
) {
	recordCode, _ := valueObject.NewActivityRecordCode("ContainerProfileUpdated")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel:        uc.recordLevel,
		RecordCode:         recordCode,
		OperatorAccountId:  &updateDto.OperatorAccountId,
		OperatorIpAddress:  &updateDto.OperatorIpAddress,
		AccountId:          &updateDto.AccountId,
		ContainerProfileId: &updateDto.ProfileId,
		RecordDetails:      updateDto,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) DeleteContainerProfile(
	deleteDto dto.DeleteContainerProfile,
) {
	recordCode, _ := valueObject.NewActivityRecordCode("ContainerProfileDeleted")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel:        uc.recordLevel,
		RecordCode:         recordCode,
		OperatorAccountId:  &deleteDto.OperatorAccountId,
		OperatorIpAddress:  &deleteDto.OperatorIpAddress,
		AccountId:          &deleteDto.AccountId,
		ContainerProfileId: &deleteDto.ProfileId,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) CreateMapping(
	createDto dto.CreateMapping,
	mappingId valueObject.MappingId,
) {
	recordCode, _ := valueObject.NewActivityRecordCode("MappingCreated")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel:       uc.recordLevel,
		RecordCode:        recordCode,
		OperatorAccountId: &createDto.OperatorAccountId,
		OperatorIpAddress: &createDto.OperatorIpAddress,
		AccountId:         &createDto.AccountId,
		MappingId:         &mappingId,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) DeleteMapping(
	deleteDto dto.DeleteMapping,
) {
	recordCode, _ := valueObject.NewActivityRecordCode("MappingDeleted")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel:       uc.recordLevel,
		RecordCode:        recordCode,
		OperatorAccountId: &deleteDto.OperatorAccountId,
		OperatorIpAddress: &deleteDto.OperatorIpAddress,
		MappingId:         &deleteDto.MappingId,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) CreateMappingTarget(
	createDto dto.CreateMappingTarget,
	targetId valueObject.MappingTargetId,
) {
	recordCode, _ := valueObject.NewActivityRecordCode("MappingTargetCreated")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel:       uc.recordLevel,
		RecordCode:        recordCode,
		OperatorAccountId: &createDto.OperatorAccountId,
		OperatorIpAddress: &createDto.OperatorIpAddress,
		MappingId:         &createDto.MappingId,
		MappingTargetId:   &targetId,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) DeleteMappingTarget(
	deleteDto dto.DeleteMappingTarget,
) {
	recordCode, _ := valueObject.NewActivityRecordCode("MappingTargetDeleted")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel:       uc.recordLevel,
		RecordCode:        recordCode,
		OperatorAccountId: &deleteDto.OperatorAccountId,
		OperatorIpAddress: &deleteDto.OperatorIpAddress,
		MappingId:         &deleteDto.MappingId,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) CreateContainerSnapshotImage(
	createDto dto.CreateContainerSnapshotImage,
	imageAccountId valueObject.AccountId,
	imageId valueObject.ContainerImageId,
) {
	recordCode, _ := valueObject.NewActivityRecordCode("ContainerSnapshotImageCreated")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel:       uc.recordLevel,
		RecordCode:        recordCode,
		OperatorAccountId: &createDto.OperatorAccountId,
		OperatorIpAddress: &createDto.OperatorIpAddress,
		AccountId:         &imageAccountId,
		ContainerId:       &createDto.ContainerId,
		ContainerImageId:  &imageId,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) DeleteContainerImage(
	deleteDto dto.DeleteContainerImage,
) {
	recordCode, _ := valueObject.NewActivityRecordCode("ContainerImageDeleted")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel:       uc.recordLevel,
		RecordCode:        recordCode,
		OperatorAccountId: &deleteDto.OperatorAccountId,
		OperatorIpAddress: &deleteDto.OperatorIpAddress,
		AccountId:         &deleteDto.AccountId,
		ContainerImageId:  &deleteDto.ImageId,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) CreateContainerImageArchiveFile(
	createDto dto.CreateContainerImageArchiveFile,
) {
	recordCode, _ := valueObject.NewActivityRecordCode("ContainerImageArchiveFileCreated")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel:       uc.recordLevel,
		RecordCode:        recordCode,
		OperatorAccountId: &createDto.OperatorAccountId,
		OperatorIpAddress: &createDto.OperatorIpAddress,
		AccountId:         &createDto.AccountId,
		ContainerImageId:  &createDto.ImageId,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) ImportContainerImageArchiveFile(
	importDto dto.ImportContainerImageArchiveFile,
	imageId valueObject.ContainerImageId,
) {
	recordCode, _ := valueObject.NewActivityRecordCode("ContainerImageArchiveFileImported")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel:       uc.recordLevel,
		RecordCode:        recordCode,
		OperatorAccountId: &importDto.OperatorAccountId,
		OperatorIpAddress: &importDto.OperatorIpAddress,
		AccountId:         &importDto.AccountId,
		ContainerImageId:  &imageId,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) DeleteContainerImageArchiveFile(
	deleteDto dto.DeleteContainerImageArchiveFile,
) {
	recordCode, _ := valueObject.NewActivityRecordCode("ContainerImageArchiveFileDeleted")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel:       uc.recordLevel,
		RecordCode:        recordCode,
		OperatorAccountId: &deleteDto.OperatorAccountId,
		OperatorIpAddress: &deleteDto.OperatorIpAddress,
		AccountId:         &deleteDto.AccountId,
		ContainerImageId:  &deleteDto.ImageId,
	}

	uc.createActivityRecord(createRecordDto)
}
