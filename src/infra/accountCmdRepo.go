package infra

import (
	"errors"
	"log/slog"
	"os"
	"os/user"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/valueObject"
	"github.com/speedianet/control/src/infra/db"
	dbModel "github.com/speedianet/control/src/infra/db/model"
	infraHelper "github.com/speedianet/control/src/infra/helper"
	"golang.org/x/crypto/bcrypt"
)

type AccountCmdRepo struct {
	persistentDbSvc *db.PersistentDatabaseService
}

func NewAccountCmdRepo(persistentDbSvc *db.PersistentDatabaseService) *AccountCmdRepo {
	return &AccountCmdRepo{persistentDbSvc: persistentDbSvc}
}

func (repo *AccountCmdRepo) updateFilesystemQuota(
	accountId valueObject.AccountId, quota valueObject.AccountQuota,
) error {
	diskBytesStr := quota.StorageBytes.String()
	inodesStr := strconv.FormatUint(quota.StorageInodes, 10)
	accountIdStr := accountId.String()

	shouldUpdateDiskQuota := quota.StorageBytes.Read() > 0
	shouldUpdateInodeQuota := quota.StorageInodes > 0
	shouldRemoveQuota := !shouldUpdateDiskQuota && !shouldUpdateInodeQuota

	xfsFlags := "-x -c 'limit -u"
	if shouldUpdateDiskQuota {
		xfsFlags += " bhard=" + diskBytesStr
	}
	if shouldUpdateInodeQuota {
		xfsFlags += " ihard=" + inodesStr
	}
	if shouldRemoveQuota {
		xfsFlags = "-x -c 'limit -u bhard=0 ihard=0"
	}

	xfsFlags += " " + accountIdStr + "' /var/data"

	_, err := infraHelper.RunCmdWithSubShell("xfs_quota " + xfsFlags)
	return err
}

func (repo *AccountCmdRepo) Create(
	createDto dto.CreateAccount,
) (accountId valueObject.AccountId, err error) {
	passHash, err := bcrypt.GenerateFromPassword(
		[]byte(createDto.Password.String()),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return accountId, err
	}

	usernameStr := createDto.Username.String()
	_, err = infraHelper.RunCmd(
		"useradd", "-m",
		"-d", "/var/data/"+usernameStr,
		"-s", "/usr/sbin/nologin",
		"-p", string(passHash),
		usernameStr,
	)
	if err != nil {
		return accountId, errors.New("UserAddFailed: " + err.Error())
	}

	userInfo, err := user.Lookup(usernameStr)
	if err != nil {
		return accountId, err
	}
	accountId, err = valueObject.NewAccountId(userInfo.Uid)
	if err != nil {
		return accountId, err
	}
	gid, err := valueObject.NewUnixGroupId(userInfo.Gid)
	if err != nil {
		return accountId, err
	}

	nowUnixTime := valueObject.NewUnixTimeNow()
	accountEntity := entity.NewAccount(
		accountId, gid, createDto.Username, *createDto.Quota,
		valueObject.NewAccountQuotaWithBlankValues(),
		nowUnixTime, nowUnixTime,
	)

	accModel, err := dbModel.Account{}.ToModel(accountEntity)
	if err != nil {
		return accountId, err
	}

	err = repo.persistentDbSvc.Handler.Create(&accModel).Error
	if err != nil {
		return accountId, err
	}

	err = repo.updateFilesystemQuota(accountId, *createDto.Quota)
	if err != nil {
		return accountId, err
	}

	err = infraHelper.EnableLingering(accountId)
	if err != nil {
		return accountId, errors.New("EnableLingeringFailed: " + err.Error())
	}

	return accountId, nil
}

func (repo *AccountCmdRepo) getUsernameById(
	accountId valueObject.AccountId,
) (valueObject.Username, error) {
	accQuery := NewAccountQueryRepo(repo.persistentDbSvc)
	accDetails, err := accQuery.ReadById(accountId)
	if err != nil {
		return "", err
	}

	return accDetails.Username, nil
}

func (repo *AccountCmdRepo) Delete(accountId valueObject.AccountId) error {
	quota := valueObject.NewAccountQuotaWithBlankValues()
	err := repo.updateFilesystemQuota(accountId, quota)
	if err != nil {
		return err
	}

	username, err := repo.getUsernameById(accountId)
	if err != nil {
		return err
	}

	err = infraHelper.DisableLingering(accountId)
	if err != nil {
		return errors.New("DisableLingeringFailed: " + err.Error())
	}

	_, err = infraHelper.RunCmd("pgrep", "-u", accountId.String())
	if err == nil {
		_, _ = infraHelper.RunCmd("pkill", "-9", "-U", accountId.String())
	}

	_, err = infraHelper.RunCmd("userdel", "-r", username.String())
	if err != nil {
		return err
	}

	model := dbModel.Account{}
	modelId := accountId.Uint64()

	relatedTables := []string{
		dbModel.AccountQuota{}.TableName(),
		dbModel.AccountQuotaUsage{}.TableName(),
	}

	for _, tableName := range relatedTables {
		err := repo.persistentDbSvc.Handler.Exec(
			"DELETE FROM "+tableName+" WHERE account_id = ?", modelId,
		).Error
		if err != nil {
			return errors.New("DeleteAccRelatedTablesDbError: " + err.Error())
		}
	}

	err = repo.persistentDbSvc.Handler.Delete(model, modelId).Error
	if err != nil {
		return errors.New("DeleteAccDbError: " + err.Error())
	}

	return nil
}

func (repo *AccountCmdRepo) UpdatePassword(
	accountId valueObject.AccountId, password valueObject.Password,
) error {
	passHash, err := bcrypt.GenerateFromPassword(
		[]byte(password.String()),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}

	username, err := repo.getUsernameById(accountId)
	if err != nil {
		return err
	}

	_, err = infraHelper.RunCmd("usermod", "-p", string(passHash), username.String())
	if err != nil {
		return errors.New("UserModFailed: " + err.Error())
	}

	accModel := dbModel.Account{ID: accountId.Uint64()}
	return repo.persistentDbSvc.Handler.
		Model(&accModel).
		Update("updated_at", time.Now()).
		Error
}

func (repo *AccountCmdRepo) UpdateApiKey(
	accountId valueObject.AccountId,
) (valueObject.AccessTokenValue, error) {
	uuid := uuid.New()
	secretKey := os.Getenv("ACC_API_KEY_SECRET")
	apiKeyPlainText := accountId.String() + ":" + uuid.String()

	encryptedApiKey, err := infraHelper.EncryptStr(secretKey, apiKeyPlainText)
	if err != nil {
		return "", err
	}

	apiKey, err := valueObject.NewAccessTokenValue(encryptedApiKey)
	if err != nil {
		return "", err
	}

	uuidHash := infraHelper.GenStrongHash(uuid.String())

	accModel := dbModel.Account{ID: accountId.Uint64()}
	updateResult := repo.persistentDbSvc.Handler.Model(&accModel).
		Update("key_hash", uuidHash)
	if updateResult.Error != nil {
		return "", err
	}

	return apiKey, nil
}

func (repo *AccountCmdRepo) updateQuotaTable(
	tableName string,
	accountId valueObject.AccountId,
	quota valueObject.AccountQuota,
) error {
	updateMap := map[string]interface{}{}

	if quota.Millicores.Uint() > 0 {
		updateMap["cpu_cores"] = quota.Millicores.Uint()
	}

	if quota.MemoryBytes.Read() >= 0 {
		updateMap["memory_bytes"] = uint64(quota.MemoryBytes.Read())
	}

	if quota.StorageBytes.Read() >= 0 {
		updateMap["disk_bytes"] = uint64(quota.StorageBytes.Read())
	}

	updateMap["inodes"] = quota.StorageInodes

	return repo.persistentDbSvc.Handler.
		Table(tableName).
		Where("account_id = ?", accountId.Uint64()).
		Updates(updateMap).Error
}

func (repo *AccountCmdRepo) UpdateQuota(
	accountId valueObject.AccountId,
	quota valueObject.AccountQuota,
) error {
	err := repo.updateFilesystemQuota(accountId, quota)
	if err != nil {
		return err
	}

	return repo.updateQuotaTable(
		dbModel.AccountQuota{}.TableName(),
		accountId, quota,
	)
}

func (repo *AccountCmdRepo) getStorageUsage(
	accountId valueObject.AccountId,
) (quotaUsage valueObject.AccountQuota, err error) {
	xfsReportUsage, err := infraHelper.RunCmdWithSubShell(
		"xfs_quota -x -c 'report -nbiN' /var/data | awk '/#" +
			accountId.String() + " /{print $1, $2, $7; exit;}'",
	)
	if err != nil {
		return quotaUsage, err
	}

	if xfsReportUsage == "" {
		return quotaUsage, errors.New("InvalidXfsReportUsage")
	}

	xfsReportUsage = strings.TrimSpace(xfsReportUsage)
	if xfsReportUsage == "" {
		return quotaUsage, errors.New("InvalidXfsReportUsage")
	}

	usageColumns := strings.Split(xfsReportUsage, " ")
	if len(usageColumns) < 3 {
		return quotaUsage, errors.New("InvalidXfsReportUsage")
	}

	storageUsageKilobytesStr := usageColumns[1]
	inodesUsageStr := usageColumns[2]
	if storageUsageKilobytesStr == "" || inodesUsageStr == "" {
		return quotaUsage, errors.New("InvalidXfsReportUsage")
	}

	storageUsageBytesStr := storageUsageKilobytesStr + "000"
	storageUsage, err := valueObject.NewByte(storageUsageBytesStr)
	if err != nil {
		return quotaUsage, err
	}

	inodesUsage, err := strconv.ParseUint(inodesUsageStr, 10, 64)
	if err != nil {
		return quotaUsage, err
	}

	millicores, _ := valueObject.NewMillicores(0)
	memoryBytes, _ := valueObject.NewByte(0)
	storagePerformanceUnits, _ := valueObject.NewStoragePerformanceUnits(0)

	return valueObject.NewAccountQuota(
		millicores, memoryBytes, storageUsage, inodesUsage, storagePerformanceUnits,
	), nil
}

func (repo *AccountCmdRepo) UpdateQuotaUsage(accountId valueObject.AccountId) error {
	storageUsage, err := repo.getStorageUsage(accountId)
	if err != nil {
		return err
	}

	containerQueryRepo := NewContainerQueryRepo(repo.persistentDbSvc)
	containers, err := containerQueryRepo.ReadByAccountId(accountId)
	if err != nil {
		return err
	}
	millicoresUsage := uint(0)
	memoryBytesUsage := int64(0)
	storagePerformanceUnitsUsage := uint(0)

	profileQueryRepo := NewContainerProfileQueryRepo(repo.persistentDbSvc)
	profileIdProfileEntityMap := map[valueObject.ContainerProfileId]entity.ContainerProfile{}
	for _, container := range containers {
		if _, exists := profileIdProfileEntityMap[container.ProfileId]; exists {
			continue
		}

		profileEntity, err := profileQueryRepo.ReadById(container.ProfileId)
		if err != nil {
			slog.Debug(
				"ReadProfileByIdError",
				slog.Uint64("profileId", container.ProfileId.Read()),
				slog.Any("error", err),
			)
			continue
		}

		profileIdProfileEntityMap[container.ProfileId] = profileEntity
	}

	for _, container := range containers {
		profileEntity, exists := profileIdProfileEntityMap[container.ProfileId]
		if !exists {
			slog.Debug(
				"ProfileNotFoundForContainer",
				slog.String("containerId", container.Id.String()),
			)
			continue
		}

		containerMillicores := profileEntity.BaseSpecs.Millicores.Uint()
		containerMemoryBytes := profileEntity.BaseSpecs.MemoryBytes.Read()
		storagePerformanceUnits := profileEntity.BaseSpecs.StoragePerformanceUnits.Uint()

		millicoresUsage += containerMillicores
		memoryBytesUsage += containerMemoryBytes
		storagePerformanceUnitsUsage += storagePerformanceUnits
	}

	millicores, _ := valueObject.NewMillicores(millicoresUsage)
	memoryBytes, _ := valueObject.NewByte(memoryBytesUsage)
	storagePerformanceUnits, _ := valueObject.NewStoragePerformanceUnits(
		storagePerformanceUnitsUsage,
	)

	quota := valueObject.NewAccountQuota(
		millicores, memoryBytes, storageUsage.StorageBytes,
		storageUsage.StorageInodes, storagePerformanceUnits,
	)

	return repo.updateQuotaTable(
		dbModel.AccountQuotaUsage{}.TableName(),
		accountId,
		quota,
	)
}
