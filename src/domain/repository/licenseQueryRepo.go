package repository

import (
	"github.com/speedianet/control/src/domain/entity"
)

type LicenseQueryRepo interface {
	GetStatus() (entity.LicenseStatus, error)
	GetErrorCount() (int, error)
}
