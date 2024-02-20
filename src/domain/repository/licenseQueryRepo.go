package repository

import (
	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/valueObject"
)

type LicenseQueryRepo interface {
	Get() (entity.LicenseInfo, error)
	GetNonceHash() (valueObject.Hash, error)
}
