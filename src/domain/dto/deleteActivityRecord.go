package dto

import (
	"github.com/speedianet/control/src/domain/valueObject"
)

type DeleteActivityRecords struct {
	RecordId           *valueObject.ActivityRecordId      `json:"recordId,omitempty"`
	Level              *valueObject.ActivityRecordLevel   `json:"level,omitempty"`
	Code               *valueObject.ActivityRecordCode    `json:"code,omitempty"`
	Message            *valueObject.ActivityRecordMessage `json:"message,omitempty"`
	OperatorAccountId  *valueObject.AccountId             `json:"operatorAccountId,omitempty"`
	OperatorIpAddress  *valueObject.IpAddress             `json:"operatorIpAddress,omitempty"`
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
	operatorAccountId *valueObject.AccountId,
	operatorIpAddress *valueObject.IpAddress,
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
		OperatorAccountId:  operatorAccountId,
		OperatorIpAddress:  operatorIpAddress,
		TargetAccountId:    targetAccountId,
		Username:           username,
		ContainerId:        containerId,
		ContainerProfileId: containerProfileId,
		ContainerImageId:   containerImageId,
		MappingId:          mappingId,
		CreatedAt:          createdAt,
	}
}
