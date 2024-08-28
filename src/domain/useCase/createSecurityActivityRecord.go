package useCase

import (
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
	go func() {
		_ = uc.activityRecordCmdRepo.Create(createDto)
	}()
}

func (uc *CreateSecurityActivityRecord) CreateSessionToken(
	loginDto dto.Login,
	recordCode valueObject.ActivityRecordCode,
) {
	createRecordDto := dto.CreateActivityRecord{
		Level:     uc.recordLevel,
		Code:      &recordCode,
		IpAddress: loginDto.IpAddress,
		Username:  &loginDto.Username,
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
		IpAddress:         &createDto.IpAddress,
		OperatorAccountId: &createDto.OperatorAccountId,
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
		IpAddress:         &updateDto.IpAddress,
		OperatorAccountId: &updateDto.OperatorAccountId,
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
		IpAddress:         &deleteDto.IpAddress,
		OperatorAccountId: &deleteDto.OperatorAccountId,
		TargetAccountId:   &deleteDto.AccountId,
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
		IpAddress:         &deleteDto.IpAddress,
		OperatorAccountId: &deleteDto.OperatorAccountId,
		TargetAccountId:   &deleteDto.AccountId,
		ContainerImageId:  &deleteDto.ImageId,
	}

	uc.createActivityRecord(createRecordDto)
}
