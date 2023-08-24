package useCase

import (
	"github.com/speedianet/sfm/src/domain/repository"
	"github.com/speedianet/sfm/src/domain/valueObject"
)

func logAction(
	serverCmdRepo repository.ServerCmdRepo,
	logPayload valueObject.ServerLogPayload,
) {
	serverCmdRepo.AddServerLog(
		valueObject.NewServerLogLevelPanic("info"),
		valueObject.NewServerLogOperationPanic("sys-install"),
		logPayload,
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
		serverCmdRepo.DeleteOneTimerSvc(svcInstallName)
		logAction(
			serverCmdRepo,
			valueObject.NewServerLogPayloadPanic(
				"The installation succeeded. The server is ready to use.",
			),
		)
		return nil
	}

	if !isInstalled {
		logAction(
			serverCmdRepo,
			valueObject.NewServerLogPayloadPanic(
				"Installation started, the server will reboot a few times. "+
					"Check /var/log/sfm.log for the installation progress.",
			),
		)

		err := sysInstallCmdRepo.Install()
		if err != nil {
			return err
		}

		serverCmdRepo.AddOneTimerSvc(
			svcInstallName,
			valueObject.NewSvcCmdPanic("/var/speedia/sfm sys-install"),
		)

		logAction(
			serverCmdRepo,
			valueObject.NewServerLogPayloadPanic(
				"Packages installed. Rebooting...",
			),
		)
		serverCmdRepo.Reboot()
	}

	logAction(
		serverCmdRepo,
		valueObject.NewServerLogPayloadPanic(
			"Formatting data disk...",
		),
	)
	err := sysInstallCmdRepo.AddDataDisk()
	if err != nil {
		return err
	}

	logAction(
		serverCmdRepo,
		valueObject.NewServerLogPayloadPanic(
			"Adding core services...",
		),
	)
	err = serverCmdRepo.AddSvc(
		valueObject.NewSvcNamePanic("sfm"),
		valueObject.NewSvcCmdPanic("/var/speedia/sfm serve"),
	)
	if err != nil {
		return err
	}

	logAction(
		serverCmdRepo,
		valueObject.NewServerLogPayloadPanic(
			"Disabling default softwares...",
		),
	)
	err = sysInstallCmdRepo.DisableDefaultSoftwares()
	if err != nil {
		return err
	}

	logAction(
		serverCmdRepo,
		valueObject.NewServerLogPayloadPanic(
			"Installation completed! Rebooting...",
		),
	)
	return nil
}
