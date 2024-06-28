package dto

import (
	"github.com/speedianet/control/src/domain/valueObject"
)

type DeleteSecurityEvents struct {
	Type      *valueObject.SecurityEventType `json:"eventType"`
	IpAddress *valueObject.IpAddress         `json:"ipAddress"`
	AccountId *valueObject.AccountId         `json:"accountId"`
	CreatedAt *valueObject.UnixTime          `json:"createdAt"`
}

func NewDeleteSecurityEvents(
	eventType *valueObject.SecurityEventType,
	ipAddress *valueObject.IpAddress,
	accountId *valueObject.AccountId,
	createdAt *valueObject.UnixTime,
) DeleteSecurityEvents {
	return DeleteSecurityEvents{
		Type:      eventType,
		IpAddress: ipAddress,
		AccountId: accountId,
		CreatedAt: createdAt,
	}
}
