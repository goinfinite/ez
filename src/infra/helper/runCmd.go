package infraHelper

import (
	"bytes"
	"encoding/json"
	"os/exec"
	"strings"
	"syscall"

	"github.com/speedianet/sfm/src/domain/valueObject"
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
	execCmd := exec.Command(command, args...)
	execCmd.SysProcAttr = &syscall.SysProcAttr{}
	execCmd.SysProcAttr.Credential = &syscall.Credential{
		Uid: uint32(accId),
	}

	return runExecCmd(execCmd)
}
