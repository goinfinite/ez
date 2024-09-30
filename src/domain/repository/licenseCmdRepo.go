package repository

import (
	"github.com/goinfinite/ez/src/domain/entity"
	"github.com/goinfinite/ez/src/domain/valueObject"
)

type LicenseCmdRepo interface {
	GenerateIntegrityHash(licenseInfo entity.LicenseInfo) (valueObject.Hash, error)
	GenerateNonceHash() (valueObject.Hash, error)
	Refresh() error
	UpdateStatus(status valueObject.LicenseStatus) error
	IncrementErrorCount() error
	ResetErrorCount() error
}
