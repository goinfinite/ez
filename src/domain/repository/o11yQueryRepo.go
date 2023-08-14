package repository

import (
	"github.com/speedianet/sfm/src/domain/entity"
)

type O11yQueryRepo interface {
	GetOverview() (entity.O11yOverview, error)
}
