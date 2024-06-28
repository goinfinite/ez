package dto

import (
	"github.com/speedianet/control/src/domain/valueObject"
)

type ReadSecurityEvents struct {
	Type      *valueObject.SecurityEventType `json:"eventType"`
	IpAddress *valueObject.IpAddress         `json:"ipAddress"`
	AccountId *valueObject.AccountId         `json:"accountId"`
	CreatedAt *valueObject.UnixTime          `json:"createdAt"`
}

func NewReadSecurityEvents(
	eventType *valueObject.SecurityEventType,
	ipAddress *valueObject.IpAddress,
	accountId *valueObject.AccountId,
	createdAt *valueObject.UnixTime,
) ReadSecurityEvents {
	return ReadSecurityEvents{
		Type:      eventType,
		IpAddress: ipAddress,
		AccountId: accountId,
		CreatedAt: createdAt,
	}
}
