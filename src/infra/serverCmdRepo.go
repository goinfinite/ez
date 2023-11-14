package infra

import (
	"encoding/json"
	"errors"
	"log"
	"os"

	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/valueObject"
	infraHelper "github.com/speedianet/control/src/infra/helper"
)

type ServerCmdRepo struct {
}

func (repo ServerCmdRepo) Reboot() error {
	_, _ = infraHelper.RunCmd("systemctl", "reboot")
	os.Exit(0)
	return nil
}

func (repo ServerCmdRepo) AddSvc(
	name valueObject.SvcName,
	cmd valueObject.SvcCmd,
) error {
	nameStr := name.String()
	cmdStr := cmd.String()
	svcFilePath := "/etc/systemd/system/" + nameStr + ".service"
	svcContent := `[Unit]
Description=` + nameStr + `
After=network.target

[Service]
User=root
WorkingDirectory=/var/speedia
ExecStart=` + cmdStr + `
Restart=always
StandardOutput=journal
StandardError=journal
SyslogIdentifier=` + nameStr + `
RestartSec=15

[Install]
WantedBy=multi-user.target
`

	err := infraHelper.UpdateFile(svcFilePath, svcContent, true)
	if err != nil {
		return errors.New("AddSvcFailed")
	}
	err = os.Chmod(svcFilePath, 0644)
	if err != nil {
		return errors.New("ChmodSvcFailed")
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
		nameStr+".service",
	)
	if err != nil {
		return errors.New("SystemctlEnableFailed")
	}

	return nil
}

func (repo ServerCmdRepo) AddOneTimerSvc(
	svcName valueObject.SvcName,
	cmd valueObject.SvcCmd,
) error {
	name := svcName.String()
	svcFilePath := "/etc/systemd/system/" + name + ".service"
	svcContent := `[Unit]
Description=` + name + `
After=network.target

[Service]
Type=oneshot
User=root
WorkingDirectory=/var/speedia
Restart=no
StandardOutput=journal
StandardError=journal
SyslogIdentifier=` + name + `
ExecStart=` + cmd.String() + `
RemainAfterExit=yes
`

	err := infraHelper.UpdateFile(svcFilePath, svcContent, true)
	if err != nil {
		return errors.New("AddOneTimerSvcFailed")
	}
	err = os.Chmod(svcFilePath, 0644)
	if err != nil {
		return errors.New("ChmodSvcFailed")
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
	err = os.Chmod(svcTimerFilePath, 0644)
	if err != nil {
		return errors.New("ChmodSvcTimerFailed")
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
		return errors.New("SystemctlEnableSvcFailed")
	}

	_, err = infraHelper.RunCmd(
		"systemctl",
		"enable",
		name+".timer",
	)
	if err != nil {
		return errors.New("SystemctlEnableTimerFailed")
	}

	return nil
}

func (repo ServerCmdRepo) DeleteOneTimerSvc(svcName valueObject.SvcName) error {
	name := svcName.String()
	_, _ = infraHelper.RunCmd("systemctl", "stop", name+".timer")
	_, _ = infraHelper.RunCmd("systemctl", "disable", name+".timer")

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
	logFilePath := "/var/log/control.log"
	jsonLogEntry, err := json.Marshal(logEntity)
	if err != nil {
		return
	}
	logContent := string(jsonLogEntry) + "\n"

	logFile, err := os.OpenFile(
		logFilePath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0640,
	)
	if err != nil {
		return
	}
	defer logFile.Close()

	_, _ = logFile.WriteString(logContent)
	log.Print(logContent)
}

func (repo ServerCmdRepo) SendServerMessage(message string) {
	_, _ = infraHelper.RunCmd("wall", "-n", message)
}
