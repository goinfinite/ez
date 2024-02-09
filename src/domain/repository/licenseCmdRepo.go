package repository

import "github.com/speedianet/control/src/domain/valueObject"

type LicenseCmdRepo interface {
	Refresh() error
	UpdateStatus(status valueObject.LicenseStatus) error
	IncrementErrorCount() error
	ResetErrorCount() error
}
