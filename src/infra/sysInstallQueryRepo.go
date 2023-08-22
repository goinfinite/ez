package infra

import (
	"os"

	infraHelper "github.com/speedianet/sfm/src/infra/helper"
)

type SysInstallQueryRepo struct {
}

func (repo SysInstallQueryRepo) IsInstalled() bool {
	_, err := infraHelper.GetFilePathWithMatch("/etc/pam.d/", "system-auth")
	return err == nil
}

func (repo SysInstallQueryRepo) IsDataDiskMounted() bool {
	_, err := os.Stat("/var/data")
	return err == nil
}
