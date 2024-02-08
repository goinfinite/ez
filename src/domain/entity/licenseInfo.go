package entity

import "github.com/speedianet/control/src/domain/valueObject"

type LicenseInfo struct {
	Method      valueObject.LicenseMethod      `json:"method"`
	Status      valueObject.LicenseStatus      `json:"status"`
	Fingerprint valueObject.LicenseFingerprint `json:"fingerprint"`
	ExpiresAt   valueObject.UnixTime           `json:"expiresAt"`
	LastCheckAt valueObject.UnixTime           `json:"lastCheckAt"`
}

func NewLicenseInfo(
	method valueObject.LicenseMethod,
	status valueObject.LicenseStatus,
	fingerprint valueObject.LicenseFingerprint,
	expiresAt valueObject.UnixTime,
	lastCheckAt valueObject.UnixTime,
) LicenseInfo {
	return LicenseInfo{
		Method:      method,
		Status:      status,
		Fingerprint: fingerprint,
		ExpiresAt:   expiresAt,
		LastCheckAt: lastCheckAt,
	}
}
