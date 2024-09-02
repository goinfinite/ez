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
	loginDto dto.Login,
	recordCode valueObject.ActivityRecordCode,
) {
	createRecordDto := dto.CreateActivityRecord{
		Level:             uc.recordLevel,
		Code:              &recordCode,
		OperatorIpAddress: loginDto.OperatorIpAddress,
		Username:          &loginDto.Username,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) CreateAccount(
	createDto dto.CreateAccount,
	accountId valueObject.AccountId,
) {
	recordCode, _ := valueObject.NewActivityRecordCode("AccountCreated")
	createRecordDto := dto.CreateActivityRecord{
		Level:             uc.recordLevel,
		Code:              &recordCode,
		OperatorAccountId: &createDto.OperatorAccountId,
		OperatorIpAddress: &createDto.OperatorIpAddress,
		TargetAccountId:   &accountId,
		Username:          &createDto.Username,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) UpdateAccount(
	updateCode valueObject.ActivityRecordCode,
	updateDto dto.UpdateAccount,
) {
	createRecordDto := dto.CreateActivityRecord{
		Level:             uc.recordLevel,
		Code:              &updateCode,
		OperatorAccountId: &updateDto.OperatorAccountId,
		OperatorIpAddress: &updateDto.OperatorIpAddress,
		TargetAccountId:   &updateDto.AccountId,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) DeleteAccount(
	deleteDto dto.DeleteAccount,
) {
	recordCode, _ := valueObject.NewActivityRecordCode("AccountDeleted")
	createRecordDto := dto.CreateActivityRecord{
		Level:             uc.recordLevel,
		Code:              &recordCode,
		OperatorAccountId: &deleteDto.OperatorAccountId,
		OperatorIpAddress: &deleteDto.OperatorIpAddress,
		TargetAccountId:   &deleteDto.AccountId,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) CreateContainerSnapshotImage(
	createDto dto.CreateContainerSnapshotImage,
	imageId valueObject.ContainerImageId,
) {
	recordCode, _ := valueObject.NewActivityRecordCode("ContainerSnapshotImageCreated")
	createRecordDto := dto.CreateActivityRecord{
		Level:             uc.recordLevel,
		Code:              &recordCode,
		OperatorAccountId: &createDto.OperatorAccountId,
		OperatorIpAddress: &createDto.OperatorIpAddress,
		TargetAccountId:   &createDto.AccountId,
		ContainerId:       &createDto.ContainerId,
		ContainerImageId:  &imageId,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) CreateContainerImageArchiveFile(
	createDto dto.CreateContainerImageArchiveFile,
) {
	recordCode, _ := valueObject.NewActivityRecordCode("ContainerImageArchiveFileCreated")
	createRecordDto := dto.CreateActivityRecord{
		Level:             uc.recordLevel,
		Code:              &recordCode,
		OperatorAccountId: &createDto.OperatorAccountId,
		OperatorIpAddress: &createDto.OperatorIpAddress,
		TargetAccountId:   &createDto.AccountId,
		ContainerImageId:  &createDto.ImageId,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) ImportContainerImageArchiveFile(
	importDto dto.ImportContainerImageArchiveFile,
) {
	recordCode, _ := valueObject.NewActivityRecordCode("ContainerImageArchiveFileImported")
	createRecordDto := dto.CreateActivityRecord{
		Level:             uc.recordLevel,
		Code:              &recordCode,
		OperatorAccountId: &importDto.OperatorAccountId,
		OperatorIpAddress: &importDto.OperatorIpAddress,
		TargetAccountId:   &importDto.AccountId,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) DeleteContainerImage(
	deleteDto dto.DeleteContainerImage,
) {
	recordCode, _ := valueObject.NewActivityRecordCode("ContainerImageDeleted")
	createRecordDto := dto.CreateActivityRecord{
		Level:             uc.recordLevel,
		Code:              &recordCode,
		OperatorAccountId: &deleteDto.OperatorAccountId,
		OperatorIpAddress: &deleteDto.OperatorIpAddress,
		TargetAccountId:   &deleteDto.AccountId,
		ContainerImageId:  &deleteDto.ImageId,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) DeleteContainerImageArchiveFile(
	deleteDto dto.DeleteContainerImageArchiveFile,
) {
	recordCode, _ := valueObject.NewActivityRecordCode("ContainerImageArchiveFileDeleted")
	createRecordDto := dto.CreateActivityRecord{
		Level:             uc.recordLevel,
		Code:              &recordCode,
		OperatorAccountId: &deleteDto.OperatorAccountId,
		OperatorIpAddress: &deleteDto.OperatorIpAddress,
		TargetAccountId:   &deleteDto.AccountId,
		ContainerImageId:  &deleteDto.ImageId,
	}

	uc.createActivityRecord(createRecordDto)
}
