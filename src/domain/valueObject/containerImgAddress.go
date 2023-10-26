package valueObject

import (
	"errors"
	"regexp"
	"strings"
)

const containerImgAddressRegex string = `^(?P<schema>https?://)?(?P<fqdn>[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?(?:\.[a-z0-9][a-z0-9-]{0,61}[a-z0-9])+)(?::(?P<port>\d{1,6}))?/(?:(?P<orgName>\w{1,128})/)?(?P<imageName>\w{1,128}):?(?P<imageTag>\w{1,128})?$`

type ContainerImgAddress string

func NewContainerImgAddress(value string) (ContainerImgAddress, error) {
	imgAddr := ContainerImgAddress(value)
	if !imgAddr.isValid() {
		return "", errors.New("InvalidContainerImgAddress")
	}

	imgAddrWithoutSchema := strings.TrimPrefix(string(imgAddr), "http://")
	imgAddrWithoutSchema = strings.TrimPrefix(imgAddrWithoutSchema, "https://")
	return ContainerImgAddress(imgAddrWithoutSchema), nil
}

func NewContainerImgAddressPanic(value string) ContainerImgAddress {
	imgAddr, err := NewContainerImgAddress(value)
	if err != nil {
		panic(err)
	}
	return imgAddr
}

func (imgAddr ContainerImgAddress) isValid() bool {
	re := regexp.MustCompile(containerImgAddressRegex)
	return re.MatchString(string(imgAddr))
}

func (imgAddr ContainerImgAddress) String() string {
	return string(imgAddr)
}
