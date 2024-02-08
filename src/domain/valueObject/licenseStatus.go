package valueObject

import (
	"errors"
	"strings"

	"golang.org/x/exp/slices"
)

type LicenseStatus string

var ValidLicenseStatuses = []string{
	"ACTIVE",
	"EXPIRED",
	"SUSPENDED",
	"REVOKED",
	"TERMINATED",
}

func NewLicenseStatus(value string) (LicenseStatus, error) {
	value = strings.TrimSpace(value)
	value = strings.ToUpper(value)
	status := LicenseStatus(value)
	if !status.isValid() {
		return "", errors.New("InvalidLicenseStatus")
	}
	return status, nil
}

func NewLicenseStatusPanic(value string) LicenseStatus {
	status, err := NewLicenseStatus(value)
	if err != nil {
		panic(err)
	}
	return status
}

func (status LicenseStatus) isValid() bool {
	return slices.Contains(ValidLicenseStatuses, status.String())
}

func (status LicenseStatus) String() string {
	return string(status)
}
