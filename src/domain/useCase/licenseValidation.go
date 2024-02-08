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

	err := licenseCmdRepo.RefreshStatus()
	if err != nil {
		log.Printf("RefreshLicenseStatusError: %s", err)

		err := licenseCmdRepo.IncrementErrorCount()
		if err != nil {
			return errors.New("IncrementLicenseErrorCountError: " + err.Error())
		}
	}

	licenseStatus, err := licenseQueryRepo.GetStatus()
	if err != nil {
		return errors.New("GetLicenseStatusError: " + err.Error())
	}

	licenseStatusStr := licenseStatus.Status.String()
	if licenseStatusStr == "ACTIVE" {
		err = licenseCmdRepo.ResetErrorCount()
		if err != nil {
			return errors.New("ResetLicenseErrorCountError: " + err.Error())
		}
		log.Print("LicenseValidatedSuccessfully")
		return nil
	}

	maxErrorCountUntilSuspension := DaysUntilSuspension * LicenseValidationsPerDay
	maxErrorCountUntilRevocation := DaysUntilRevocation * LicenseValidationsPerDay

	errorCount, err := licenseQueryRepo.GetErrorCount()
	if err != nil {
		return errors.New("GetLicenseErrorCountError: " + err.Error())
	}

	var newLicenseStatus *valueObject.LicenseStatus
	switch {
	case errorCount > maxErrorCountUntilRevocation:
		log.Print("LicenseErrorCountExceedsRevocationTolerance")
		status, _ := valueObject.NewLicenseStatus("REVOKED")
		newLicenseStatus = &status
	case errorCount > maxErrorCountUntilSuspension:
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
