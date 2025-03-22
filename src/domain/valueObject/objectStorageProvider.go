package valueObject

import (
	"errors"
	"strings"

	voHelper "github.com/goinfinite/ez/src/domain/valueObject/helper"
)

type ObjectStorageProvider string

const (
	ObjectStorageProviderCustom       ObjectStorageProvider = "custom"
	ObjectStorageProviderAkamai       ObjectStorageProvider = "akamai"
	ObjectStorageProviderAlibaba      ObjectStorageProvider = "alibaba"
	ObjectStorageProviderAws          ObjectStorageProvider = "aws"
	ObjectStorageProviderCloudFlare   ObjectStorageProvider = "cloudflare"
	ObjectStorageProviderDigitalOcean ObjectStorageProvider = "digitalocean"
	ObjectStorageProviderGoogleCloud  ObjectStorageProvider = "google-cloud"
	ObjectStorageProviderLinode       ObjectStorageProvider = "linode"
	ObjectStorageProviderMagalu       ObjectStorageProvider = "magalu"
	ObjectStorageProviderR2           ObjectStorageProvider = "r2"
	ObjectStorageProviderWasabi       ObjectStorageProvider = "wasabi"
)

var ObjectStorageProviderStrList = []string{
	ObjectStorageProviderCustom.String(),
	ObjectStorageProviderAkamai.String(),
	ObjectStorageProviderAlibaba.String(),
	ObjectStorageProviderAws.String(),
	ObjectStorageProviderCloudFlare.String(),
	ObjectStorageProviderDigitalOcean.String(),
	ObjectStorageProviderGoogleCloud.String(),
	ObjectStorageProviderLinode.String(),
	ObjectStorageProviderMagalu.String(),
	ObjectStorageProviderR2.String(),
	ObjectStorageProviderWasabi.String(),
}

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
	case ObjectStorageProviderAlibaba,
		ObjectStorageProviderAws,
		ObjectStorageProviderCloudFlare,
		ObjectStorageProviderDigitalOcean,
		ObjectStorageProviderGoogleCloud,
		ObjectStorageProviderLinode,
		ObjectStorageProviderMagalu,
		ObjectStorageProviderWasabi,
		ObjectStorageProviderCustom:
		return stringValueVo, nil
	case ObjectStorageProviderAkamai:
		return ObjectStorageProviderLinode, nil
	case ObjectStorageProviderR2:
		return ObjectStorageProviderCloudFlare, nil
	default:
		return provider, errors.New("InvalidObjectStorageProvider")
	}
}

func (vo ObjectStorageProvider) String() string {
	return string(vo)
}
