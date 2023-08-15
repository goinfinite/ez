package infra

import (
	"os"

	infraHelper "github.com/speedianet/sfm/src/infra/helper"
)

type SysInstallCmdRepo struct {
}

func (repo SysInstallCmdRepo) Install() error {
	necessaryPkgs := []string{
		"git",
		"wget",
		"curl",
		"cyrus-sasl",
		"pam-devel",
		"gcc",
		"make",
		"tar",
		"procps",
	}
	err := infraHelper.InstallPkgs(necessaryPkgs)
	if err != nil {
		return err
	}

	os.Symlink("/etc/pam.d/common-auth", "/etc/pam.d/system-auth")

	infraHelper.RunCmd("transactional-update", "reboot")

	return nil
}
