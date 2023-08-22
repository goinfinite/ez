package useCase

import (
	"errors"

	"github.com/speedianet/sfm/src/domain/repository"
)

func SysInstall(
	sysInstallQueryRepo repository.SysInstallQueryRepo,
	sysInstallCmdRepo repository.SysInstallCmdRepo,
) error {
	isInstalled := sysInstallQueryRepo.IsInstalled()
	isDataDiskMounted := sysInstallQueryRepo.IsDataDiskMounted()

	if isInstalled && isDataDiskMounted {
		return errors.New("SysAlreadyInstalled")
	}

	if !isInstalled {
		err := sysInstallCmdRepo.Install()
		if err != nil {
			return err
		}
	}

	return sysInstallCmdRepo.AddDataDisk()
}
