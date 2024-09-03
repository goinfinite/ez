package dto

import (
	"github.com/speedianet/control/src/domain/valueObject"
)

type ReadActivityRecords struct {
	RecordId           *valueObject.ActivityRecordId    `json:"recordId,omitempty"`
	RecordLevel        *valueObject.ActivityRecordLevel `json:"recordLevel,omitempty"`
	RecordCode         *valueObject.ActivityRecordCode  `json:"recordCode,omitempty"`
	OperatorAccountId  *valueObject.AccountId           `json:"operatorAccountId,omitempty"`
	OperatorIpAddress  *valueObject.IpAddress           `json:"operatorIpAddress,omitempty"`
	AccountId          *valueObject.AccountId           `json:"accountId,omitempty"`
	ContainerId        *valueObject.ContainerId         `json:"containerId,omitempty"`
	ContainerProfileId *valueObject.ContainerProfileId  `json:"containerProfileId,omitempty"`
	ContainerImageId   *valueObject.ContainerImageId    `json:"containerImageId,omitempty"`
	MappingId          *valueObject.MappingId           `json:"mappingId,omitempty"`
	ScheduledTaskId    *valueObject.ScheduledTaskId     `json:"scheduledTaskId,omitempty"`
	CreatedBeforeAt    *valueObject.UnixTime            `json:"createdBeforeAt,omitempty"`
	CreatedAfterAt     *valueObject.UnixTime            `json:"createdAfterAt,omitempty"`
}

func NewReadActivityRecords(
	recordId *valueObject.ActivityRecordId,
	recordLevel *valueObject.ActivityRecordLevel,
	recordCode *valueObject.ActivityRecordCode,
	operatorAccountId *valueObject.AccountId,
	operatorIpAddress *valueObject.IpAddress,
	accountId *valueObject.AccountId,
	containerId *valueObject.ContainerId,
	containerProfileId *valueObject.ContainerProfileId,
	containerImageId *valueObject.ContainerImageId,
	mappingId *valueObject.MappingId,
	scheduledTaskId *valueObject.ScheduledTaskId,
	createdBeforeAt *valueObject.UnixTime,
	createdAfterAt *valueObject.UnixTime,
) ReadActivityRecords {
	return ReadActivityRecords{
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
		ScheduledTaskId:    scheduledTaskId,
		CreatedBeforeAt:    createdBeforeAt,
		CreatedAfterAt:     createdAfterAt,
	}
}
