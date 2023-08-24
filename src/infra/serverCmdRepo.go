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
WorkingDirectory=/speedia
ExecStart=` + cmdStr + `
Restart=always
RestartSec=15

[Install]
WantedBy=multi-user.target
`

	err := infraHelper.UpdateFile(svcFilePath, svcContent, true)
	if err != nil {
		return errors.New("AddSvcFailed")
	}
	os.Chmod(svcFilePath, 0755)

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
		"--now",
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
ExecStart=` + cmd.String() + `
RemainAfterExit=yes
`

	err := infraHelper.UpdateFile(svcFilePath, svcContent, true)
	if err != nil {
		return errors.New("AddOneTimerSvcFailed")
	}
	os.Chmod(svcFilePath, 0755)

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
	os.Chmod(svcTimerFilePath, 0755)

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

func (repo ServerCmdRepo) DeleteOneTimerSvc(svcName valueObject.SvcName) error {
	name := svcName.String()
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
		0640,
	)
	if err != nil {
		return
	}
	defer logFile.Close()

	logFile.WriteString(logContent)
	log.Print(logContent)
}

func (repo ServerCmdRepo) SendServerMessage(message string) {
	infraHelper.RunCmd("wall", "-n", message)
}
