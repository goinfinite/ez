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

	serverCmdRepo.AddServerLog(
		valueObject.NewServerLogLevelPanic("info"),
		valueObject.NewServerLogOperationPanic("sys-install"),
		valueObject.NewServerLogPayloadPanic("System installation started."),
	)

	svcInstallName := valueObject.NewSvcNamePanic("sys-install-continue")

	if !isInstalled {
		err := sysInstallCmdRepo.Install()
		if err != nil {
			return err
		}

		err = sysInstallCmdRepo.DisableDefaultSoftwares()
		if err != nil {
			return err
		}

		serverCmdRepo.AddServerLog(
			valueObject.NewServerLogLevelPanic("info"),
			valueObject.NewServerLogOperationPanic("sys-install"),
			valueObject.NewServerLogPayloadPanic(
				"Packages installed. The server will now reboot "+
					"and then continue the process automatically.",
			),
		)
		serverCmdRepo.AddOneTimerSvc(
			svcInstallName,
			valueObject.NewSvcCmdPanic("/var/speedia/sfm sys-install"),
		)
		serverCmdRepo.SendServerMessage("Server is rebooting!")
		serverCmdRepo.Reboot()
	}

	serverCmdRepo.SendServerMessage(
		"[elevator music ♬] SFM is still installing... " +
			"please await and don't interact with the server.",
	)
	serverCmdRepo.DeleteOneTimerSvc(svcInstallName)
	err := sysInstallCmdRepo.AddDataDisk()
	if err != nil {
		return err
	}

	err = serverCmdRepo.AddSvc(
		valueObject.NewSvcNamePanic("sfm"),
		valueObject.NewSvcCmdPanic("/var/speedia/sfm serve"),
	)
	if err != nil {
		return err
	}

	successMessage := "SFM has been installed and is running. " +
		"The server is now ready to be used."
	serverCmdRepo.AddServerLog(
		valueObject.NewServerLogLevelPanic("info"),
		valueObject.NewServerLogOperationPanic("sys-install"),
		valueObject.NewServerLogPayloadPanic(successMessage),
	)
	serverCmdRepo.SendServerMessage(successMessage)

	return nil
}
