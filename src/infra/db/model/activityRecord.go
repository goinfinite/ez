package dbModel

import (
	"time"

	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/valueObject"
)

type ActivityRecord struct {
	ID                 uint64 `gorm:"primarykey"`
	Level              string `gorm:"not null"`
	Code               *string
	Message            *string
	IpAddress          *string
	OperatorAccountId  *uint64
	TargetAccountId    *uint64
	Username           *string
	ContainerId        *string
	ContainerProfileId *uint64
	ContainerImageId   *string
	MappingId          *uint64
	CreatedAt          time.Time `gorm:"not null"`
}

func (ActivityRecord) TableName() string {
	return "activity_records"
}

func NewActivityRecord(
	id uint64,
	level string,
	code, message, ipAddress *string,
	operatorAccountId, targetAccountId *uint64,
	username, containerId *string,
	containerProfileId *uint64,
	containerImageId *string,
	mappingId *uint64,

) ActivityRecord {
	model := ActivityRecord{
		Level:              level,
		Code:               code,
		Message:            message,
		IpAddress:          ipAddress,
		OperatorAccountId:  operatorAccountId,
		TargetAccountId:    targetAccountId,
		Username:           username,
		ContainerId:        containerId,
		ContainerProfileId: containerProfileId,
		ContainerImageId:   containerImageId,
		MappingId:          mappingId,
	}

	if id != 0 {
		model.ID = id
	}

	return model
}

func (model ActivityRecord) ToEntity() (recordEntity entity.ActivityRecord, err error) {
	id, err := valueObject.NewActivityRecordId(model.ID)
	if err != nil {
		return recordEntity, err
	}

	level, err := valueObject.NewActivityRecordLevel(model.Level)
	if err != nil {
		return recordEntity, err
	}

	var codePtr *valueObject.ActivityRecordCode
	if model.Code != nil {
		code, err := valueObject.NewActivityRecordCode(*model.Code)
		if err != nil {
			return recordEntity, err
		}
		codePtr = &code
	}

	var messagePtr *valueObject.ActivityRecordMessage
	if model.Message != nil {
		message, err := valueObject.NewActivityRecordMessage(*model.Message)
		if err != nil {
			return recordEntity, err
		}
		messagePtr = &message
	}

	var ipAddressPtr *valueObject.IpAddress
	if model.IpAddress != nil {
		ipAddress, err := valueObject.NewIpAddress(*model.IpAddress)
		if err != nil {
			return recordEntity, err
		}
		ipAddressPtr = &ipAddress
	}

	var operatorAccountIdPtr *valueObject.AccountId
	if model.OperatorAccountId != nil {
		operatorAccountId, err := valueObject.NewAccountId(*model.OperatorAccountId)
		if err != nil {
			return recordEntity, err
		}
		operatorAccountIdPtr = &operatorAccountId
	}

	var targetAccountIdPtr *valueObject.AccountId
	if model.TargetAccountId != nil {
		targetAccountId, err := valueObject.NewAccountId(*model.TargetAccountId)
		if err != nil {
			return recordEntity, err
		}
		targetAccountIdPtr = &targetAccountId
	}

	var usernamePtr *valueObject.Username
	if model.Username != nil {
		username, err := valueObject.NewUsername(*model.Username)
		if err != nil {
			return recordEntity, err
		}
		usernamePtr = &username
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

	createdAt := valueObject.NewUnixTimeWithGoTime(model.CreatedAt)

	return entity.NewActivityRecord(
		id, level, codePtr, messagePtr, ipAddressPtr, operatorAccountIdPtr,
		targetAccountIdPtr, usernamePtr, containerIdPtr, containerProfileIdPtr,
		containerImageIdPtr, mappingIdPtr, createdAt,
	)
}
