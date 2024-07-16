package dto

import (
	"github.com/speedianet/control/src/domain/valueObject"
)

type DeleteActivityRecords struct {
	Id                 *valueObject.ActivityRecordId      `json:"id,omitempty"`
	Level              *valueObject.ActivityRecordLevel   `json:"level,omitempty"`
	Code               *valueObject.ActivityRecordCode    `json:"code,omitempty"`
	Message            *valueObject.ActivityRecordMessage `json:"message,omitempty"`
	IpAddress          *valueObject.IpAddress             `json:"ipAddress,omitempty"`
	OperatorAccountId  *valueObject.AccountId             `json:"operatorAccountId,omitempty"`
	TargetAccountId    *valueObject.AccountId             `json:"targetAccountId,omitempty"`
	Username           *valueObject.Username              `json:"username,omitempty"`
	ContainerId        *valueObject.ContainerId           `json:"containerId,omitempty"`
	ContainerProfileId *valueObject.ContainerProfileId    `json:"containerProfileId,omitempty"`
	MappingId          *valueObject.MappingId             `json:"mappingId,omitempty"`
	CreatedAt          *valueObject.UnixTime              `json:"createdAt,omitempty"`
}

func NewDeleteActivityRecords(
	id *valueObject.ActivityRecordId,
	level *valueObject.ActivityRecordLevel,
	code *valueObject.ActivityRecordCode,
	message *valueObject.ActivityRecordMessage,
	ipAddress *valueObject.IpAddress,
	operatorAccountId *valueObject.AccountId,
	targetAccountId *valueObject.AccountId,
	username *valueObject.Username,
	containerId *valueObject.ContainerId,
	containerProfileId *valueObject.ContainerProfileId,
	mappingId *valueObject.MappingId,
	createdAt *valueObject.UnixTime,
) DeleteActivityRecords {
	return DeleteActivityRecords{
		Id:                 id,
		Level:              level,
		Code:               code,
		Message:            message,
		IpAddress:          ipAddress,
		OperatorAccountId:  operatorAccountId,
		TargetAccountId:    targetAccountId,
		Username:           username,
		ContainerId:        containerId,
		ContainerProfileId: containerProfileId,
		MappingId:          mappingId,
		CreatedAt:          createdAt,
	}
}
