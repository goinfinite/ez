package repository

import (
	"github.com/speedianet/control/src/domain/entity"
)

type LicenseQueryRepo interface {
	Get() (entity.LicenseInfo, error)
	GetErrorCount() (int, error)
}
