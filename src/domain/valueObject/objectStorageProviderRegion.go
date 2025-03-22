package valueObject

import (
	"errors"
	"net"
	"regexp"
	"strings"

	voHelper "github.com/goinfinite/ez/src/domain/valueObject/helper"
)

type ObjectStorageProviderRegion string

var KnownObjectStorageProviderRegions = map[ObjectStorageProvider][]ObjectStorageProviderRegion{
	ObjectStorageProviderAlibaba: {
		"oss-cn-hangzhou.aliyuncs.com", "oss-cn-shanghai.aliyuncs.com",
		"oss-cn-qingdao.aliyuncs.com",
	},
	ObjectStorageProviderAws: {
		"us-east-1", "us-east-2", "us-west-1", "us-west-2",
		"ap-south-1", "ap-northeast-1", "ap-northeast-2", "ap-southeast-1",
		"ap-southeast-2", "ca-central-1", "eu-central-1", "eu-west-1",
		"eu-west-2", "eu-west-3", "sa-east-1",
	},
	ObjectStorageProviderCloudFlare: {"auto"},
	ObjectStorageProviderDigitalOcean: {
		"nyc3.digitaloceanspaces.com", "ams3.digitaloceanspaces.com",
		"sfo2.digitaloceanspaces.com", "sfo3.digitaloceanspaces.com",
		"sgp1.digitaloceanspaces.com", "lon1.digitaloceanspaces.com",
		"fra1.digitaloceanspaces.com", "tor1.digitaloceanspaces.com",
		"blr1.digitaloceanspaces.com", "syd1.digitaloceanspaces.com",
	},
	ObjectStorageProviderGoogleCloud: {"storage.googleapis.com"},
	ObjectStorageProviderLinode: {
		"nl-ams-1.linodeobjects.com", "us-southeast-1.linodeobjects.com",
		"in-maa-1.linodeobjects.com", "us-ord-1.linodeobjects.com",
		"eu-central-1.linodeobjects.com", "id-cgk-1.linodeobjects.com",
		"gb-lon-1.linodeobjects.com", "us-lax-1.linodeobjects.com",
		"es-mad-1.linodeobjects.com", "au-mel-1.linodeobjects.com",
		"us-mia-1.linodeobjects.com", "it-mil-1.linodeobjects.com",
		"us-east-1.linodeobjects.com", "jp-osa-1.linodeobjects.com",
		"fr-par-1.linodeobjects.com", "br-gru-1.linodeobjects.com",
		"us-sea-1.linodeobjects.com", "ap-south-1.linodeobjects.com",
		"sg-sin-1.linodeobjects.com", "se-sto-1.linodeobjects.com",
		"us-iad-1.linodeobjects.com",
	},
	ObjectStorageProviderMagalu: {"br-se1.magaluobjects.com", "br-ne1.magaluobjects.com"},
	ObjectStorageProviderWasabi: {
		"us-west-1", "us-east-1", "us-east-2", "us-central-1", "ca-central-1",
		"eu-west-1", "eu-west-2", "eu-west-3", "eu-central-1", "eu-central-2",
		"eu-south-1", "ap-northeast-1", "ap-northeast-2", "ap-southeast-1",
		"ap-southeast-2",
	},
}

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
