package dto

import (
	"github.com/speedianet/control/src/domain/valueObject"
)

type GetSecurityEvents struct {
	Type      *valueObject.SecurityEventType `json:"eventType"`
	IpAddress *valueObject.IpAddress         `json:"ipAddress"`
	AccountId *valueObject.AccountId         `json:"accountId"`
	CreatedAt *valueObject.UnixTime          `json:"createdAt"`
}

func NewGetSecurityEvents(
	eventType *valueObject.SecurityEventType,
	ipAddress *valueObject.IpAddress,
	accountId *valueObject.AccountId,
	createdAt *valueObject.UnixTime,
) GetSecurityEvents {
	return GetSecurityEvents{
		Type:      eventType,
		IpAddress: ipAddress,
		AccountId: accountId,
		CreatedAt: createdAt,
	}
}
