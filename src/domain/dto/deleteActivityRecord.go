package dto

import (
	"github.com/speedianet/control/src/domain/valueObject"
)

type DeleteActivityRecords struct {
	RecordId          *valueObject.ActivityRecordId          `json:"recordId,omitempty"`
	RecordLevel       *valueObject.ActivityRecordLevel       `json:"recordLevel,omitempty"`
	RecordCode        *valueObject.ActivityRecordCode        `json:"recordCode,omitempty"`
	AffectedResources []valueObject.SystemResourceIdentifier `json:"affectedResources,omitempty"`
	OperatorAccountId *valueObject.AccountId                 `json:"operatorAccountId,omitempty"`
	OperatorIpAddress *valueObject.IpAddress                 `json:"operatorIpAddress,omitempty"`
	CreatedBeforeAt   *valueObject.UnixTime                  `json:"createdBeforeAt,omitempty"`
	CreatedAfterAt    *valueObject.UnixTime                  `json:"createdAfterAt,omitempty"`
}

func NewDeleteActivityRecords(
	recordId *valueObject.ActivityRecordId,
	recordLevel *valueObject.ActivityRecordLevel,
	recordCode *valueObject.ActivityRecordCode,
	affectedResources []valueObject.SystemResourceIdentifier,
	operatorAccountId *valueObject.AccountId,
	operatorIpAddress *valueObject.IpAddress,
	createdBeforeAt *valueObject.UnixTime,
	createdAfterAt *valueObject.UnixTime,
) DeleteActivityRecords {
	return DeleteActivityRecords{
		RecordId:          recordId,
		RecordLevel:       recordLevel,
		RecordCode:        recordCode,
		AffectedResources: affectedResources,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
		CreatedBeforeAt:   createdBeforeAt,
		CreatedAfterAt:    createdAfterAt,
	}
}
