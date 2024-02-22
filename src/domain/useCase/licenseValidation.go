package useCase

import (
	"errors"
	"log"
	"strings"

	"github.com/speedianet/control/src/domain/repository"
	"github.com/speedianet/control/src/domain/valueObject"
)

const (
	DaysUntilSuspension      int = 7
	DaysUntilRevocation      int = 14
	LicenseValidationsPerDay int = 4
)

func isLicenseNonceValid(
	nonceHash valueObject.Hash,
	fingerprint valueObject.LicenseFingerprint,
) bool {
	fingerprintParts := strings.Split(fingerprint.String(), "-")
	if len(fingerprintParts) != 3 {
		return false
	}

	fingerprintNonceHash, err := valueObject.NewHash(fingerprintParts[2])
	if err != nil {
		return false
	}

	return nonceHash == fingerprintNonceHash
}

func LicenseValidation(
	licenseQueryRepo repository.LicenseQueryRepo,
	licenseCmdRepo repository.LicenseCmdRepo,
) error {
	log.Print("LicenseValidationStarted")

	refreshOk := true
	err := licenseCmdRepo.Refresh()
	if err != nil {
		log.Printf("RefreshLicenseInfoError: %s", err)
		refreshOk = false

		err := licenseCmdRepo.IncrementErrorCount()
		if err != nil {
			return errors.New("IncrementLicenseErrorCountError: " + err.Error())
		}
	}

	licenseInfo, err := licenseQueryRepo.Get()
	if err != nil {
		return errors.New("GetLicenseInfoError: " + err.Error())
	}

	freshNonceHash, err := licenseCmdRepo.GenerateNonceHash()
	if err != nil {
		return errors.New("GenerateLicenseNonceHashError: " + err.Error())
	}

	suspendedStatus, _ := valueObject.NewLicenseStatus("SUSPENDED")

	if refreshOk && !isLicenseNonceValid(freshNonceHash, licenseInfo.Fingerprint) {
		err = licenseCmdRepo.UpdateStatus(suspendedStatus)
		if err != nil {
			return errors.New("UpdateLicenseStatusError: " + err.Error())
		}

		return errors.New("LicenseIntegrityCheckFailed")
	}

	integrityHash, err := licenseCmdRepo.GenerateIntegrityHash(licenseInfo)
	if err != nil {
		return errors.New("GenerateLicenseIntegrityHashError: " + err.Error())
	}

	persistedIntegrityHash, err := licenseQueryRepo.GetIntegrityHash()
	if err != nil {
		return errors.New("GetPersistedLicenseIntegrityHashError: " + err.Error())
	}

	if integrityHash != persistedIntegrityHash {
		err = licenseCmdRepo.UpdateStatus(suspendedStatus)
		if err != nil {
			return errors.New("UpdateLicenseStatusError: " + err.Error())
		}

		return errors.New("LicenseIntegrityCheckFailed")
	}

	if refreshOk && licenseInfo.Status.String() == "ACTIVE" {
		err = licenseCmdRepo.ResetErrorCount()
		if err != nil {
			return errors.New("ResetLicenseErrorCountError: " + err.Error())
		}
		log.Print("LicenseValidatedSuccessfully")
		return nil
	}

	maxErrorCountUntilSuspension := uint(DaysUntilSuspension * LicenseValidationsPerDay)
	maxErrorCountUntilRevocation := uint(DaysUntilRevocation * LicenseValidationsPerDay)

	var newLicenseStatus *valueObject.LicenseStatus
	switch {
	case licenseInfo.ErrorCount > maxErrorCountUntilRevocation:
		log.Print("LicenseErrorCountExceedsRevocationTolerance")
		status, _ := valueObject.NewLicenseStatus("REVOKED")
		newLicenseStatus = &status
	case licenseInfo.ErrorCount > maxErrorCountUntilSuspension:
		log.Print("LicenseErrorCountExceedsSuspensionTolerance")
		newLicenseStatus = &suspendedStatus
	}

	if newLicenseStatus != nil {
		err = licenseCmdRepo.UpdateStatus(*newLicenseStatus)
		if err != nil {
			return errors.New("UpdateLicenseStatusError: " + err.Error())
		}
		return errors.New("InvalidLicenseStatusSystemShutdown")
	}

	log.Print("LicenseValidationFinished")
	return nil
}
