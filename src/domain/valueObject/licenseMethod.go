package valueObject

import (
	"errors"
	"strings"

	"golang.org/x/exp/slices"
)

type LicenseMethod string

var ValidLicenseMethods = []string{
	"ip",
	"key",
}

func NewLicenseMethod(value string) (LicenseMethod, error) {
	value = strings.TrimSpace(value)
	value = strings.ToLower(value)
	LicenseMethod := LicenseMethod(value)
	if !LicenseMethod.isValid() {
		return "", errors.New("InvalidLicenseMethod")
	}
	return LicenseMethod, nil
}

func NewLicenseMethodPanic(value string) LicenseMethod {
	LicenseMethod, err := NewLicenseMethod(value)
	if err != nil {
		panic(err)
	}
	return LicenseMethod
}

func (LicenseMethod LicenseMethod) isValid() bool {
	return slices.Contains(ValidLicenseMethods, LicenseMethod.String())
}

func (LicenseMethod LicenseMethod) String() string {
	return string(LicenseMethod)
}
