package infra

import (
	"os"

	infraHelper "github.com/speedianet/control/src/infra/helper"
)

type SysInstallQueryRepo struct {
}

func (repo SysInstallQueryRepo) IsInstalled() bool {
	out, err := infraHelper.RunCmd("grep", "alias control=", "/root/.bashrc")
	if err != nil || out == "" {
		return false
	}

	return true
}

func (repo SysInstallQueryRepo) IsDataDiskMounted() bool {
	_, err := os.Stat("/var/data")
	return err == nil
}
