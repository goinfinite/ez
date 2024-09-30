package infra

import (
	"os"

	infraEnvs "github.com/goinfinite/ez/src/infra/envs"
	infraHelper "github.com/goinfinite/ez/src/infra/helper"
)

type SysInstallQueryRepo struct {
}

func (repo SysInstallQueryRepo) IsInstalled() bool {
	out, err := infraHelper.RunCmd("grep", "alias ez=", "/root/.bashrc")
	if err != nil || out == "" {
		return false
	}

	return true
}

func (repo SysInstallQueryRepo) IsDataDiskMounted() bool {
	_, err := os.Stat(infraEnvs.UserDataDirectory)
	return err == nil
}
