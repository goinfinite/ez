package useCase

import (
	"github.com/speedianet/control/src/domain/repository"
	"github.com/speedianet/control/src/domain/valueObject"
)

func logAction(
	serverCmdRepo repository.ServerCmdRepo,
	logPayload string,
) {
	serverCmdRepo.AddServerLog(
		valueObject.NewServerLogLevelPanic("info"),
		valueObject.NewServerLogOperationPanic("sys-install"),
		valueObject.NewServerLogPayloadPanic(logPayload),
	)
}

func SysInstall(
	sysInstallQueryRepo repository.SysInstallQueryRepo,
	sysInstallCmdRepo repository.SysInstallCmdRepo,
	serverCmdRepo repository.ServerCmdRepo,
) error {
	isInstalled := sysInstallQueryRepo.IsInstalled()
	isDataDiskMounted := sysInstallQueryRepo.IsDataDiskMounted()

	svcInstallName := valueObject.NewSvcNamePanic("sys-install-continue")

	if isInstalled && isDataDiskMounted {
		_ = serverCmdRepo.DeleteOneTimerSvc(svcInstallName)
		logAction(
			serverCmdRepo,
			"Installation succeeded. The server is now ready to be used.",
		)
		return nil
	}

	if !isInstalled {
		logAction(
			serverCmdRepo,
			"Installation started. The server will reboot a few times. "+
				"Check /var/log/control.log for the installation progress.",
		)

		err := sysInstallCmdRepo.Install()
		if err != nil {
			logAction(serverCmdRepo, err.Error())
			return err
		}

		logAction(serverCmdRepo, "Packages installed.")

		logAction(serverCmdRepo, "Disabling default softwares...")
		err = sysInstallCmdRepo.DisableDefaultSoftwares()
		if err != nil {
			logAction(serverCmdRepo, err.Error())
			return err
		}

		logAction(serverCmdRepo, "Default softwares disabled. Rebooting...")

		serverCmdRepo.AddOneTimerSvc(
			svcInstallName,
			valueObject.NewSvcCmdPanic("/var/speedia/control sys-install"),
		)

		serverCmdRepo.Reboot()
		return nil
	}

	logAction(serverCmdRepo, "Formatting data disk...")
	err := sysInstallCmdRepo.AddDataDisk()
	if err != nil {
		logAction(serverCmdRepo, err.Error())
		return err
	}

	logAction(serverCmdRepo, "Adding core services...")
	err = serverCmdRepo.AddSvc(
		valueObject.NewSvcNamePanic("control"),
		valueObject.NewSvcCmdPanic("/var/speedia/control serve"),
	)
	if err != nil {
		logAction(serverCmdRepo, err.Error())
		return err
	}

	logAction(serverCmdRepo, "Installation completed. Rebooting...")
	serverCmdRepo.Reboot()
	return nil
}
