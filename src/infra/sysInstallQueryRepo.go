package infra

import (
	"os"

	infraHelper "github.com/speedianet/sfm/src/infra/helper"
)

type SysInstallQueryRepo struct {
}

func (repo SysInstallQueryRepo) IsInstalled() bool {
	filePath, err := infraHelper.GetFilePathWithMatch("/usr/bin", "sfm")
	if err != nil || filePath == "" {
		return false
	}

	return true
}

func (repo SysInstallQueryRepo) IsDataDiskMounted() bool {
	_, err := os.Stat("/var/data")
	return err == nil
}
