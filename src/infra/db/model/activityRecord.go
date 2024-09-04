package dbModel

import (
	"time"

	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/valueObject"
)

type ActivityRecord struct {
	ID                 uint64 `gorm:"primarykey"`
	RecordLevel        string `gorm:"not null"`
	RecordCode         string `gorm:"not null"`
	OperatorAccountId  *uint64
	OperatorIpAddress  *string
	AccountId          *uint64
	ContainerId        *string
	ContainerProfileId *uint64
	ContainerImageId   *string
	MappingId          *uint64
	MappingTargetId    *uint64
	ScheduledTaskId    *uint64
	RecordDetails      *string
	CreatedAt          time.Time `gorm:"not null"`
}

func (ActivityRecord) TableName() string {
	return "activity_records"
}

func NewActivityRecord(
	recordId uint64,
	recordLevel, recordCode string,
	operatorAccountId *uint64,
	operatorIpAddress *string,
	accountId *uint64,
	containerId *string,
	containerProfileId *uint64,
	containerImageId *string,
	mappingId, mappingTargetId, scheduledTaskId *uint64,
	recordDetails *string,
) ActivityRecord {
	model := ActivityRecord{
		RecordLevel:        recordLevel,
		RecordCode:         recordCode,
		OperatorAccountId:  operatorAccountId,
		OperatorIpAddress:  operatorIpAddress,
		AccountId:          accountId,
		ContainerId:        containerId,
		ContainerProfileId: containerProfileId,
		ContainerImageId:   containerImageId,
		MappingId:          mappingId,
		MappingTargetId:    mappingTargetId,
		ScheduledTaskId:    scheduledTaskId,
		RecordDetails:      recordDetails,
	}

	if recordId != 0 {
		model.ID = recordId
	}

	return model
}

func (model ActivityRecord) ToEntity() (recordEntity entity.ActivityRecord, err error) {
	recordId, err := valueObject.NewActivityRecordId(model.ID)
	if err != nil {
		return recordEntity, err
	}

	recordLevel, err := valueObject.NewActivityRecordLevel(model.RecordLevel)
	if err != nil {
		return recordEntity, err
	}

	recordCode, err := valueObject.NewActivityRecordCode(model.RecordCode)
	if err != nil {
		return recordEntity, err
	}

	var operatorAccountIdPtr *valueObject.AccountId
	if model.OperatorAccountId != nil {
		operatorAccountId, err := valueObject.NewAccountId(*model.OperatorAccountId)
		if err != nil {
			return recordEntity, err
		}
		operatorAccountIdPtr = &operatorAccountId
	}

	var operatorIpAddressPtr *valueObject.IpAddress
	if model.OperatorIpAddress != nil {
		operatorIpAddress, err := valueObject.NewIpAddress(*model.OperatorIpAddress)
		if err != nil {
			return recordEntity, err
		}
		operatorIpAddressPtr = &operatorIpAddress
	}

	var accountIdPtr *valueObject.AccountId
	if model.AccountId != nil {
		accountId, err := valueObject.NewAccountId(*model.AccountId)
		if err != nil {
			return recordEntity, err
		}
		accountIdPtr = &accountId
	}

	var containerIdPtr *valueObject.ContainerId
	if model.ContainerId != nil {
		containerId, err := valueObject.NewContainerId(*model.ContainerId)
		if err != nil {
			return recordEntity, err
		}
		containerIdPtr = &containerId
	}

	var containerProfileIdPtr *valueObject.ContainerProfileId
	if model.ContainerProfileId != nil {
		containerProfileId, err := valueObject.NewContainerProfileId(*model.ContainerProfileId)
		if err != nil {
			return recordEntity, err
		}
		containerProfileIdPtr = &containerProfileId
	}

	var containerImageIdPtr *valueObject.ContainerImageId
	if model.ContainerImageId != nil {
		containerImageId, err := valueObject.NewContainerImageId(*model.ContainerImageId)
		if err != nil {
			return recordEntity, err
		}
		containerImageIdPtr = &containerImageId
	}

	var mappingIdPtr *valueObject.MappingId
	if model.MappingId != nil {
		mappingId, err := valueObject.NewMappingId(*model.MappingId)
		if err != nil {
			return recordEntity, err
		}
		mappingIdPtr = &mappingId
	}

	var mappingTargetIdPtr *valueObject.MappingTargetId
	if model.MappingTargetId != nil {
		mappingTargetId, err := valueObject.NewMappingTargetId(*model.MappingTargetId)
		if err != nil {
			return recordEntity, err
		}
		mappingTargetIdPtr = &mappingTargetId
	}

	var scheduledTaskIdPtr *valueObject.ScheduledTaskId
	if model.ScheduledTaskId != nil {
		scheduledTaskId, err := valueObject.NewScheduledTaskId(*model.ScheduledTaskId)
		if err != nil {
			return recordEntity, err
		}
		scheduledTaskIdPtr = &scheduledTaskId
	}

	var recordDetails interface{}
	if model.RecordDetails != nil {
		recordDetails = *model.RecordDetails
	}

	createdAt := valueObject.NewUnixTimeWithGoTime(model.CreatedAt)

	return entity.NewActivityRecord(
		recordId, recordLevel, recordCode, operatorAccountIdPtr, operatorIpAddressPtr,
		accountIdPtr, containerIdPtr, containerProfileIdPtr, containerImageIdPtr,
		mappingIdPtr, mappingTargetIdPtr, scheduledTaskIdPtr, recordDetails, createdAt,
	)
}
