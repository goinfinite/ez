package repository

import (
	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/entity"
)

type SecurityQueryRepo interface {
	ReadEvents(readDto dto.ReadSecurityEvents) ([]entity.SecurityEvent, error)
}
