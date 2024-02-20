package entity

import "github.com/speedianet/control/src/domain/valueObject"

type LicenseInfo struct {
	Method      valueObject.LicenseMethod      `json:"method"`
	Status      valueObject.LicenseStatus      `json:"status"`
	Fingerprint valueObject.LicenseFingerprint `json:"fingerprint"`
	ErrorCount  uint                           `json:"errorCount"`
	ExpiresAt   valueObject.UnixTime           `json:"expiresAt"`
	LastCheckAt valueObject.UnixTime           `json:"lastCheckAt"`
	UpdatedAt   valueObject.UnixTime           `json:"updatedAt"`
}

func NewLicenseInfo(
	method valueObject.LicenseMethod,
	status valueObject.LicenseStatus,
	fingerprint valueObject.LicenseFingerprint,
	errorCount uint,
	expiresAt valueObject.UnixTime,
	lastCheckAt valueObject.UnixTime,
	updatedAt valueObject.UnixTime,
) LicenseInfo {
	return LicenseInfo{
		Method:      method,
		Status:      status,
		Fingerprint: fingerprint,
		ErrorCount:  errorCount,
		ExpiresAt:   expiresAt,
		LastCheckAt: lastCheckAt,
		UpdatedAt:   updatedAt,
	}
}
