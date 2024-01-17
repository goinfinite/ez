package valueObject

import (
	"errors"
	"regexp"
	"strings"
)

const containerImgAddressRegex string = `^(?P<schema>https?://)?(?P<fqdn>[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?(?:\.[a-z0-9][a-z0-9-]{0,61}[a-z0-9])+)(?::(?P<port>\d{1,6}))?/(?:(?P<orgName>\w{1,128})/)?(?P<imageName>\w{1,128}):?(?P<imageTag>[\w\.\_\-]{1,128})?$`

type ContainerImgAddress string

func NewContainerImgAddress(value string) (ContainerImgAddress, error) {
	imageAddr := ContainerImgAddress(value)
	if !imageAddr.isValid() {
		return "", errors.New("InvalidContainerImageAddress")
	}

	imageAddrWithoutSchema := strings.TrimPrefix(string(imageAddr), "http://")
	imageAddrWithoutSchema = strings.TrimPrefix(imageAddrWithoutSchema, "https://")
	return ContainerImgAddress(imageAddrWithoutSchema), nil
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
