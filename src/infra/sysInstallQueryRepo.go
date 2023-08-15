package infra

import infraHelper "github.com/speedianet/sfm/src/infra/helper"

type SysInstallQueryRepo struct {
}

func (repo SysInstallQueryRepo) IsInstalled() bool {
	_, err := infraHelper.GetFilePathWithMatch("/etc/pam.d/", "system-auth")
	return err == nil
}
