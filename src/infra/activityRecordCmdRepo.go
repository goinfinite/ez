package infra

import (
	"encoding/json"

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
	var operatorAccountIdPtr *uint64
	if createDto.OperatorAccountId != nil {
		operatorAccountId := createDto.OperatorAccountId.Uint64()
		operatorAccountIdPtr = &operatorAccountId
	}

	var operatorIpAddressPtr *string
	if createDto.OperatorIpAddress != nil {
		operatorIpAddress := createDto.OperatorIpAddress.String()
		operatorIpAddressPtr = &operatorIpAddress
	}

	var accountIdPtr *uint64
	if createDto.AccountId != nil {
		accountId := createDto.AccountId.Uint64()
		accountIdPtr = &accountId
	}

	var containerIdPtr *string
	if createDto.ContainerId != nil {
		containerId := createDto.ContainerId.String()
		containerIdPtr = &containerId
	}

	var containerProfileIdPtr *uint64
	if createDto.ContainerProfileId != nil {
		containerProfileId := createDto.ContainerProfileId.Uint64()
		containerProfileIdPtr = &containerProfileId
	}

	var containerImageIdPtr *string
	if createDto.ContainerImageId != nil {
		containerImageId := createDto.ContainerImageId.String()
		containerImageIdPtr = &containerImageId
	}

	var mappingIdPtr *uint64
	if createDto.MappingId != nil {
		mappingId := createDto.MappingId.Uint64()
		mappingIdPtr = &mappingId
	}

	var mappingTargetIdPtr *uint64
	if createDto.MappingTargetId != nil {
		mappingTargetId := createDto.MappingTargetId.Uint64()
		mappingTargetIdPtr = &mappingTargetId
	}

	var scheduledTaskIdPtr *uint64
	if createDto.ScheduledTaskId != nil {
		scheduledTaskId := createDto.ScheduledTaskId.Uint64()
		scheduledTaskIdPtr = &scheduledTaskId
	}

	var recordDetails *string
	if createDto.RecordDetails != nil {
		recordDetailsBytes, err := json.Marshal(createDto.RecordDetails)
		if err != nil {
			return err
		}
		recordDetailsStr := string(recordDetailsBytes)
		recordDetails = &recordDetailsStr
	}

	activityRecordModel := dbModel.NewActivityRecord(
		0, createDto.RecordLevel.String(), createDto.RecordCode.String(),
		operatorAccountIdPtr, operatorIpAddressPtr, accountIdPtr, containerIdPtr,
		containerProfileIdPtr, containerImageIdPtr, mappingIdPtr, mappingTargetIdPtr,
		scheduledTaskIdPtr, recordDetails,
	)

	return repo.trailDbSvc.Handler.Create(&activityRecordModel).Error
}

func (repo *ActivityRecordCmdRepo) Delete(deleteDto dto.DeleteActivityRecords) error {
	deleteModel := dbModel.ActivityRecord{}
	if deleteDto.RecordId != nil {
		deleteModel.ID = deleteDto.RecordId.Uint64()
	}

	if deleteDto.RecordLevel != nil {
		deleteModel.RecordLevel = deleteDto.RecordLevel.String()
	}

	if deleteDto.RecordCode != nil {
		deleteModel.RecordCode = deleteDto.RecordCode.String()
	}

	if deleteDto.OperatorAccountId != nil {
		operatorAccountId := deleteDto.OperatorAccountId.Uint64()
		deleteModel.OperatorAccountId = &operatorAccountId
	}

	if deleteDto.OperatorIpAddress != nil {
		operatorIpAddressStr := deleteDto.OperatorIpAddress.String()
		deleteModel.OperatorIpAddress = &operatorIpAddressStr
	}

	if deleteDto.AccountId != nil {
		accountId := deleteDto.AccountId.Uint64()
		deleteModel.AccountId = &accountId
	}

	if deleteDto.ContainerId != nil {
		containerIdStr := deleteDto.ContainerId.String()
		deleteModel.ContainerId = &containerIdStr
	}

	if deleteDto.ContainerProfileId != nil {
		containerProfileId := deleteDto.ContainerProfileId.Uint64()
		deleteModel.ContainerProfileId = &containerProfileId
	}

	if deleteDto.ContainerImageId != nil {
		containerImageIdStr := deleteDto.ContainerImageId.String()
		deleteModel.ContainerImageId = &containerImageIdStr
	}

	if deleteDto.MappingId != nil {
		mappingId := deleteDto.MappingId.Uint64()
		deleteModel.MappingId = &mappingId
	}

	if deleteDto.MappingTargetId != nil {
		mappingTargetId := deleteDto.MappingTargetId.Uint64()
		deleteModel.MappingTargetId = &mappingTargetId
	}

	if deleteDto.ScheduledTaskId != nil {
		scheduledTaskId := deleteDto.ScheduledTaskId.Uint64()
		deleteModel.ScheduledTaskId = &scheduledTaskId
	}

	dbQuery := repo.trailDbSvc.Handler.Where(&deleteModel)

	if deleteDto.CreatedBeforeAt != nil {
		dbQuery.Where("created_at < ?", deleteDto.CreatedBeforeAt.GetAsGoTime())
	}
	if deleteDto.CreatedAfterAt != nil {
		dbQuery.Where("created_at > ?", deleteDto.CreatedAfterAt.GetAsGoTime())
	}

	return dbQuery.Delete(&dbModel.ActivityRecord{}).Error
}
