package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/control/src/domain/repository"
	"github.com/speedianet/control/src/domain/valueObject"
)

const (
	DaysUntilSuspension      int = 7
	DaysUntilRevocation      int = 14
	LicenseValidationsPerDay int = 4
)

func LicenseValidation(
	licenseQueryRepo repository.LicenseQueryRepo,
	licenseCmdRepo repository.LicenseCmdRepo,
) error {
	log.Print("LicenseValidationStarted")

	err := licenseCmdRepo.Refresh()
	if err != nil {
		log.Printf("RefreshLicenseInfoError: %s", err)

		err := licenseCmdRepo.IncrementErrorCount()
		if err != nil {
			return errors.New("IncrementLicenseErrorCountError: " + err.Error())
		}
	}

	licenseInfo, err := licenseQueryRepo.Get()
	if err != nil {
		return errors.New("GetLicenseInfoError: " + err.Error())
	}

	licenseStatusStr := licenseInfo.Status.String()
	if licenseStatusStr == "ACTIVE" {
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
		status, _ := valueObject.NewLicenseStatus("SUSPENDED")
		newLicenseStatus = &status
	}

	if newLicenseStatus != nil {
		err = licenseCmdRepo.UpdateStatus(*newLicenseStatus)
		if err != nil {
			return errors.New("UpdateLicenseStatusError: " + err.Error())
		}
	}

	log.Print("LicenseValidationFinished")
	return nil
}
