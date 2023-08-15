package valueObject

import (
	"errors"
	"regexp"
)

const containerImgAddressRegex string = `^(?P<registryScheme>https?://)?(?P<registryHostname>\w{1,10}\.\w{1,10}\.?\w{0,10})?/?(?P<orgName>\w{1,100})/(?P<repoName>\w{1,100}):?(?P<repoTag>\w{1,100})$`

type ContainerImgAddress string

func NewContainerImgAddress(value string) (ContainerImgAddress, error) {
	cntrImgAddr := ContainerImgAddress(value)
	if !cntrImgAddr.isValid() {
		return "", errors.New("InvalidContainerImgAddress")
	}
	return cntrImgAddr, nil
}

func NewContainerImgAddressPanic(value string) ContainerImgAddress {
	cntrImgAddr := ContainerImgAddress(value)
	if !cntrImgAddr.isValid() {
		panic("InvalidContainerImgAddress")
	}
	return cntrImgAddr
}

func (cntrImgAddr ContainerImgAddress) isValid() bool {
	re := regexp.MustCompile(containerImgAddressRegex)
	return re.MatchString(string(cntrImgAddr))
}

func (cntrImgAddr ContainerImgAddress) getRegexParts() map[string]string {
	re := regexp.MustCompile(containerImgAddressRegex)
	match := re.FindStringSubmatch(string(cntrImgAddr))

	groupNames := re.SubexpNames()
	groupMap := make(map[string]string)

	for i, name := range groupNames {
		if i != 0 && name != "" {
			groupMap[name] = match[i]
		}
	}

	return groupMap
}

func (cntrImgAddr ContainerImgAddress) GetRegistryScheme() string {
	parts := cntrImgAddr.getRegexParts()
	return parts["registryScheme"]
}

func (cntrImgAddr ContainerImgAddress) GetRegistryHostname() string {
	parts := cntrImgAddr.getRegexParts()
	return parts["registryHostname"]
}

func (cntrImgAddr ContainerImgAddress) GetOrgName() string {
	parts := cntrImgAddr.getRegexParts()
	return parts["orgName"]
}

func (cntrImgAddr ContainerImgAddress) GetRepoName() string {
	parts := cntrImgAddr.getRegexParts()
	return parts["repoName"]
}

func (cntrImgAddr ContainerImgAddress) GetRepoTag() string {
	parts := cntrImgAddr.getRegexParts()
	return parts["repoTag"]
}

func (cntrImgAddr ContainerImgAddress) String() string {
	return string(cntrImgAddr)
}
