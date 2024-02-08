package useCase

import (
	"log"

	"github.com/speedianet/control/src/domain/repository"
	"github.com/speedianet/control/src/domain/valueObject"
)

const (
	DaysUntilSuspension int = 7
	DaysUntilRevocation int = 14
	LicenseChecksPerDay int = 4
)

func AutoCheckLicense(
	licenseQueryRepo repository.LicenseQueryRepo,
	licenseCmdRepo repository.LicenseCmdRepo,
) {
	err := licenseCmdRepo.RefreshStatus()
	if err != nil {
		log.Printf("RefreshLicenseStatusError: %s", err)

		err := licenseCmdRepo.IncrementErrorCount()
		if err != nil {
			panic("IncrementLicenseErrorCountError")
		}
	}

	licenseStatus, err := licenseQueryRepo.GetStatus()
	if err != nil {
		panic("GetLicenseStatusError")
	}

	licenseStatusStr := licenseStatus.Status.String()
	if licenseStatusStr == "ACTIVE" {
		err = licenseCmdRepo.ResetErrorCount()
		if err != nil {
			panic("ResetLicenseErrorCountError")
		}
		return
	}

	maxErrorCountUntilSuspension := DaysUntilSuspension * LicenseChecksPerDay
	maxErrorCountUntilRevocation := DaysUntilRevocation * LicenseChecksPerDay

	errorCount, err := licenseQueryRepo.GetErrorCount()
	if err != nil {
		panic("GetLicenseErrorCountError")
	}

	var newLicenseStatus *valueObject.LicenseStatus
	switch {
	case errorCount > maxErrorCountUntilRevocation:
		log.Print("LicenseErrorCountExceedsRevokeLimit")
		status, _ := valueObject.NewLicenseStatus("REVOKED")
		newLicenseStatus = &status
	case errorCount > maxErrorCountUntilSuspension:
		log.Print("LicenseErrorCountExceedsSuspendLimit")
		status, _ := valueObject.NewLicenseStatus("SUSPENDED")
		newLicenseStatus = &status
	}

	if newLicenseStatus != nil {
		err = licenseCmdRepo.UpdateStatus(*newLicenseStatus)
		if err != nil {
			panic("UpdateLicenseStatusError")
		}
	}
}
