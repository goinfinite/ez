package dbModel

import (
	"time"

	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/valueObject"
)

type LicenseInfo struct {
	ID          uint `gorm:"primarykey"`
	Method      string
	Status      string
	Fingerprint string
	ErrorCount  uint `gorm:"default:0"`
	ExpiresAt   time.Time
	LastCheckAt time.Time
	UpdatedAt   time.Time
}

func (LicenseInfo) TableName() string {
	return "license_info"
}

func (LicenseInfo) InitialEntries() []interface{} {
	initialEntry := LicenseInfo{
		ID:          1,
		Method:      "ip",
		Status:      "ACTIVE",
		Fingerprint: "fingerprint",
		ErrorCount:  0,
		ExpiresAt:   time.Now(),
		LastCheckAt: time.Now(),
		UpdatedAt:   time.Now(),
	}

	return []interface{}{
		initialEntry,
	}
}

func NewLicenseInfo(
	method string,
	status string,
	fingerprint string,
	errorCount uint,
	expiresAt time.Time,
	lastCheckAt time.Time,
) LicenseInfo {
	licenseInfoModel := LicenseInfo{
		ID:          1,
		Method:      method,
		Status:      status,
		Fingerprint: fingerprint,
		ErrorCount:  errorCount,
		ExpiresAt:   expiresAt,
		LastCheckAt: lastCheckAt,
	}

	return licenseInfoModel
}

func (LicenseInfo) ToModel(entity entity.LicenseInfo) LicenseInfo {
	expiresAt := time.Unix(entity.ExpiresAt.Read(), 0)
	lastCheckAt := time.Unix(entity.LastCheckAt.Read(), 0)

	return NewLicenseInfo(
		entity.Method.String(),
		entity.Status.String(),
		entity.Fingerprint.String(),
		entity.ErrorCount,
		expiresAt,
		lastCheckAt,
	)
}

func (model LicenseInfo) ToEntity() (entity.LicenseInfo, error) {
	var licenseInfo entity.LicenseInfo

	licenseMethod, err := valueObject.NewLicenseMethod(model.Method)
	if err != nil {
		return licenseInfo, err
	}

	licenseStatus, err := valueObject.NewLicenseStatus(model.Status)
	if err != nil {
		return licenseInfo, err
	}

	licenseFingerprint, err := valueObject.NewLicenseFingerprint(model.Fingerprint)
	if err != nil {
		return licenseInfo, err
	}

	expiresAt := valueObject.NewUnixTimeWithGoTime(model.ExpiresAt)
	lastCheckAt := valueObject.NewUnixTimeWithGoTime(model.LastCheckAt)
	updatedAt := valueObject.NewUnixTimeWithGoTime(model.UpdatedAt)

	licenseInfo = entity.NewLicenseInfo(
		licenseMethod,
		licenseStatus,
		licenseFingerprint,
		model.ErrorCount,
		expiresAt,
		lastCheckAt,
		updatedAt,
	)

	return licenseInfo, nil
}
