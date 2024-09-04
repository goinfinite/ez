package entity

import (
	"github.com/speedianet/control/src/domain/valueObject"
)

type ActivityRecord struct {
	RecordId           valueObject.ActivityRecordId    `json:"recordId"`
	RecordLevel        valueObject.ActivityRecordLevel `json:"recordLevel"`
	RecordCode         valueObject.ActivityRecordCode  `json:"recordCode,omitempty"`
	OperatorAccountId  *valueObject.AccountId          `json:"operatorAccountId,omitempty"`
	OperatorIpAddress  *valueObject.IpAddress          `json:"operatorIpAddress,omitempty"`
	AccountId          *valueObject.AccountId          `json:"accountId,omitempty"`
	ContainerId        *valueObject.ContainerId        `json:"containerId,omitempty"`
	ContainerProfileId *valueObject.ContainerProfileId `json:"containerProfileId,omitempty"`
	ContainerImageId   *valueObject.ContainerImageId   `json:"containerImageId,omitempty"`
	MappingId          *valueObject.MappingId          `json:"mappingId,omitempty"`
	MappingTargetId    *valueObject.MappingTargetId    `json:"mappingTargetId,omitempty"`
	ScheduledTaskId    *valueObject.ScheduledTaskId    `json:"scheduledTaskId,omitempty"`
	RecordDetails      interface{}                     `json:"recordDetails,omitempty"`
	CreatedAt          valueObject.UnixTime            `json:"createdAt"`
}

func NewActivityRecord(
	recordId valueObject.ActivityRecordId,
	recordLevel valueObject.ActivityRecordLevel,
	recordCode valueObject.ActivityRecordCode,
	operatorAccountId *valueObject.AccountId,
	operatorIpAddress *valueObject.IpAddress,
	accountId *valueObject.AccountId,
	containerId *valueObject.ContainerId,
	containerProfileId *valueObject.ContainerProfileId,
	containerImageId *valueObject.ContainerImageId,
	mappingId *valueObject.MappingId,
	mappingTargetId *valueObject.MappingTargetId,
	scheduledTaskId *valueObject.ScheduledTaskId,
	recordDetails interface{},
	createdAt valueObject.UnixTime,
) (activityRecord ActivityRecord, err error) {
	return ActivityRecord{
		RecordId:           recordId,
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
		CreatedAt:          createdAt,
	}, nil
}
