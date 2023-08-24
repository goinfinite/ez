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
		"procps",
		"xfsprogs",
		"util-linux-tty-tools",
	}
	err := infraHelper.InstallPkgs(necessaryPkgs)
	if err != nil {
		return err
	}

	os.MkdirAll("/var/speedia", 0755)
	_, err = infraHelper.RunCmd(
		"bash",
		"-c",
		"echo \"alias sfm='/var/speedia/sfm'\" >> /root/.bashrc",
	)
	if err != nil {
		return errors.New("AddSfmAliasFailed")
	}

	return nil
}

func (repo SysInstallCmdRepo) DisableDefaultSoftwares() error {
	_, err := infraHelper.RunCmd(
		"sed",
		"-i",
		"s/security=selinux selinux=1/selinux=0/g",
		"/etc/default/grub",
	)
	if err != nil {
		return errors.New("DisableSelinuxFailed")
	}

	_, err = infraHelper.RunCmd("transactional-update", "grub.cfg")
	if err != nil {
		return errors.New("UpdateGrubFailed")
	}

	return nil
}

func (repo SysInstallCmdRepo) getAdditionalDisk() (string, error) {
	primaryPart, err := infraHelper.RunCmd(
		"bash",
		"-c",
		"mount | awk '/on \\/ type/{print $1}'",
	)
	if err != nil {
		return "", errors.New("GetPrimaryPartError")
	}

	primaryDiskId, err := infraHelper.RunCmd(
		"lsblk", primaryPart, "-n", "--output", "PKNAME",
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
		"lsblk", addDisk, "-n", "--output", "FSTYPE",
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

	addDiskUuid, err := infraHelper.RunCmd(
		"lsblk", addDisk, "-n", "--output", "UUID",
	)
	if err != nil {
		return errors.New("GetAddDiskUuidError")
	}

	_, err = infraHelper.RunCmd(
		"bash",
		"-c",
		"echo 'UUID="+addDiskUuid+" /var/data xfs defaults,uquota 0 0' >> /etc/fstab",
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
