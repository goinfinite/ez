package repository

import (
	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/entity"
)

type SecurityQueryRepo interface {
	GetEvents(getDto dto.GetSecurityEvents) ([]entity.SecurityEvent, error)
}
