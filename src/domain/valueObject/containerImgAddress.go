package valueObject

import (
	"errors"
	"regexp"
	"strings"

	voHelper "github.com/speedianet/control/src/domain/valueObject/helper"
)

const containerImgAddressRegex string = `^(?P<schema>https?://)?(?P<fqdn>[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?(?:\.[a-z0-9][a-z0-9-]{0,61}[a-z0-9])+)?(?::(?P<port>\d{1,6}))?/?(?:(?P<orgName>\w{1,128})/)?(?P<imageName>\w{1,128}):?(?P<imageTag>[\w\.\_\-]{1,128})?$`

type ContainerImgAddress string

func NewContainerImgAddress(value string) (ContainerImgAddress, error) {
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
	value = "https://" + value

	slashesCount := strings.Count(value, "/")
	minSlashesAmount := 3
	if slashesCount < minSlashesAmount {
		return "", errors.New("ImageAddressDoesNotContainEnoughSlashes")
	}

	imageAddr := ContainerImgAddress(value)
	if !imageAddr.isValid() {
		return "", errors.New("InvalidContainerImageAddress")
	}

	return imageAddr, nil
}

func NewContainerImgAddressPanic(value string) ContainerImgAddress {
	imageAddr, err := NewContainerImgAddress(value)
	if err != nil {
		panic(err)
	}
	return imageAddr
}

func (imageAddr ContainerImgAddress) isValid() bool {
	re := regexp.MustCompile(containerImgAddressRegex)
	return re.MatchString(string(imageAddr))
}

func (imageAddr ContainerImgAddress) String() string {
	return string(imageAddr)
}
