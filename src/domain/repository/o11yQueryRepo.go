package repository

import (
	"github.com/goinfinite/fleet/src/domain/entity"
)

type O11yQueryRepo interface {
	GetOverview() (entity.O11yOverview, error)
}
