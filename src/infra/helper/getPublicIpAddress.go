package infraHelper

import (
	"errors"
	"io"
	"net/http"

	"github.com/goinfinite/fleet/src/domain/valueObject"
)

func GetPublicIpAddress() (valueObject.IpAddress, error) {
	resp, err := http.Get("https://goinfinite.net/ip")
	if err != nil {
		return "", errors.New("GetPublicIpAddressFailed")
	}
	defer resp.Body.Close()

	ip, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("ReadPublicIpAddressFailed")
	}

	return valueObject.NewIpAddress(string(ip))
}
