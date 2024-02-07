package valueObject

import (
	"errors"
	"regexp"
)

const licenseFingerprintRegex = `^\w{6,256}\-?\w{0,256}\-?\w{0,256}$`

type LicenseFingerprint string

func NewLicenseFingerprint(value string) (LicenseFingerprint, error) {
	fingerprint := LicenseFingerprint(value)
	if !fingerprint.isValid() {
		return "", errors.New("InvalidLicenseFingerprint")
	}

	return fingerprint, nil
}

func NewLicenseFingerprintPanic(value string) LicenseFingerprint {
	fingerprint, err := NewLicenseFingerprint(value)
	if err != nil {
		panic(err)
	}
	return fingerprint
}

func (fingerprint LicenseFingerprint) isValid() bool {
	fingerprintRegex := regexp.MustCompile(licenseFingerprintRegex)
	return fingerprintRegex.MatchString(string(fingerprint))
}

func (fingerprint LicenseFingerprint) String() string {
	return string(fingerprint)
}
