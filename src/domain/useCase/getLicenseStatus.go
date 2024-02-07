package useCase

import (
	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/repository"
)

func GetLicenseStatus(
	licenseQueryRepo repository.LicenseQueryRepo,
) (entity.LicenseStatus, error) {
	return licenseQueryRepo.GetStatus()
}
