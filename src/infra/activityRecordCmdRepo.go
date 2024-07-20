package infra

import (
	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/infra/db"
	dbModel "github.com/speedianet/control/src/infra/db/model"
)

type ActivityRecordCmdRepo struct {
	trailDbSvc *db.TrailDatabaseService
}

func NewActivityRecordCmdRepo(
	trailDbSvc *db.TrailDatabaseService,
) *ActivityRecordCmdRepo {
	return &ActivityRecordCmdRepo{trailDbSvc: trailDbSvc}
}

func (repo *ActivityRecordCmdRepo) Create(createDto dto.CreateActivityRecord) error {
	var codePtr *string
	if createDto.Code != nil {
		code := createDto.Code.String()
		codePtr = &code
	}

	var messagePtr *string
	if createDto.Message != nil {
		message := createDto.Message.String()
		messagePtr = &message
	}

	var ipAddressPtr *string
	if createDto.IpAddress != nil {
		ipAddress := createDto.IpAddress.String()
		ipAddressPtr = &ipAddress
	}

	var operatorAccountIdPtr *uint
	if createDto.OperatorAccountId != nil {
		operatorAccountId := uint(createDto.OperatorAccountId.Uint64())
		operatorAccountIdPtr = &operatorAccountId
	}

	var targetAccountIdPtr *uint
	if createDto.TargetAccountId != nil {
		targetAccountId := uint(createDto.TargetAccountId.Uint64())
		targetAccountIdPtr = &targetAccountId
	}

	var usernamePtr *string
	if createDto.Username != nil {
		username := createDto.Username.String()
		usernamePtr = &username
	}

	var containerIdPtr *string
	if createDto.ContainerId != nil {
		containerId := createDto.ContainerId.String()
		containerIdPtr = &containerId
	}

	var containerProfileIdPtr *uint
	if createDto.ContainerProfileId != nil {
		containerProfileId := uint(createDto.ContainerProfileId.Read())
		containerProfileIdPtr = &containerProfileId
	}

	var mappingIdPtr *uint
	if createDto.MappingId != nil {
		mappingId := uint(createDto.MappingId.Read())
		mappingIdPtr = &mappingId
	}

	securityEventModel := dbModel.NewActivityRecord(
		0, createDto.Level.String(), codePtr, messagePtr, ipAddressPtr,
		operatorAccountIdPtr, targetAccountIdPtr, usernamePtr, containerIdPtr,
		containerProfileIdPtr, mappingIdPtr,
	)

	return repo.trailDbSvc.Handler.Create(&securityEventModel).Error
}

func (repo *ActivityRecordCmdRepo) Delete(deleteDto dto.DeleteActivityRecords) error {
	deleteModel := dbModel.ActivityRecord{}
	if deleteDto.Id != nil {
		deleteModel.ID = uint(deleteDto.Id.Read())
	}

	if deleteDto.Level != nil {
		deleteModel.Level = deleteDto.Level.String()
	}

	if deleteDto.Code != nil {
		codeStr := deleteDto.Code.String()
		deleteModel.Code = &codeStr
	}

	if deleteDto.Message != nil {
		messageStr := deleteDto.Message.String()
		deleteModel.Message = &messageStr
	}

	if deleteDto.IpAddress != nil {
		ipAddressStr := deleteDto.IpAddress.String()
		deleteModel.IpAddress = &ipAddressStr
	}

	if deleteDto.OperatorAccountId != nil {
		operatorAccountId := uint(deleteDto.OperatorAccountId.Uint64())
		deleteModel.OperatorAccountId = &operatorAccountId
	}

	if deleteDto.TargetAccountId != nil {
		targetAccountId := uint(deleteDto.TargetAccountId.Uint64())
		deleteModel.TargetAccountId = &targetAccountId
	}

	if deleteDto.Username != nil {
		usernameStr := deleteDto.Username.String()
		deleteModel.Username = &usernameStr
	}

	if deleteDto.ContainerId != nil {
		containerIdStr := deleteDto.ContainerId.String()
		deleteModel.ContainerId = &containerIdStr
	}

	if deleteDto.ContainerProfileId != nil {
		containerProfileId := uint(deleteDto.ContainerProfileId.Read())
		deleteModel.ContainerProfileId = &containerProfileId
	}

	if deleteDto.MappingId != nil {
		mappingId := uint(deleteDto.MappingId.Read())
		deleteModel.MappingId = &mappingId
	}

	dbQuery := repo.trailDbSvc.Handler.Where(&deleteModel)

	if deleteDto.CreatedAt != nil {
		dbQuery.Where("created_at >= ?", deleteDto.CreatedAt.GetAsGoTime())
	}

	return dbQuery.Delete(&dbModel.ActivityRecord{}).Error
}
