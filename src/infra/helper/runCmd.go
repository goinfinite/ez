package infraHelper

import (
	"bytes"
	"encoding/json"
	"errors"
	"os/exec"
	"os/user"
	"strings"
	"syscall"

	"github.com/goinfinite/fleet/src/domain/valueObject"
)

type CommandError struct {
	StdErr   string `json:"stdErr"`
	ExitCode int    `json:"exitCode"`
}

func (e *CommandError) Error() string {
	errJSON, _ := json.Marshal(e)
	return string(errJSON)
}

func runExecCmd(execCmd *exec.Cmd) (string, error) {
	var stdout, stderr bytes.Buffer
	execCmd.Stdout = &stdout
	execCmd.Stderr = &stderr

	err := execCmd.Run()
	stdOut := strings.TrimSpace(stdout.String())
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return stdOut, &CommandError{
				StdErr:   stderr.String(),
				ExitCode: exitErr.ExitCode(),
			}
		}
		return stdOut, err
	}

	return stdOut, nil
}

func RunCmd(command string, args ...string) (string, error) {
	execCmd := exec.Command(command, args...)
	return runExecCmd(execCmd)
}

func RunCmdAsUser(
	accId valueObject.AccountId,
	command string,
	args ...string,
) (string, error) {
	userInfo, err := user.LookupId(accId.String())
	if err != nil {
		return "", errors.New("AccountIdNotFound")
	}

	gId, err := valueObject.NewGroupId(userInfo.Gid)
	if err != nil {
		return "", errors.New("GroupIdNotFound")
	}

	execCmd := exec.Command(command, args...)
	execCmd.SysProcAttr = &syscall.SysProcAttr{}
	execCmd.Dir = userInfo.HomeDir
	execCmd.SysProcAttr.Credential = &syscall.Credential{
		Uid: uint32(accId),
		Gid: uint32(gId),
	}
	execCmd.Env = []string{
		"HOME=" + userInfo.HomeDir,
		"LOGNAME=" + userInfo.Username,
		"USER=" + userInfo.Username,
		"SHELL=/usr/sbin/nologin",
		"PATH=/usr/sbin:/usr/bin:/sbin:/bin",
		"MAIL=/var/mail/" + userInfo.Username,
		"LANG=en_US.UTF-8",
		"PWD=" + userInfo.HomeDir,
		"TERM=xterm-256color",
		"XDG_RUNTIME_DIR=/run/user/" + accId.String(),
		"DBUS_SESSION_BUS_ADDRESS=unix:path=/run/user/" + accId.String() + "/bus",
	}
	return runExecCmd(execCmd)
}
