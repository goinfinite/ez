package useCase

import (
	"log/slog"

	"github.com/speedianet/control/src/domain/repository"
	"github.com/speedianet/control/src/domain/valueObject"
)

func SysInstall(
	sysInstallQueryRepo repository.SysInstallQueryRepo,
	sysInstallCmdRepo repository.SysInstallCmdRepo,
	serverCmdRepo repository.ServerCmdRepo,
) error {
	isInstalled := sysInstallQueryRepo.IsInstalled()
	isDataDiskMounted := sysInstallQueryRepo.IsDataDiskMounted()

	svcInstallName := valueObject.NewServiceNamePanic("sys-install-continue")

	if isInstalled && isDataDiskMounted {
		_ = serverCmdRepo.DeleteOneTimerSvc(svcInstallName)
		slog.Info("Installation succeeded. The server is now ready to be used.")
		return nil
	}

	if !isInstalled {
		slog.Info("Installation started. The server will reboot a few times. " +
			"Check /var/log/control.log for the installation progress.")

		err := sysInstallCmdRepo.Install()
		if err != nil {
			slog.Error(err.Error())
			return err
		}

		slog.Info("Packages installed. Disabling default softwares...")
		err = sysInstallCmdRepo.DisableDefaultSoftwares()
		if err != nil {
			slog.Error(err.Error())
			return err
		}

		slog.Info("Default softwares disabled. Rebooting...")

		serverCmdRepo.AddOneTimerSvc(
			svcInstallName,
			valueObject.NewSvcCmdPanic("/var/speedia/control sys-install"),
		)

		serverCmdRepo.Reboot()
		return nil
	}

	slog.Info("Formatting data disk...")
	err := sysInstallCmdRepo.AddDataDisk()
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	slog.Info("Adding core services...")
	err = serverCmdRepo.AddSvc(
		valueObject.NewServiceNamePanic("control"),
		valueObject.NewSvcCmdPanic("/var/speedia/control serve"),
	)
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	slog.Info("Installation completed. Rebooting...")
	serverCmdRepo.Reboot()
	return nil
}
