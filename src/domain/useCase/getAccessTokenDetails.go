package useCase

import (
	"errors"

	"github.com/goinfinite/fleet/src/domain/dto"
	"github.com/goinfinite/fleet/src/domain/repository"
	"github.com/goinfinite/fleet/src/domain/valueObject"
)

func GetAccessTokenDetails(
	authQueryRepo repository.AuthQueryRepo,
	accessToken valueObject.AccessTokenStr,
	trustedIpAddress []valueObject.IpAddress,
	ipAddress valueObject.IpAddress,
) (dto.AccessTokenDetails, error) {
	accessTokenDetails, err := authQueryRepo.GetAccessTokenDetails(accessToken)
	if err != nil {
		return dto.AccessTokenDetails{}, err
	}

	if accessTokenDetails.IpAddress == nil {
		return accessTokenDetails, nil
	}

	for _, trustedIp := range trustedIpAddress {
		if trustedIp.String() == ipAddress.String() {
			return accessTokenDetails, nil
		}
	}

	if accessTokenDetails.IpAddress.String() != ipAddress.String() {
		return dto.AccessTokenDetails{}, errors.New("IpAddressChanged")
	}

	return accessTokenDetails, nil
}
