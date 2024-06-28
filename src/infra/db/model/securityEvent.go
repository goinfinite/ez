package dbModel

import (
	"time"

	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/valueObject"
)

type SecurityEvent struct {
	ID        uint   `gorm:"primarykey"`
	Type      string `gorm:"not null"`
	Details   *string
	IpAddress *string
	AccountId *uint
	CreatedAt time.Time `gorm:"not null"`
}

func (SecurityEvent) TableName() string {
	return "security_events"
}

func NewSecurityEvent(
	id uint,
	eventType string,
	details *string,
	ipAddress *string,
	accountId *uint,
) SecurityEvent {
	model := SecurityEvent{
		Type: eventType,
	}

	if id != 0 {
		model.ID = id
	}

	if details != nil {
		model.Details = details
	}

	if ipAddress != nil {
		model.IpAddress = ipAddress
	}

	if accountId != nil {
		model.AccountId = accountId
	}

	return model
}

func (model SecurityEvent) ToEntity() (eventEntity entity.SecurityEvent, err error) {
	id, err := valueObject.NewSecurityEventId(model.ID)
	if err != nil {
		return eventEntity, err
	}

	eventType, err := valueObject.NewSecurityEventType(model.Type)
	if err != nil {
		return eventEntity, err
	}

	var detailsPtr *valueObject.SecurityEventDetails
	if model.Details != nil {
		details, err := valueObject.NewSecurityEventDetails(*model.Details)
		if err != nil {
			return eventEntity, err
		}
		detailsPtr = &details
	}

	var ipAddressPtr *valueObject.IpAddress
	if model.IpAddress != nil {
		ipAddress, err := valueObject.NewIpAddress(*model.IpAddress)
		if err != nil {
			return eventEntity, err
		}
		ipAddressPtr = &ipAddress
	}

	var accountIdPtr *valueObject.AccountId
	if model.AccountId != nil {
		accountId, err := valueObject.NewAccountId(*model.AccountId)
		if err != nil {
			return eventEntity, err
		}
		accountIdPtr = &accountId
	}

	createdAt := valueObject.NewUnixTimeWithGoTime(model.CreatedAt)

	return entity.NewSecurityEvent(
		id, eventType, detailsPtr, ipAddressPtr, accountIdPtr, createdAt,
	), nil
}
