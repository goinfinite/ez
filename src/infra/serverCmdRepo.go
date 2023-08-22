package infra

import (
	"encoding/json"
	"errors"
	"log"
	"os"

	"github.com/speedianet/sfm/src/domain/entity"
	"github.com/speedianet/sfm/src/domain/valueObject"
	infraHelper "github.com/speedianet/sfm/src/infra/helper"
)

type ServerCmdRepo struct {
}

func (repo ServerCmdRepo) Reboot() error {
	infraHelper.RunCmd("systemctl", "reboot")
	return nil
}

func (repo ServerCmdRepo) AddOneTimerSvc(name string, cmd string) error {
	svcFilePath := "/etc/systemd/system/" + name + ".service"
	svcContent := `[Unit]
Description=` + name + `
After=network.target

[Service]
Type=oneshot
ExecStart=` + cmd + `
RemainAfterExit=yes
`

	err := infraHelper.UpdateFile(svcFilePath, svcContent, true)
	if err != nil {
		return errors.New("AddOneTimerSvcFailed")
	}

	svcTimerFilePath := "/etc/systemd/system/" + name + ".timer"
	svcTimerContent := `[Unit]
Description=` + name + `

[Timer]
OnBootSec=10s
Unit=` + name + `.service

[Install]
WantedBy=multi-user.target
`

	err = infraHelper.UpdateFile(svcTimerFilePath, svcTimerContent, true)
	if err != nil {
		return errors.New("AddOneTimerSvcTimerFailed")
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
		name+".service",
	)
	if err != nil {
		return errors.New("SystemctlEnableFailed")
	}

	return nil
}

func (repo ServerCmdRepo) DeleteOneTimerSvc(name string) error {
	infraHelper.RunCmd("systemctl", "stop", name+".service")
	infraHelper.RunCmd("systemctl", "disable", name+".service")

	err := os.Remove("/etc/systemd/system/" + name + ".service")
	if err != nil {
		return errors.New("RemoveSvcFailed")
	}

	err = os.Remove("/etc/systemd/system/" + name + ".timer")
	if err != nil {
		return errors.New("RemoveSvcTimerFailed")
	}

	_, err = infraHelper.RunCmd(
		"systemctl",
		"daemon-reload",
	)
	if err != nil {
		return errors.New("SystemctlDaemonReloadFailed")
	}

	return nil
}

func (repo ServerCmdRepo) AddServerLog(
	level valueObject.ServerLogLevel,
	operation valueObject.ServerLogOperation,
	payload valueObject.ServerLogPayload,
) {
	logEntity := entity.NewServerLog(level, operation, payload)
	logFilePath := "/var/log/sfm.log"
	jsonLogEntry, err := json.Marshal(logEntity)
	if err != nil {
		return
	}
	logContent := string(jsonLogEntry) + "\n"

	logFile, err := os.OpenFile(
		logFilePath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)
	if err != nil {
		return
	}
	defer logFile.Close()

	logFile.WriteString(logContent)
	log.Print(logContent)
}
