package valueObject

import (
	"errors"
	"regexp"
	"strings"

	voHelper "github.com/speedianet/control/src/domain/valueObject/helper"
)

const containerImgAddressRegex string = `^(?P<schema>https?://)?(?P<fqdn>[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?(?:\.[a-z0-9][a-z0-9-]{0,61}[a-z0-9])+)?(?::(?P<port>\d{1,6}))?/?(?:(?P<orgName>[\w\_\-]{1,128})/)?(?P<imageName>[\w\_\-]{1,128}):?(?P<imageTag>[\w\.\_\-]{1,128})?$`

type ContainerImageAddress string

func NewContainerImageAddress(value interface{}) (ContainerImageAddress, error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return "", errors.New("ContainerImageAddressMustBeString")
	}

	stringValue = strings.TrimSpace(stringValue)
	stringValue = strings.ToLower(stringValue)

	valueParts := voHelper.FindNamedGroupsMatches(containerImgAddressRegex, stringValue)
	if len(valueParts) == 0 {
		return "", errors.New("UnknownImageAddressFormat")
	}

	if valueParts["schema"] != "" {
		stringValue = strings.TrimPrefix(stringValue, valueParts["schema"])
	}

	if valueParts["fqdn"] == "" {
		stringValue = "docker.io/" + stringValue
	}

	if !strings.Contains(stringValue, "/") {
		return "", errors.New("ImageAddressMustContainOrgAndImageName")
	}

	imageAddress := ContainerImageAddress(stringValue)
	if !imageAddress.isValid() {
		return "", errors.New("InvalidContainerImageAddress")
	}

	return imageAddress, nil
}

func NewContainerImageAddressPanic(value string) ContainerImageAddress {
	imageAddress, err := NewContainerImageAddress(value)
	if err != nil {
		panic(err)
	}
	return imageAddress
}

func (imageAddress ContainerImageAddress) isValid() bool {
	re := regexp.MustCompile(containerImgAddressRegex)
	return re.MatchString(string(imageAddress))
}

func (imageAddress ContainerImageAddress) String() string {
	return string(imageAddress)
}

func (imageAddress ContainerImageAddress) getParts() map[string]string {
	return voHelper.FindNamedGroupsMatches(containerImgAddressRegex, string(imageAddress))
}

func (imageAddress ContainerImageAddress) GetFqdn() (Fqdn, error) {
	return NewFqdn(imageAddress.getParts()["fqdn"])
}

func (imageAddress ContainerImageAddress) GetOrgName() (RegistryPublisherName, error) {
	orgNameStr, exists := imageAddress.getParts()["orgName"]
	if !exists || orgNameStr == "" || orgNameStr == "_" {
		orgNameStr = "library"
	}

	return NewRegistryPublisherName(orgNameStr)
}

func (imageAddress ContainerImageAddress) GetImageName() (RegistryImageName, error) {
	return NewRegistryImageName(imageAddress.getParts()["imageName"])
}

func (imageAddress ContainerImageAddress) GetImageTag() (RegistryImageTag, error) {
	imageTagStr, exists := imageAddress.getParts()["imageTag"]
	if !exists || imageTagStr == "" {
		imageTagStr = "latest"
	}

	return NewRegistryImageTag(imageTagStr)
}
