package dto

import (
	"github.com/speedianet/control/src/domain/valueObject"
)

type DeleteActivityRecords struct {
	RecordId           *valueObject.ActivityRecordId      `json:"recordId,omitempty"`
	Level              *valueObject.ActivityRecordLevel   `json:"level,omitempty"`
	Code               *valueObject.ActivityRecordCode    `json:"code,omitempty"`
	Message            *valueObject.ActivityRecordMessage `json:"message,omitempty"`
	OperatorIpAddress  *valueObject.IpAddress             `json:"operatorIpAddress,omitempty"`
	OperatorAccountId  *valueObject.AccountId             `json:"operatorAccountId,omitempty"`
	TargetAccountId    *valueObject.AccountId             `json:"targetAccountId,omitempty"`
	Username           *valueObject.Username              `json:"username,omitempty"`
	ContainerId        *valueObject.ContainerId           `json:"containerId,omitempty"`
	ContainerProfileId *valueObject.ContainerProfileId    `json:"containerProfileId,omitempty"`
	ContainerImageId   *valueObject.ContainerImageId      `json:"containerImageId,omitempty"`
	MappingId          *valueObject.MappingId             `json:"mappingId,omitempty"`
	CreatedAt          *valueObject.UnixTime              `json:"createdAt,omitempty"`
}

func NewDeleteActivityRecords(
	recordId *valueObject.ActivityRecordId,
	level *valueObject.ActivityRecordLevel,
	code *valueObject.ActivityRecordCode,
	message *valueObject.ActivityRecordMessage,
	operatorIpAddress *valueObject.IpAddress,
	operatorAccountId *valueObject.AccountId,
	targetAccountId *valueObject.AccountId,
	username *valueObject.Username,
	containerId *valueObject.ContainerId,
	containerProfileId *valueObject.ContainerProfileId,
	containerImageId *valueObject.ContainerImageId,
	mappingId *valueObject.MappingId,
	createdAt *valueObject.UnixTime,
) DeleteActivityRecords {
	return DeleteActivityRecords{
		RecordId:           recordId,
		Level:              level,
		Code:               code,
		Message:            message,
		OperatorIpAddress:  operatorIpAddress,
		OperatorAccountId:  operatorAccountId,
		TargetAccountId:    targetAccountId,
		Username:           username,
		ContainerId:        containerId,
		ContainerProfileId: containerProfileId,
		ContainerImageId:   containerImageId,
		MappingId:          mappingId,
		CreatedAt:          createdAt,
	}
}
