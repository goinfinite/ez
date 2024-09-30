package repository

import (
	"github.com/goinfinite/ez/src/domain/entity"
	"github.com/goinfinite/ez/src/domain/valueObject"
)

type LicenseQueryRepo interface {
	Read() (entity.LicenseInfo, error)
	ReadIntegrityHash() (valueObject.Hash, error)
}
