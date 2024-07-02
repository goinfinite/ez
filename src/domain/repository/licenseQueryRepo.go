package repository

import (
	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/valueObject"
)

type LicenseQueryRepo interface {
	Read() (entity.LicenseInfo, error)
	ReadIntegrityHash() (valueObject.Hash, error)
}
