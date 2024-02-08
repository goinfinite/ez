package repository

import "github.com/speedianet/control/src/domain/valueObject"

type LicenseCmdRepo interface {
	RefreshStatus() error
	UpdateStatus(status valueObject.LicenseStatus) error
	IncrementErrorCount() error
	ResetErrorCount() error
}
