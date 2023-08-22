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
		serverCmdRepo.SendServerMessage("Server is rebooting!")
		serverCmdRepo.Reboot()
	}

	serverCmdRepo.SendServerMessage(
		"[Elevator Music ♬] SFM is still installing... " +
			"please await and don't interact with the server.",
	)
	serverCmdRepo.DeleteOneTimerSvc("sys-install-continue")
	err := sysInstallCmdRepo.AddDataDisk()
	if err != nil {
		return err
	}

	successMessage := "SFM has been installed and is now ready to be used."
	serverCmdRepo.AddServerLog(
		valueObject.NewServerLogLevelPanic("info"),
		valueObject.NewServerLogOperationPanic("sys-install"),
		valueObject.NewServerLogPayloadPanic(successMessage),
	)
	serverCmdRepo.SendServerMessage(successMessage)

	return nil
}
