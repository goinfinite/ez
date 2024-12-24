package valueObject

import (
	"errors"
	"net"
	"regexp"
	"strings"

	voHelper "github.com/goinfinite/ez/src/domain/valueObject/helper"
)

type ObjectStorageBucketName string

const objectStorageBucketNameRegexExpression = `^[a-z0-9][a-z0-9\.\-]{0,256}[a-z0-9]$`

func NewObjectStorageBucketName(value interface{}) (
	bucketName ObjectStorageBucketName, err error,
) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return bucketName, errors.New("ObjectStorageBucketNameMustBeString")
	}
	stringValue = strings.ToLower(stringValue)

	if strings.Contains(stringValue, "..") {
		return bucketName, errors.New("ObjectStorageBucketNameCannotContainDoublePeriods")
	}

	if net.ParseIP(stringValue) != nil {
		return bucketName, errors.New("ObjectStorageBucketNameCannotBeAnIpAddress")
	}

	re := regexp.MustCompile(objectStorageBucketNameRegexExpression)
	if !re.MatchString(stringValue) {
		return bucketName, errors.New("InvalidObjectStorageBucketName")
	}

	return ObjectStorageBucketName(stringValue), nil
}

func (vo ObjectStorageBucketName) String() string {
	return string(vo)
}
