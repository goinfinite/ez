package dto

import (
	"github.com/speedianet/control/src/domain/valueObject"
)

type CreateActivityRecord struct {
	RecordLevel        valueObject.ActivityRecordLevel `json:"recordLevel"`
	RecordCode         valueObject.ActivityRecordCode  `json:"recordCode"`
	OperatorAccountId  *valueObject.AccountId          `json:"operatorAccountId,omitempty"`
	OperatorIpAddress  *valueObject.IpAddress          `json:"operatorIpAddress,omitempty"`
	AccountId          *valueObject.AccountId          `json:"accountId,omitempty"`
	ContainerId        *valueObject.ContainerId        `json:"containerId,omitempty"`
	ContainerProfileId *valueObject.ContainerProfileId `json:"containerProfileId,omitempty"`
	ContainerImageId   *valueObject.ContainerImageId   `json:"containerImageId,omitempty"`
	MappingId          *valueObject.MappingId          `json:"mappingId,omitempty"`
	ScheduledTaskId    *valueObject.ScheduledTaskId    `json:"scheduledTaskId,omitempty"`
	RecordDetails      interface{}                     `json:"recordDetails,omitempty"`
}
