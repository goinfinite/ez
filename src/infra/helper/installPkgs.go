package infraHelper

import (
	"errors"
	"log"
)

func InstallPkgs(packages []string) error {
	installPackages := append(
		[]string{"pkg", "install", "-y"},
		packages...,
	)

	var installErr error
	nAttempts := 3
	for i := 0; i < nAttempts; i++ {
		out, err := RunCmd("transactional-update", installPackages...)
		log.Printf("InstallPkgOut: %s", out)
		if err == nil {
			break
		}

		log.Printf("InstallPkgError: %s", err.Error())

		if i == nAttempts-1 {
			installErr = errors.New("InstallAttemptsFailed")
		}
	}

	return installErr
}
