package infraHelper

import (
	"os"
	"testing"

	testHelpers "github.com/speedianet/sfm/src/devUtils"
	"github.com/speedianet/sfm/src/domain/valueObject"
)

func TestRunCmd(t *testing.T) {
	testHelpers.LoadEnvVars()

	t.Run("RunCmd", func(t *testing.T) {
		command := "echo"
		args := []string{"hello", "world"}

		stdOut, err := RunCmd(command, args...)
		if err != nil {
			t.Errorf("Expected no error, got %s", err)
		}

		expectedStdOut := "hello world"
		if stdOut != expectedStdOut {
			t.Errorf("Expected %s, got %s", expectedStdOut, stdOut)
		}
	})

	t.Run("RunCmdWithError", func(t *testing.T) {
		command := "nonexistentcommand"
		args := []string{"hello", "world"}

		_, err := RunCmd(command, args...)
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})

	t.Run("RunCmdAsUser", func(t *testing.T) {
		accId := valueObject.NewAccountIdPanic(os.Getenv("DUMMY_USER_ID"))
		command := "whoami"
		args := []string{}

		stdOut, err := RunCmdAsUser(accId, command, args...)
		if err != nil {
			t.Errorf("Expected no error, got %s", err)
		}

		expectedStdOut := os.Getenv("DUMMY_USER_NAME")
		if stdOut != expectedStdOut {
			t.Errorf("Expected %s, got %s", expectedStdOut, stdOut)
		}
	})
}