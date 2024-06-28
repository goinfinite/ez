package dto

import (
	"github.com/speedianet/control/src/domain/valueObject"
)

type CreateSecurityEvent struct {
	Type      valueObject.SecurityEventType     `json:"eventType"`
	Details   *valueObject.SecurityEventDetails `json:"details"`
	IpAddress *valueObject.IpAddress            `json:"ipAddress"`
	AccountId *valueObject.AccountId            `json:"accountId"`
}

func NewCreateSecurityEvent(
	eventType valueObject.SecurityEventType,
	details *valueObject.SecurityEventDetails,
	ipAddress *valueObject.IpAddress,
	accountId *valueObject.AccountId,
) CreateSecurityEvent {
	return CreateSecurityEvent{
		Type:      eventType,
		Details:   details,
		IpAddress: ipAddress,
		AccountId: accountId,
	}
}
