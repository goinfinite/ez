package valueObject

import (
	"errors"
	"regexp"

	voHelper "github.com/goinfinite/ez/src/domain/valueObject/helper"
)

type ObjectStorageProviderSecretAccessKey string

const objectStorageProviderSecretAccessKeyRegexExpression = `^[A-Za-z0-9/+=]{10,256}$`

func NewObjectStorageProviderSecretAccessKey(value interface{}) (
	secretKey ObjectStorageProviderSecretAccessKey, err error,
) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return secretKey, errors.New("ObjectStorageProviderSecretAccessKeyMustBeString")
	}

	re := regexp.MustCompile(objectStorageProviderSecretAccessKeyRegexExpression)
	if !re.MatchString(stringValue) {
		return secretKey, errors.New("InvalidObjectStorageProviderSecretAccessKey")
	}

	return ObjectStorageProviderSecretAccessKey(stringValue), nil
}

func (vo ObjectStorageProviderSecretAccessKey) String() string {
	return string(vo)
}
