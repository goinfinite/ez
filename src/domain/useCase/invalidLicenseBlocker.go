package useCase

import (
	"errors"

	"github.com/speedianet/control/src/domain/repository"
)

func InvalidLicenseBlocker(
	licenseQueryRepo repository.LicenseQueryRepo,
	containerQueryRepo repository.ContainerQueryRepo,
	containerCmdRepo repository.ContainerCmdRepo,
) error {
	licenseInfo, err := licenseQueryRepo.Read()
	if err != nil {
		return errors.New("GetLicenseStatusError: " + err.Error())
	}

	licenseStatusStr := licenseInfo.Status.String()
	switch licenseStatusStr {
	case "ACTIVE":
		return nil
	case "SUSPENDED":
		return errors.New("LicenseSuspendedNoActionAllowed")
	case "REVOKED":
		err = StopAllContainers(containerQueryRepo, containerCmdRepo)
		if err != nil {
			return errors.New("StopAllContainersError: " + err.Error())
		}
		return errors.New("LicenseRevokedSystemDisabled")
	}

	return nil
}
