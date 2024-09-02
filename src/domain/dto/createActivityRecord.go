package dto

import (
	"github.com/speedianet/control/src/domain/valueObject"
)

type CreateActivityRecord struct {
	Level              valueObject.ActivityRecordLevel    `json:"level"`
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
}
