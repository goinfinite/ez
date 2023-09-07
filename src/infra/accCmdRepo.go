package infra

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"
	"log"
	"os"
	"os/exec"
	"os/user"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/speedianet/sfm/src/domain/dto"
	"github.com/speedianet/sfm/src/domain/entity"
	"github.com/speedianet/sfm/src/domain/valueObject"
	"github.com/speedianet/sfm/src/infra/db"
	dbModel "github.com/speedianet/sfm/src/infra/db/model"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/sha3"
)

type AccCmdRepo struct {
}

func (repo AccCmdRepo) Add(addAccount dto.AddAccount) error {
	passHash, err := bcrypt.GenerateFromPassword(
		[]byte(addAccount.Password.String()),
		bcrypt.DefaultCost,
	)
	if err != nil {
		log.Printf("PasswordHashError: %s", err)
		return errors.New("PasswordHashError")
	}

	addAccountCmd := exec.Command(
		"useradd",
		"-m",
		"-s", "/bin/bash",
		"-p", string(passHash),
		addAccount.Username.String(),
	)

	err = addAccountCmd.Run()
	if err != nil {
		log.Printf("AccountAddError: %s", err)
		return errors.New("AccountAddError")
	}

	userInfo, err := user.Lookup(addAccount.Username.String())
	if err != nil {
		return errors.New("AccountLookupError")
	}
	accId, err := valueObject.NewAccountId(userInfo.Uid)
	if err != nil {
		return errors.New("AccountIdParseError")
	}
	gid, err := valueObject.NewGroupId(userInfo.Gid)
	if err != nil {
		return errors.New("GroupIdParseError")
	}

	dbSvc, err := db.DatabaseService()
	if err != nil {
		return err
	}

	nowUnixTime := valueObject.UnixTime(time.Now().Unix())
	accEntity := entity.NewAccount(
		accId,
		gid,
		addAccount.Username,
		addAccount.Quota,
		valueObject.NewAccountQuotaWithBlankValues(),
		nowUnixTime,
		nowUnixTime,
	)

	accModel, err := dbModel.Account{}.ToModel(accEntity)
	if err != nil {
		log.Printf("AccountModelParseError: %s", err)
		return errors.New("AccountModelParseError")
	}

	dbResult := dbSvc.Create(&accModel)
	if dbResult.Error != nil {
		log.Printf("AddAccountDbError: %s", dbResult.Error)
		return errors.New("AddAccountDbError")
	}

	return nil
}

func getUsernameById(accountId valueObject.AccountId) (valueObject.Username, error) {
	accQuery := AccQueryRepo{}
	accDetails, err := accQuery.GetById(accountId)
	if err != nil {
		log.Printf("GetUserDetailsError: %s", err)
		return "", errors.New("GetUserDetailsError")
	}

	return accDetails.Username, nil
}

func (repo AccCmdRepo) Delete(accountId valueObject.AccountId) error {
	username, err := getUsernameById(accountId)
	if err != nil {
		return err
	}

	delUserCmd := exec.Command(
		"userdel",
		"-r",
		username.String(),
	)

	err = delUserCmd.Run()
	if err != nil {
		log.Printf("UserDeleteError: %s", err)
		return errors.New("UserDeleteError")
	}

	return nil
}

func (repo AccCmdRepo) UpdatePassword(
	accountId valueObject.AccountId,
	password valueObject.Password,
) error {
	passHash, err := bcrypt.GenerateFromPassword(
		[]byte(password.String()),
		bcrypt.DefaultCost,
	)
	if err != nil {
		log.Printf("PasswordHashError: %s", err)
		return errors.New("PasswordHashError")
	}

	username, err := getUsernameById(accountId)
	if err != nil {
		return err
	}

	updateAccountCmd := exec.Command(
		"usermod",
		"-p", string(passHash),
		username.String(),
	)

	err = updateAccountCmd.Run()
	if err != nil {
		log.Printf("PasswordUpdateError: %s", err)
		return errors.New("PasswordUpdateError")
	}

	return nil
}

func encryptApiKey(plainTextApiKey string) (valueObject.AccessTokenStr, error) {
	secretKey := os.Getenv("UAK_SECRET")
	secretKeyBytes, err := base64.RawURLEncoding.DecodeString(secretKey)
	if err != nil {
		log.Printf("ApiKeySecretKeyError: %s", err)
		return "", errors.New("ApiKeySecretKeyError")
	}

	block, err := aes.NewCipher(secretKeyBytes)
	if err != nil {
		log.Printf("ApiKeyCipherError: %s", err)
		return "", errors.New("ApiKeyCipherError")
	}

	plainTextApiKeyBytes := []byte(plainTextApiKey)
	cipherText := make([]byte, aes.BlockSize+len(plainTextApiKeyBytes))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		log.Printf("ApiKeyIvGenerationError: %s", err)
		return "", errors.New("ApiKeyIvGenerationError")
	}

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainTextApiKeyBytes)

	newApiKey, err := valueObject.NewAccessTokenStr(
		base64.StdEncoding.EncodeToString(cipherText),
	)
	if err != nil {
		log.Printf("ApiKeyEncodingError: %s", err)
		return "", errors.New("ApiKeyEncodingError")
	}

	return newApiKey, nil
}

func storeNewKeyHash(
	accountId valueObject.AccountId,
	uuid uuid.UUID,
) error {
	var lock sync.Mutex
	lock.Lock()
	defer lock.Unlock()

	keysHashFile := ".accountApiKeys"

	if _, err := os.Stat(keysHashFile); err == nil {
		purgeOldKeyCmd := exec.Command(
			"sed",
			"-i",
			"/"+accountId.String()+":/d",
			keysHashFile,
		)
		err := purgeOldKeyCmd.Run()
		if err != nil {
			log.Printf("PurgeOldKeyError: %s", err)
			return errors.New("PurgeOldKeyError")
		}
	}

	hash := sha3.New256()
	hash.Write([]byte(uuid.String()))
	hashString := hex.EncodeToString(hash.Sum(nil))

	file, err := os.OpenFile(keysHashFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0400)
	if err != nil {
		log.Printf("KeysFileUnreadable: %v", err)
		return errors.New("KeysFileUnreadable")
	}
	defer file.Close()

	_, err = file.WriteString(accountId.String() + ":" + hashString + "\n")
	if err != nil {
		log.Printf("UserKeysWriteError: %v", err)
		return errors.New("UserKeysWriteError")
	}

	return nil
}

func (repo AccCmdRepo) UpdateApiKey(
	accountId valueObject.AccountId,
) (valueObject.AccessTokenStr, error) {
	uuid := uuid.New()
	plainTextApiKey := accountId.String() + ":" + uuid.String()
	newApiKey, err := encryptApiKey(plainTextApiKey)
	if err != nil {
		return "", err
	}

	err = storeNewKeyHash(accountId, uuid)
	if err != nil {
		return "", err
	}

	return newApiKey, nil
}
