package entity

import "github.com/speedianet/control/src/domain/valueObject"

type SecurityEvent struct {
	Id        valueObject.SecurityEventId       `json:"id"`
	Type      valueObject.SecurityEventType     `json:"eventType"`
	Details   *valueObject.SecurityEventDetails `json:"details"`
	IpAddress *valueObject.IpAddress            `json:"ipAddress"`
	AccountId *valueObject.AccountId            `json:"accountId"`
	CreatedAt valueObject.UnixTime              `json:"createdAt"`
}

func NewSecurityEvent(
	id valueObject.SecurityEventId,
	eventType valueObject.SecurityEventType,
	details *valueObject.SecurityEventDetails,
	ipAddress *valueObject.IpAddress,
	accountId *valueObject.AccountId,
	createdAt valueObject.UnixTime,
) SecurityEvent {
	return SecurityEvent{
		Id:        id,
		Type:      eventType,
		Details:   details,
		IpAddress: ipAddress,
		AccountId: accountId,
		CreatedAt: createdAt,
	}
}
