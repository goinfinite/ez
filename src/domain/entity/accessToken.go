package entity

import "github.com/goinfinite/ez/src/domain/valueObject"

type AccessToken struct {
	Type      valueObject.AccessTokenType  `json:"type"`
	ExpiresIn valueObject.UnixTime         `json:"expiresIn"`
	TokenStr  valueObject.AccessTokenValue `json:"tokenStr"`
}

func NewAccessToken(
	tokenType valueObject.AccessTokenType,
	expiresIn valueObject.UnixTime,
	tokenStr valueObject.AccessTokenValue,
) AccessToken {
	return AccessToken{
		Type:      tokenType,
		ExpiresIn: expiresIn,
		TokenStr:  tokenStr,
	}
}
