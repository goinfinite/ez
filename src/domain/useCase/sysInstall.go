package useCase

import (
	"errors"

	"github.com/speedianet/sfm/src/domain/repository"
	"github.com/speedianet/sfm/src/domain/valueObject"
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

		serverCmdRepo.AddServerLog(
			valueObject.NewServerLogLevelPanic("info"),
			valueObject.NewServerLogOperationPanic("sys-install"),
			valueObject.NewServerLogPayloadPanic(
				"Packages installed, the system will now reboot "+
					"and then continue the process.",
			),
		)
		serverCmdRepo.AddOneTimerSvc("sys-install-continue", "sfm sys-install")
		serverCmdRepo.Reboot()
	}

	serverCmdRepo.DeleteOneTimerSvc("sys-install-continue")
	err := sysInstallCmdRepo.AddDataDisk()
	if err != nil {
		return err
	}

	serverCmdRepo.AddServerLog(
		valueObject.NewServerLogLevelPanic("info"),
		valueObject.NewServerLogOperationPanic("sys-install"),
		valueObject.NewServerLogPayloadPanic(
			"Data disk mounted, the system is now ready to use.",
		),
	)

	return nil
}
