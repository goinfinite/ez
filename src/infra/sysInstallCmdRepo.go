package infra

import (
	"errors"
	"os"

	infraHelper "github.com/speedianet/control/src/infra/helper"
)

const (
	ControlMainDir = "/var/speedia"
)

type SysInstallCmdRepo struct {
}

func (repo SysInstallCmdRepo) installNginx() error {
	necessaryPkgs := []string{
		"nginx",
	}
	err := infraHelper.InstallPkgs(necessaryPkgs)
	if err != nil {
		return err
	}

	_, err = infraHelper.RunCmd("transactional-update", "apply")
	if err != nil {
		return err
	}

	nginxConf := `
user nginx;
pid /run/nginx.pid;
worker_processes auto;
worker_rlimit_nofile 65535;

load_module lib64/nginx/modules/ngx_stream_module.so;

events {
	multi_accept on;
	worker_connections 8192;
	use epoll;
}

http {
    charset utf-8;
    sendfile on;
    tcp_nopush on;
    tcp_nodelay on;
    server_tokens off;
    log_not_found off;
    types_hash_max_size 2048;
    types_hash_bucket_size 64;
    client_max_body_size 1G;

    include mime.types;
    default_type application/octet-stream;

    access_log off;
    error_log /var/log/nginx/error.log warn;

	include /var/nginx/http.d/*.conf;
}

stream {
	include /var/nginx/stream.d/*.conf;
}
`

	err = infraHelper.UpdateFile(
		"/etc/nginx/nginx.conf",
		nginxConf,
		true,
	)
	if err != nil {
		return errors.New("UpdateNginxConfFailed")
	}

	err = infraHelper.MakeDir("/var/nginx/http.d")
	if err != nil {
		return errors.New("MakeNginxHttpDirFailed")
	}

	err = infraHelper.MakeDir("/var/nginx/stream.d")
	if err != nil {
		return errors.New("MakeNginxStreamDirFailed")
	}

	_, err = infraHelper.RunCmd(
		"chown",
		"-R",
		"nginx:nginx",
		"/var/nginx",
	)
	if err != nil {
		return errors.New("ChownNginxDirFailed")
	}

	_, err = infraHelper.RunCmd(
		"systemctl",
		"enable",
		"nginx.service",
		"--now",
	)
	if err != nil {
		return errors.New("EnableNginxSvcFailed")
	}

	return nil
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
		"dmidecode",
		"whois",
		"bind-utils",
	}
	err := infraHelper.InstallPkgs(necessaryPkgs)
	if err != nil {
		return err
	}

	_, err = infraHelper.RunCmd("transactional-update", "apply")
	if err != nil {
		return err
	}

	_ = os.MkdirAll(ControlMainDir, 0755)
	_, err = infraHelper.RunCmd(
		"bash",
		"-c",
		"echo \"alias control='"+ControlMainDir+"/control'\" >> /root/.bashrc",
	)
	if err != nil {
		return errors.New("AddControlAliasFailed")
	}

	//cspell:disable
	hidepidSvc := `[Unit]
Description=Hide Pids on /proc

[Service]
Type=oneshot
ExecStart=/bin/mount -o remount,rw,nosuid,nodev,noexec,relatime,hidepid=invisible /proc

[Timer]
OnBootSec=60

[Install]
WantedBy=multi-user.target
`
	//cspell:enable

	err = infraHelper.UpdateFile(
		"/etc/systemd/system/hidepid.service",
		hidepidSvc,
		true,
	)
	if err != nil {
		return errors.New("UpdateHidepidSvcFailed")
	}

	_, err = infraHelper.RunCmd(
		"systemctl",
		"daemon-reload",
	)
	if err != nil {
		return errors.New("SystemctlDaemonReloadFailed")
	}

	_, err = infraHelper.RunCmd(
		"systemctl",
		"enable",
		"hidepid.service",
		"--now",
	)
	if err != nil {
		return errors.New("EnableHidepidSvcFailed")
	}

	err = repo.installNginx()
	if err != nil {
		return err
	}

	_ = os.MkdirAll(ControlMainDir+"/pki", 0755)
	err = infraHelper.GenSelfSignedSslCert(ControlMainDir+"/pki", "control")
	if err != nil {
		return err
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

	_, err = infraHelper.RunCmd(
		"sed",
		"-i",
		"s/GRUB_TIMEOUT=10/GRUB_TIMEOUT=5/g",
		"/etc/default/grub",
	)
	if err != nil {
		return errors.New("ReduceGrubTimeoutFailed")
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
	if err != nil || primaryPart == "" {
		return "", errors.New("GetPrimaryPartError")
	}

	primaryDiskId, err := infraHelper.RunCmd(
		"lsblk", primaryPart, "-n", "--output", "PKNAME",
	)
	if err != nil || primaryDiskId == "" {
		return "", errors.New("GetPrimaryDiskIdError")
	}

	additionalDisk, err := infraHelper.RunCmd(
		"bash",
		"-c",
		"lsblk -ndp -e 7 --output KNAME | grep -v '/dev/"+primaryDiskId+"' | head -n1",
	)
	if err != nil || additionalDisk == "" {
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
		return errors.New("AddDiskAlreadyHaveFilesystem")
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
	if err != nil || addDiskUuid == "" {
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
