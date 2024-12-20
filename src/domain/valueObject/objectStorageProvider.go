package valueObject

import (
	"errors"
	"strings"

	voHelper "github.com/goinfinite/ez/src/domain/valueObject/helper"
)

type ObjectStorageProvider string

const (
	ObjectStorageProviderAkamai       ObjectStorageProvider = "akamai"
	ObjectStorageProviderAws          ObjectStorageProvider = "aws"
	ObjectStorageProviderAzure        ObjectStorageProvider = "azure"
	ObjectStorageProviderBackblaze    ObjectStorageProvider = "backblaze"
	ObjectStorageProviderCloudFlare   ObjectStorageProvider = "cloudflare"
	ObjectStorageProviderDigitalOcean ObjectStorageProvider = "digitalocean"
	ObjectStorageProviderGoogleCloud  ObjectStorageProvider = "google-cloud"
	ObjectStorageProviderLinode       ObjectStorageProvider = "linode"
	ObjectStorageProviderMagalu       ObjectStorageProvider = "magalu"
	ObjectStorageProviderWasabi       ObjectStorageProvider = "wasabi"
	ObjectStorageProviderCustom       ObjectStorageProvider = "custom"
)

func NewObjectStorageProvider(value interface{}) (
	provider ObjectStorageProvider, err error,
) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return provider, errors.New("ObjectStorageProviderMustBeString")
	}
	stringValue = strings.ToLower(stringValue)

	stringValueVo := ObjectStorageProvider(stringValue)
	switch stringValueVo {
	case ObjectStorageProviderAkamai,
		ObjectStorageProviderAws,
		ObjectStorageProviderAzure,
		ObjectStorageProviderBackblaze,
		ObjectStorageProviderCloudFlare,
		ObjectStorageProviderDigitalOcean,
		ObjectStorageProviderGoogleCloud,
		ObjectStorageProviderMagalu,
		ObjectStorageProviderWasabi,
		ObjectStorageProviderCustom:
		return stringValueVo, nil
	case ObjectStorageProviderLinode:
		return ObjectStorageProviderAkamai, nil
	default:
		return provider, errors.New("InvalidObjectStorageProvider")
	}
}

func (vo ObjectStorageProvider) String() string {
	return string(vo)
}
