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

	if !slices.Contains(ValidLicenseMethods, value) {
		return "", errors.New("InvalidLicenseMethod")
	}
	return LicenseMethod(value), nil
}

func NewLicenseMethodPanic(value string) LicenseMethod {
	licenseMethod, err := NewLicenseMethod(value)
	if err != nil {
		panic(err)
	}
	return licenseMethod
}

func (lm LicenseMethod) String() string {
	return string(lm)
}
