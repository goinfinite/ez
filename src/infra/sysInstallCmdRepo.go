package infra

import (
	"errors"
	"math/rand"
	"os"
	"strconv"
	"time"

	infraHelper "github.com/speedianet/control/src/infra/helper"
)

const (
	ControlMainDir = "/var/speedia"
)

type SysInstallCmdRepo struct {
}

func (repo SysInstallCmdRepo) installNginx() error {
	necessaryPkgs := []string{"nginx"}
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

	err = infraHelper.UpdateFile("/etc/nginx/nginx.conf", nginxConf, true)
	if err != nil {
		return errors.New("UpdateNginxConfFailed: " + err.Error())
	}

	err = infraHelper.MakeDir("/var/nginx/http.d")
	if err != nil {
		return errors.New("MakeNginxHttpDirFailed: " + err.Error())
	}

	err = infraHelper.MakeDir("/var/nginx/stream.d")
	if err != nil {
		return errors.New("MakeNginxStreamDirFailed: " + err.Error())
	}

	_, err = infraHelper.RunCmd("chown", "-R", "nginx:nginx", "/var/nginx")
	if err != nil {
		return errors.New("ChownNginxDirFailed: " + err.Error())
	}

	_, err = infraHelper.RunCmd("systemctl", "enable", "nginx.service", "--now")
	if err != nil {
		return errors.New("EnableNginxSvcFailed: " + err.Error())
	}

	return nil
}

func (repo SysInstallCmdRepo) Install() error {
	//cspell:disable
	necessaryPkgs := []string{
		"git", "wget", "curl", "cyrus-sasl", "procps", "xfsprogs",
		"util-linux-tty-tools", "dmidecode", "whois", "bind-utils", "jq",
	}
	//cspell:enable
	err := infraHelper.InstallPkgs(necessaryPkgs)
	if err != nil {
		return err
	}

	_, err = infraHelper.RunCmd("transactional-update", "apply")
	if err != nil {
		return err
	}

	_ = os.MkdirAll(ControlMainDir, 0755)
	_, err = infraHelper.RunCmdWithSubShell(
		"echo \"alias control='" + ControlMainDir + "/control'\" >> /root/.bashrc",
	)
	if err != nil {
		return errors.New("AddControlAliasFailed: " + err.Error())
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

	err = infraHelper.UpdateFile("/etc/systemd/system/hidepid.service", hidepidSvc, true)
	if err != nil {
		return errors.New("UpdateHidepidSvcFailed: " + err.Error())
	}

	//cspell:disable
	cgroupControllersSvc := `[Unit]
Description=CGroup Controllers
After=network.target

[Service]
Type=oneshot
ExecStart=/bin/bash -c 'echo "+cpu +cpuset +io +memory +pids" > /sys/fs/cgroup/cgroup.subtree_control'

[Install]
WantedBy=multi-user.target
`
	//cspell:enable

	err = infraHelper.UpdateFile(
		"/etc/systemd/system/cgroup-controllers.service",
		cgroupControllersSvc,
		true,
	)
	if err != nil {
		return errors.New("UpdateCgroupControllersSvcFailed: " + err.Error())
	}

	err = infraHelper.MakeDir("/etc/systemd/system/user@.service.d")
	if err != nil {
		return errors.New("MakeSystemdUserSliceDirFailed: " + err.Error())
	}

	//cspell:disable
	cgroupDelegateConf := `[Service]
Delegate=cpu cpuset io memory pids
`
	//cspell:enable

	err = infraHelper.UpdateFile(
		"/etc/systemd/system/user@.service.d/delegate.conf",
		cgroupDelegateConf,
		true,
	)
	if err != nil {
		return errors.New("UpdateCgroupDelegateConfFailed: " + err.Error())
	}

	err = infraHelper.MakeDir("/etc/systemd/system/transactional-update.timer.d")
	if err != nil {
		return errors.New("MakeSystemdTransUpdateTimerDirFailed: " + err.Error())
	}

	randomMorningHour := strconv.Itoa(rand.Intn(7))
	//cspell:disable
	updateTimerOverride := `[Timer]
OnCalendar=
OnCalendar=Sat *-*-* 0` + randomMorningHour + `:00:00
`
	//cspell:enable

	err = infraHelper.UpdateFile(
		"/etc/systemd/system/transactional-update.timer.d/timer.conf",
		updateTimerOverride,
		true,
	)
	if err != nil {
		return errors.New("UpdateTransUpdateTimerOverrideFailed: " + err.Error())
	}

	_, err = infraHelper.RunCmd("systemctl", "daemon-reload")
	if err != nil {
		return errors.New("SystemctlDaemonReloadFailed: " + err.Error())
	}

	_, err = infraHelper.RunCmd("systemctl", "enable", "hidepid.service", "--now")
	if err != nil {
		return errors.New("EnableHidepidSvcFailed: " + err.Error())
	}

	_, err = infraHelper.RunCmd(
		"systemctl", "enable", "cgroup-controllers.service", "--now",
	)
	if err != nil {
		return errors.New("EnableCgroupControllersSvcFailed: " + err.Error())
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
		"s/security=selinux selinux=1/selinux=0 systemd.unified_cgroup_hierarchy=1/g",
		"/etc/default/grub",
	)
	if err != nil {
		return errors.New("DisableSelinuxFailed: " + err.Error())
	}

	_, err = infraHelper.RunCmd(
		"sed",
		"-i",
		"s/GRUB_TIMEOUT=10/GRUB_TIMEOUT=5/g",
		"/etc/default/grub",
	)
	if err != nil {
		return errors.New("ReduceGrubTimeoutFailed: " + err.Error())
	}

	_, err = infraHelper.RunCmd("transactional-update", "grub.cfg")
	if err != nil {
		return errors.New("UpdateGrubFailed: " + err.Error())
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
		return "", errors.New("GetPrimaryPartError: " + err.Error())
	}

	if primaryPart == "" {
		return "", errors.New("PrimaryPartNotFound")
	}

	primaryDiskId, err := infraHelper.RunCmd(
		"lsblk", primaryPart, "-n", "--output", "PKNAME",
	)
	if err != nil {
		return "", errors.New("GetPrimaryDiskIdError: " + err.Error())
	}

	if primaryDiskId == "" {
		return "", errors.New("PrimaryDiskIdNotFound")
	}

	additionalDisk, err := infraHelper.RunCmd(
		"bash",
		"-c",
		"lsblk -ndp -e 7 --output KNAME | grep -v '/dev/"+primaryDiskId+"' | head -n1",
	)
	if err != nil {
		return "", errors.New("GetAddDiskError: " + err.Error())
	}

	if additionalDisk == "" {
		return "", errors.New("AddDiskNotFound")
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
		return errors.New("GetAddDiskFilesystemError: " + err.Error())
	}

	if addDiskFilesystem != "" {
		return errors.New("AddDiskAlreadyHaveFilesystem")
	}

	_, err = infraHelper.RunCmd("mkfs.xfs", addDisk)
	if err != nil {
		return errors.New("MkfsDataDiskFailed: " + err.Error())
	}

	time.Sleep(5 * time.Second)

	addDiskUuid, err := infraHelper.RunCmd(
		"lsblk", addDisk, "-n", "--output", "UUID",
	)
	if err != nil {
		return errors.New("GetAddDiskUuidError: " + err.Error())
	}

	if addDiskUuid == "" {
		return errors.New("AddDiskUuidEmpty")
	}

	err = infraHelper.MakeDir("/var/data")
	if err != nil {
		return errors.New("MkdirDataDirFailed: " + err.Error())
	}

	_, err = infraHelper.RunCmdWithSubShell(
		"echo 'UUID=" + addDiskUuid + " /var/data xfs defaults,uquota,prjquota 0 0' >> /etc/fstab",
	)
	if err != nil {
		return errors.New("AddDataDiskToFsTabFailed: " + err.Error())
	}

	_, err = infraHelper.RunCmd("mount", "-a")
	if err != nil {
		return errors.New("MountDataDiskFailed: " + err.Error())
	}

	_, err = infraHelper.RunCmd("transactional-update", "apply")
	if err != nil {
		return err
	}

	return nil
}
