package useCase

import (
	"errors"

	"github.com/speedianet/sfm/src/domain/repository"
)

func SysInstall(
	sysInstallQueryRepo repository.SysInstallQueryRepo,
	sysInstallCmdRepo repository.SysInstallCmdRepo,
	serverCmdRepo repository.ServerCmdRepo,
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

		serverCmdRepo.AddOneTimerSvc("sys-install-continue", "sfm sys-install")
		serverCmdRepo.Reboot()
	}

	serverCmdRepo.DeleteOneTimerSvc("sys-install-continue")
	return sysInstallCmdRepo.AddDataDisk()
}
