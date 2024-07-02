package infra

import (
	"encoding/hex"
	"errors"
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/valueObject"
	"github.com/speedianet/control/src/infra/db"
	dbModel "github.com/speedianet/control/src/infra/db/model"
	infraHelper "github.com/speedianet/control/src/infra/helper"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/sha3"
)

type AuthQueryRepo struct {
	persistentDbSvc *db.PersistentDatabaseService
}

func NewAuthQueryRepo(persistentDbSvc *db.PersistentDatabaseService) *AuthQueryRepo {
	return &AuthQueryRepo{persistentDbSvc: persistentDbSvc}
}

func (repo *AuthQueryRepo) IsLoginValid(login dto.Login) bool {
	storedPassHash, err := infraHelper.RunCmd(
		"bash",
		"-c",
		"getent shadow "+login.Username.String()+" | awk -F: '{print $2}'",
	)
	if err != nil {
		return false
	}

	if len(storedPassHash) == 0 {
		return false
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(storedPassHash),
		[]byte(login.Password.String()),
	)
	return err == nil
}

func (repo *AuthQueryRepo) getSessionTokenClaims(
	sessionToken valueObject.AccessTokenValue,
) (jwt.MapClaims, error) {
	var claims jwt.MapClaims

	parsedToken, err := jwt.Parse(
		sessionToken.String(),
		func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
	if err != nil {
		if err.(*jwt.ValidationError).Errors == jwt.ValidationErrorExpired {
			return claims, errors.New("SessionTokenExpired")
		}

		return claims, errors.New("SessionTokenParseError: " + err.Error())
	}

	claims, areClaimsReadable := parsedToken.Claims.(jwt.MapClaims)
	if !areClaimsReadable {
		return claims, errors.New("SessionTokenClaimsUnreadable")
	}

	return claims, nil
}

func (repo *AuthQueryRepo) getTokenDetailsFromSession(
	sessionTokenClaims jwt.MapClaims,
) (dto.AccessTokenDetails, error) {
	issuedIp, err := valueObject.NewIpAddress(
		sessionTokenClaims["originalIp"].(string),
	)
	if err != nil {
		return dto.AccessTokenDetails{}, errors.New("OriginalIpUnreadable")
	}

	accountId, err := valueObject.NewAccountId(sessionTokenClaims["accountId"])
	if err != nil {
		return dto.AccessTokenDetails{}, errors.New("AccountIdUnreadable")
	}

	return dto.NewAccessTokenDetails(
		valueObject.NewAccessTokenTypePanic("sessionToken"),
		accountId,
		&issuedIp,
	), nil
}

func (repo *AuthQueryRepo) getKeyHash(
	accountId valueObject.AccountId,
) (string, error) {
	accModel := dbModel.Account{ID: uint(accountId.Read())}
	err := repo.persistentDbSvc.Handler.Model(&accModel).First(&accModel).Error
	if err != nil {
		return "", errors.New("AccountNotFound")
	}

	if accModel.KeyHash == nil {
		return "", errors.New("UserKeyHashNotFound")
	}

	return *accModel.KeyHash, nil
}

func (repo *AuthQueryRepo) getTokenDetailsFromApiKey(
	token valueObject.AccessTokenValue,
) (dto.AccessTokenDetails, error) {
	secretKey := os.Getenv("ACC_API_KEY_SECRET")
	decryptedApiKey, err := infraHelper.DecryptStr(secretKey, token.String())
	if err != nil {
		return dto.AccessTokenDetails{}, errors.New("ApiKeyDecryptionError")
	}

	// keyFormat: accountId:UUIDv4
	keyParts := strings.Split(decryptedApiKey, ":")
	if len(keyParts) != 2 {
		return dto.AccessTokenDetails{}, errors.New("ApiKeyFormatError")
	}

	accountId, err := valueObject.NewAccountId(keyParts[0])
	if err != nil {
		return dto.AccessTokenDetails{}, errors.New("AccountIdUnreadable")
	}
	uuid := keyParts[1]

	uuidHash := sha3.New256()
	uuidHash.Write([]byte(uuid))
	uuidHashStr := hex.EncodeToString(uuidHash.Sum(nil))

	storedUuidHash, err := repo.getKeyHash(accountId)
	if err != nil {
		return dto.AccessTokenDetails{}, errors.New("UserKeyHashUnreadable")
	}

	if uuidHashStr != storedUuidHash {
		return dto.AccessTokenDetails{}, errors.New("UserKeyHashMismatch")
	}

	return dto.NewAccessTokenDetails(
		valueObject.NewAccessTokenTypePanic("accountApiKey"),
		accountId,
		nil,
	), nil
}

func (repo *AuthQueryRepo) ReadAccessTokenDetails(
	token valueObject.AccessTokenValue,
) (dto.AccessTokenDetails, error) {
	var tokenDetails dto.AccessTokenDetails

	sessionTokenClaims, err := repo.getSessionTokenClaims(token)
	if err != nil {
		if err.Error() == "SessionTokenExpired" {
			return tokenDetails, err
		}

		return repo.getTokenDetailsFromApiKey(token)
	}

	return repo.getTokenDetailsFromSession(sessionTokenClaims)
}
