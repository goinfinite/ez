package repository

import (
	"github.com/speedianet/control/src/domain/entity"
)

type O11yQueryRepo interface {
	ReadOverview() (entity.O11yOverview, error)
}
