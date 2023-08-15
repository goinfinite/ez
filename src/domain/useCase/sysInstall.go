package useCase

import (
	"errors"

	"github.com/speedianet/sfm/src/domain/repository"
)

func SysInstall(
	sysInstallQueryRepo repository.SysInstallQueryRepo,
	sysInstallCmdRepo repository.SysInstallCmdRepo,
) error {
	if sysInstallQueryRepo.IsInstalled() {
		return errors.New("SysAlreadyInstalled")
	}

	return sysInstallCmdRepo.Install()
}
