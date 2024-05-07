package valueObject

import (
	"errors"
	"regexp"
	"strings"

	voHelper "github.com/speedianet/control/src/domain/valueObject/helper"
)

const containerImageAddressRegex string = `^(?P<schema>https?://)?(?P<fqdn>[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?(?:\.[a-z0-9][a-z0-9-]{0,61}[a-z0-9])+)?(?::(?P<port>\d{1,6}))?/?(?:(?P<orgName>[\w\_\-]{1,128})/)?(?P<imageName>[\w\_\-]{1,128}):?(?P<imageTag>[\w\.\_\-]{1,128})?$`

type ContainerImageAddress string

func NewContainerImageAddress(value interface{}) (ContainerImageAddress, error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return "", errors.New("ContainerImageAddressMustBeString")
	}

	stringValue = strings.TrimSpace(stringValue)
	stringValue = strings.ToLower(stringValue)

	valueParts := voHelper.FindNamedGroupsMatches(containerImageAddressRegex, stringValue)
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

	re := regexp.MustCompile(containerImageAddressRegex)
	isValid := re.MatchString(stringValue)
	if !isValid {
		return "", errors.New("InvalidContainerImageAddress")
	}

	return ContainerImageAddress(stringValue), nil
}

func NewContainerImageAddressPanic(value string) ContainerImageAddress {
	imageAddress, err := NewContainerImageAddress(value)
	if err != nil {
		panic(err)
	}
	return imageAddress
}

func (vo ContainerImageAddress) String() string {
	return string(vo)
}

func (vo ContainerImageAddress) getParts() map[string]string {
	return voHelper.FindNamedGroupsMatches(containerImageAddressRegex, string(vo))
}

func (vo ContainerImageAddress) GetFqdn() (Fqdn, error) {
	return NewFqdn(vo.getParts()["fqdn"])
}

func (vo ContainerImageAddress) GetOrgName() (RegistryPublisherName, error) {
	orgNameStr, exists := vo.getParts()["orgName"]
	if !exists || orgNameStr == "" || orgNameStr == "_" {
		orgNameStr = "library"
	}

	return NewRegistryPublisherName(orgNameStr)
}

func (vo ContainerImageAddress) GetImageName() (RegistryImageName, error) {
	return NewRegistryImageName(vo.getParts()["imageName"])
}

func (vo ContainerImageAddress) GetImageTag() (RegistryImageTag, error) {
	imageTagStr, exists := vo.getParts()["imageTag"]
	if !exists || imageTagStr == "" {
		imageTagStr = "latest"
	}

	return NewRegistryImageTag(imageTagStr)
}

func (vo ContainerImageAddress) IsSpeediaOs() bool {
	return strings.Contains(vo.String(), "speedia/os")
}
