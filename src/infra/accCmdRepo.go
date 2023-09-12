package infra

import (
	"encoding/hex"
	"errors"
	"log"
	"os"
	"os/exec"
	"os/user"
	"time"

	"github.com/google/uuid"
	"github.com/speedianet/sfm/src/domain/dto"
	"github.com/speedianet/sfm/src/domain/entity"
	"github.com/speedianet/sfm/src/domain/valueObject"
	"github.com/speedianet/sfm/src/infra/db"
	dbModel "github.com/speedianet/sfm/src/infra/db/model"
	infraHelper "github.com/speedianet/sfm/src/infra/helper"
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

func getUsernameById(accId valueObject.AccountId) (valueObject.Username, error) {
	accQuery := AccQueryRepo{}
	accDetails, err := accQuery.GetById(accId)
	if err != nil {
		log.Printf("GetUserDetailsError: %s", err)
		return "", errors.New("GetUserDetailsError")
	}

	return accDetails.Username, nil
}

func (repo AccCmdRepo) Delete(accId valueObject.AccountId) error {
	username, err := getUsernameById(accId)
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

	dbSvc, err := db.DatabaseService()
	if err != nil {
		return err
	}

	err = dbModel.Account{ID: uint(accId.Get())}.Delete(dbSvc)
	if err != nil {
		log.Printf("DeleteAccountDbError: %s", err)
		return errors.New("DeleteAccountDbError")
	}

	return nil
}

func (repo AccCmdRepo) UpdatePassword(
	accId valueObject.AccountId,
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

	username, err := getUsernameById(accId)
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

	dbSvc, err := db.DatabaseService()
	if err != nil {
		return err
	}

	err = dbSvc.Model(&dbModel.Account{ID: uint(accId.Get())}).
		Update("updated_at", time.Now()).Error
	if err != nil {
		log.Printf("UpdateAccountDbError: %s", err)
		return errors.New("UpdateAccountDbError")
	}

	return nil
}

func (repo AccCmdRepo) UpdateApiKey(
	accId valueObject.AccountId,
) (valueObject.AccessTokenStr, error) {
	uuid := uuid.New()
	secretKey := os.Getenv("ACC_API_KEY_SECRET")
	apiKeyPlainText := accId.String() + ":" + uuid.String()

	encryptedApiKey, err := infraHelper.EncryptStr(secretKey, apiKeyPlainText)
	if err != nil {
		return "", errors.New("ApiKeyEncryptionError")
	}

	apiKey, err := valueObject.NewAccessTokenStr(encryptedApiKey)
	if err != nil {
		return "", errors.New("ApiKeyParseError")
	}

	hash := sha3.New256()
	hash.Write([]byte(uuid.String()))
	uuidHash := hex.EncodeToString(hash.Sum(nil))

	dbSvc, err := db.DatabaseService()
	if err != nil {
		return "", err
	}

	err = dbSvc.Model(&dbModel.Account{ID: uint(accId.Get())}).
		Update("key_hash", uuidHash).Error
	if err != nil {
		log.Printf("UpdateAccountDbError: %s", err)
		return "", errors.New("UpdateAccountDbError")
	}

	return apiKey, nil
}

func (repo AccCmdRepo) updateQuotaTable(
	tableName string,
	accId valueObject.AccountId,
	quota valueObject.AccountQuota,
) error {
	dbSvc, err := db.DatabaseService()
	if err != nil {
		return err
	}

	updateMap := map[string]interface{}{}

	if quota.CpuCores.Get() > 0 {
		updateMap["cpu_cores"] = quota.CpuCores.Get()
	}

	if quota.MemoryBytes.Get() > 0 {
		updateMap["memory_bytes"] = uint64(quota.MemoryBytes.Get())
	}

	if quota.DiskBytes.Get() > 0 {
		updateMap["disk_bytes"] = uint64(quota.DiskBytes.Get())
	}

	if quota.Inodes.Get() > 0 {
		updateMap["inodes"] = quota.Inodes.Get()
	}

	err = dbSvc.Table(tableName).
		Where("account_id = ?", uint(accId.Get())).
		Updates(updateMap).Error
	if err != nil {
		log.Printf("UpdateAccountQuotaDbError: %s", err)
		return errors.New("UpdateAccountQuotaDbError")
	}

	return nil
}

func (repo AccCmdRepo) UpdateQuota(
	accId valueObject.AccountId,
	quota valueObject.AccountQuota,
) error {
	return repo.updateQuotaTable(
		dbModel.AccountQuota{}.TableName(),
		accId,
		quota,
	)
}

func (repo AccCmdRepo) UpdateQuotaUsage(
	accId valueObject.AccountId,
	quota valueObject.AccountQuota,
) error {
	return repo.updateQuotaTable(
		dbModel.AccountQuotaUsage{}.TableName(),
		accId,
		quota,
	)
}
