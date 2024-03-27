package infraHelper

import (
	"errors"
	"strings"

	"github.com/speedianet/control/src/domain/valueObject"
	"github.com/speedianet/control/src/infra/db"
)

const PublicIpTransientKey string = "PublicIp"

func GetPublicIpAddress(
	transientDbSvc *db.TransientDatabaseService,
) (valueObject.IpAddress, error) {
	var ipAddress valueObject.IpAddress

	cachedIpAddressStr, err := transientDbSvc.Get(PublicIpTransientKey)
	if err == nil {
		return valueObject.NewIpAddress(cachedIpAddressStr)
	}

	rawIpEntry, err := RunCmd(
		"dig", "+short", "TXT", "o-o.myaddr.l.google.com", "@ns1.google.com",
	)
	if err != nil {
		rawIpEntry, err = RunCmd(
			"dig", "+short", "TXT", "CH", "whoami.cloudflare", "@1.1.1.1",
		)
		if err != nil {
			return ipAddress, errors.New("GetPublicIpFailed: " + err.Error())
		}
	}

	rawIpEntry = strings.Trim(rawIpEntry, `"`)
	if rawIpEntry == "" {
		return ipAddress, errors.New("GetPublicIpFailed: NoIpEntry")
	}

	ipAddress, err = valueObject.NewIpAddress(rawIpEntry)
	if err != nil {
		return ipAddress, err
	}

	err = transientDbSvc.Set(PublicIpTransientKey, ipAddress.String())
	if err != nil {
		return ipAddress, errors.New("PersistPublicIpFailed: " + err.Error())
	}

	return ipAddress, nil
}
