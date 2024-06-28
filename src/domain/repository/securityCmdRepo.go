package repository

import (
	"github.com/speedianet/control/src/domain/dto"
)

type SecurityCmdRepo interface {
	CreateEvent(createDto dto.CreateSecurityEvent) error
	DeleteEvents(deleteDto dto.DeleteSecurityEvents) error
}
