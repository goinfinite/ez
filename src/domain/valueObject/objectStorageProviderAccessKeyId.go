package valueObject

import (
	"errors"
	"regexp"

	voHelper "github.com/goinfinite/ez/src/domain/valueObject/helper"
)

type ObjectStorageProviderAccessKeyId string

const objectStorageProviderAccessKeyIdRegexExpression = `^[A-Za-z0-9]{10,256}$`

func NewObjectStorageProviderAccessKeyId(value interface{}) (
	keyId ObjectStorageProviderAccessKeyId, err error,
) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return keyId, errors.New("ObjectStorageProviderAccessKeyIdMustBeString")
	}

	re := regexp.MustCompile(objectStorageProviderAccessKeyIdRegexExpression)
	if !re.MatchString(stringValue) {
		return keyId, errors.New("InvalidObjectStorageProviderAccessKeyId")
	}

	return ObjectStorageProviderAccessKeyId(stringValue), nil
}

func (vo ObjectStorageProviderAccessKeyId) String() string {
	return string(vo)
}
