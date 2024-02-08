package useCase

import (
	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/repository"
)

func GetLicenseInfo(
	licenseQueryRepo repository.LicenseQueryRepo,
) (entity.LicenseInfo, error) {
	return licenseQueryRepo.Get()
}