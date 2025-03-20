package infra

import (
	"errors"
	"log/slog"
	"os"
	"strings"

	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/valueObject"
	"github.com/goinfinite/ez/src/infra/db"
	dbModel "github.com/goinfinite/ez/src/infra/db/model"
	infraHelper "github.com/goinfinite/ez/src/infra/helper"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

var (
	errSessionTokenExpired          = errors.New("SessionTokenExpired")
	errSessionTokenSignatureInvalid = errors.New("SessionTokenSignatureInvalid")
	errSessionTokenParseError       = errors.New("SessionTokenParseError")
	errSessionTokenClaimsUnreadable = errors.New("SessionTokenClaimsUnreadable")
)

type AuthQueryRepo struct {
	persistentDbSvc *db.PersistentDatabaseService
}

func NewAuthQueryRepo(persistentDbSvc *db.PersistentDatabaseService) *AuthQueryRepo {
	return &AuthQueryRepo{persistentDbSvc: persistentDbSvc}
}

func (repo *AuthQueryRepo) IsLoginValid(createDto dto.CreateSessionToken) bool {
	storedPassHash, err := infraHelper.RunCmdWithSubShell(
		"getent shadow " + createDto.Username.String() + " | awk -F: '{print $2}'",
	)
	if err != nil {
		slog.Debug(
			"GetentShadowError",
			slog.String("username", createDto.Username.String()),
			slog.Any("error", err),
		)
		return false
	}

	if len(storedPassHash) == 0 {
		return false
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(storedPassHash),
		[]byte(createDto.Password.String()),
	)
	return err == nil
}

func (repo *AuthQueryRepo) getSessionTokenClaims(
	sessionToken valueObject.AccessTokenValue,
) (claims jwt.MapClaims, err error) {
	parsedToken, err := jwt.Parse(
		sessionToken.String(),
		func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
	if err != nil {
		switch errorEnum := err.(*jwt.ValidationError).Errors; errorEnum {
		case jwt.ValidationErrorExpired:
			return claims, errSessionTokenExpired
		case jwt.ValidationErrorSignatureInvalid:
			return claims, errSessionTokenSignatureInvalid
		default:
			return claims, errSessionTokenParseError
		}
	}

	claims, areClaimsReadable := parsedToken.Claims.(jwt.MapClaims)
	if !areClaimsReadable {
		return claims, errSessionTokenClaimsUnreadable
	}

	return claims, nil
}

func (repo *AuthQueryRepo) getTokenDetailsFromSession(
	sessionTokenClaims jwt.MapClaims,
) (tokenDetails dto.AccessTokenDetails, err error) {
	tokenType, _ := valueObject.NewAccessTokenType("sessionToken")

	accountId, err := valueObject.NewAccountId(sessionTokenClaims["accountId"])
	if err != nil {
		return tokenDetails, errors.New("AccountIdUnreadable")
	}

	issuedIp, err := valueObject.NewIpAddress(sessionTokenClaims["originalIp"])
	if err != nil {
		return tokenDetails, errors.New("OriginalIpUnreadable")
	}

	return dto.NewAccessTokenDetails(tokenType, accountId, &issuedIp), nil
}

func (repo *AuthQueryRepo) getKeyHash(
	accountId valueObject.AccountId,
) (string, error) {
	accountModel := dbModel.Account{ID: accountId.Uint64()}
	err := repo.persistentDbSvc.Handler.Model(&accountModel).First(&accountModel).Error
	if err != nil {
		return "", errors.New("AccountNotFound")
	}

	if accountModel.KeyHash == nil {
		return "", errors.New("UserKeyHashNotFound")
	}

	return *accountModel.KeyHash, nil
}

func (repo *AuthQueryRepo) getTokenDetailsFromApiKey(
	token valueObject.AccessTokenValue,
) (tokenDetails dto.AccessTokenDetails, err error) {
	secretKey := os.Getenv("ACC_API_KEY_SECRET")
	decryptedApiKey, err := infraHelper.DecryptStr(secretKey, token.String())
	if err != nil {
		return tokenDetails, errors.New("ApiKeyDecryptionError")
	}

	// keyFormat: accountId:UUIDv4
	keyParts := strings.Split(decryptedApiKey, ":")
	if len(keyParts) != 2 {
		return tokenDetails, errors.New("ApiKeyFormatError")
	}

	accountId, err := valueObject.NewAccountId(keyParts[0])
	if err != nil {
		return tokenDetails, errors.New("AccountIdUnreadable")
	}

	uuidHash := infraHelper.GenStrongHash(keyParts[1])

	storedUuidHash, err := repo.getKeyHash(accountId)
	if err != nil {
		return tokenDetails, errors.New("UserKeyHashUnreadable")
	}

	if uuidHash != storedUuidHash {
		return tokenDetails, errors.New("UserKeyHashMismatch")
	}

	tokenType, _ := valueObject.NewAccessTokenType("accountApiKey")

	return dto.NewAccessTokenDetails(tokenType, accountId, nil), nil
}

func (repo *AuthQueryRepo) ReadAccessTokenDetails(
	token valueObject.AccessTokenValue,
) (tokenDetails dto.AccessTokenDetails, err error) {
	sessionTokenClaims, err := repo.getSessionTokenClaims(token)
	if err != nil {
		isLikelyApiKey := err == errSessionTokenParseError
		if !isLikelyApiKey {
			return tokenDetails, err
		}

		return repo.getTokenDetailsFromApiKey(token)
	}

	return repo.getTokenDetailsFromSession(sessionTokenClaims)
}
