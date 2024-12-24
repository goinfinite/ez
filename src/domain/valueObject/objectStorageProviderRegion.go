package valueObject

import (
	"errors"
	"net"
	"regexp"
	"strings"

	voHelper "github.com/goinfinite/ez/src/domain/valueObject/helper"
)

type ObjectStorageProviderRegion string

const objectStorageProviderRegionRegexExpression = `^[a-z0-9][a-z0-9\.\-]{0,256}[a-z0-9]$`

func NewObjectStorageProviderRegion(value interface{}) (
	providerRegion ObjectStorageProviderRegion, err error,
) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return providerRegion, errors.New("ObjectStorageProviderRegionMustBeString")
	}
	stringValue = strings.ToLower(stringValue)

	if strings.Contains(stringValue, "..") {
		return providerRegion, errors.New("ObjectStorageProviderRegionCannotContainDoublePeriods")
	}

	if net.ParseIP(stringValue) != nil {
		return providerRegion, errors.New("ObjectStorageProviderRegionCannotBeAnIpAddress")
	}

	re := regexp.MustCompile(objectStorageProviderRegionRegexExpression)
	if !re.MatchString(stringValue) {
		return providerRegion, errors.New("InvalidObjectStorageProviderRegion")
	}

	return ObjectStorageProviderRegion(stringValue), nil
}

func (vo ObjectStorageProviderRegion) String() string {
	return string(vo)
}
