package infra

import (
	"errors"
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
		"xfsprogs",
		"util-linux-tty-tools",
	}
	err := infraHelper.InstallPkgs(necessaryPkgs)
	if err != nil {
		return err
	}

	os.Symlink("/etc/pam.d/common-auth", "/etc/pam.d/system-auth")

	return nil
}

func (repo SysInstallCmdRepo) getAdditionalDisk() (string, error) {
	primaryPart, err := infraHelper.RunCmd(
		"bash",
		"-c",
		"mount | awk '/on / type/{print $1}'",
	)
	if err != nil {
		return "", errors.New("GetPrimaryPartError")
	}

	primaryDiskId, err := infraHelper.RunCmd(
		"bash",
		"-c",
		"lsblk "+primaryPart+" -n --output PKNAME",
	)
	if err != nil {
		return "", errors.New("GetPrimaryDiskIdError")
	}

	additionalDisk, err := infraHelper.RunCmd(
		"bash",
		"-c",
		"lsblk -ndp -e 7 --output KNAME | grep -v '/dev/"+primaryDiskId+"' | head -n1",
	)
	if err != nil {
		return "", errors.New("GetAddDiskError")
	}

	return additionalDisk, nil
}

func (repo SysInstallCmdRepo) AddDataDisk() error {
	addDisk, err := repo.getAdditionalDisk()
	if err != nil {
		return err
	}

	addDiskFilesystem, err := infraHelper.RunCmd(
		"lsblk",
		addDisk,
		"-n",
		"--output",
		"FSTYPE",
	)
	if err != nil {
		return errors.New("GetAddDiskFilesystemError")
	}

	if addDiskFilesystem != "" {
		return errors.New("AddDiskCannotHaveFilesystem")
	}

	_, err = infraHelper.RunCmd("mkfs.xfs", addDisk)
	if err != nil {
		return errors.New("MkfsDataDiskFailed")
	}

	_, err = infraHelper.RunCmd("mkdir", "/var/data")
	if err != nil {
		return errors.New("MkdirDataDirFailed")
	}

	_, err = infraHelper.RunCmd(
		"bash",
		"-c",
		"echo '"+addDisk+" /var/data xfs defaults,uquota 0 0' >> /etc/fstab",
	)
	if err != nil {
		return errors.New("AddDataDiskToFsTabFailed")
	}

	_, err = infraHelper.RunCmd("mount", "-a")
	if err != nil {
		return errors.New("MountDataDiskFailed")
	}

	return nil
}
