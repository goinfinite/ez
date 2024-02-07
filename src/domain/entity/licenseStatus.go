package entity

import "github.com/speedianet/control/src/domain/valueObject"

type LicenseStatus struct {
	Method      valueObject.LicenseMethod      `json:"method"`
	Status      valueObject.LicenseStatus      `json:"status"`
	Fingerprint valueObject.LicenseFingerprint `json:"fingerprint"`
	ExpiresAt   valueObject.UnixTime           `json:"expiresAt"`
	LastCheckAt valueObject.UnixTime           `json:"lastCheckAt"`
}

func NewLicenseStatus(
	method valueObject.LicenseMethod,
	status valueObject.LicenseStatus,
	fingerprint valueObject.LicenseFingerprint,
	expiresAt valueObject.UnixTime,
	lastCheckAt valueObject.UnixTime,
) LicenseStatus {
	return LicenseStatus{
		Method:      method,
		Status:      status,
		Fingerprint: fingerprint,
		ExpiresAt:   expiresAt,
		LastCheckAt: lastCheckAt,
	}
}
