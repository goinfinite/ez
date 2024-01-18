package valueObject

import (
	"errors"
	"regexp"
	"strings"

	voHelper "github.com/speedianet/control/src/domain/valueObject/helper"
)

const containerImgAddressRegex string = `^(?P<schema>https?://)?(?P<fqdn>[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?(?:\.[a-z0-9][a-z0-9-]{0,61}[a-z0-9])+)?(?::(?P<port>\d{1,6}))?/?(?:(?P<orgName>\w{1,128})/)?(?P<imageName>\w{1,128}):?(?P<imageTag>[\w\.\_\-]{1,128})?$`

type ContainerImageAddress string

func NewContainerImageAddress(value string) (ContainerImageAddress, error) {
	valueParts := voHelper.FindNamedGroupsMatches(containerImgAddressRegex, value)
	if len(valueParts) == 0 {
		return "", errors.New("UnknownImageAddressFormat")
	}

	if valueParts["schema"] != "" {
		value = strings.TrimPrefix(value, valueParts["schema"])
	}

	if valueParts["fqdn"] == "" {
		value = "docker.io/" + value
	}

	if !strings.Contains(value, "/") {
		return "", errors.New("ImageAddressMustContainOrgAndImageName")
	}

	imageAddress := ContainerImageAddress(value)
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
