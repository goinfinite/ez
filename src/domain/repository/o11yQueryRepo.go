package repository

import (
	"github.com/goinfinite/ez/src/domain/entity"
)

type O11yQueryRepo interface {
	ReadOverview() (entity.O11yOverview, error)
}
