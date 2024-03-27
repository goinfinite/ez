package infraHelper

import (
	"errors"
	"io"
	"net/http"

	"github.com/speedianet/control/src/domain/valueObject"
	"github.com/speedianet/control/src/infra/db"
)

const PublicIpTransientKey string = "PublicIp"

func GetPublicIpAddress(
	transientDbSvc *db.TransientDatabaseService,
) (valueObject.IpAddress, error) {
	cachedIpAddressStr, err := transientDbSvc.Get(PublicIpTransientKey)
	if err == nil {
		return valueObject.NewIpAddress(cachedIpAddressStr)
	}

	resp, err := http.Get("https://speedia.net/ip")
	if err != nil {
		return "", errors.New("GetPublicIpAddressFailed")
	}
	defer resp.Body.Close()

	ipAddressBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("ReadPublicIpAddressFailed")
	}

	ipAddressStr := string(ipAddressBytes)
	ipAddress, err := valueObject.NewIpAddress(ipAddressStr)
	if err != nil {
		return "", err
	}

	err = transientDbSvc.Set(PublicIpTransientKey, ipAddress.String())
	if err != nil {
		return ipAddress, errors.New("FailedToPersistPublicIp: " + err.Error())
	}

	return ipAddress, nil
}
