package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/repository"
	"github.com/speedianet/control/src/domain/valueObject"
)

func GetAccessTokenDetails(
	authQueryRepo repository.AuthQueryRepo,
	accessToken valueObject.AccessTokenStr,
	trustedIpAddress []valueObject.IpAddress,
	ipAddress valueObject.IpAddress,
) (dto.AccessTokenDetails, error) {
	var tokenDetails dto.AccessTokenDetails

	accessTokenDetails, err := authQueryRepo.GetAccessTokenDetails(accessToken)
	if err != nil {
		log.Printf("GetAccessTokenDetailsError: %s", err)
		return tokenDetails, errors.New("GetAccessTokenDetailsInfraError")
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
		return tokenDetails, errors.New("AccessTokenIpMismatch")
	}

	return accessTokenDetails, nil
}
