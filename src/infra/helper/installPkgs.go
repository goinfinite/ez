package infraHelper

import (
	"errors"
	"log/slog"
)

func InstallPkgs(packages []string) error {
	installPackages := append(
		[]string{"pkg", "install", "-y"},
		packages...,
	)

	var installErr error
	nAttempts := 3
	for i := 0; i < nAttempts; i++ {
		_, err := RunCmd("transactional-update", installPackages...)
		if err == nil {
			break
		}

		slog.Error("InstallPkgError", slog.Any("error", err))

		if i == nAttempts-1 {
			installErr = errors.New("InstallAttemptsFailed")
		}
	}

	return installErr
}
