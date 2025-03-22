package infra

import (
	"errors"
	"log/slog"
	"os"
	"os/user"
	"strconv"
	"strings"
	"time"

	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/entity"
	"github.com/goinfinite/ez/src/domain/valueObject"
	"github.com/goinfinite/ez/src/infra/db"
	dbModel "github.com/goinfinite/ez/src/infra/db/model"
	infraEnvs "github.com/goinfinite/ez/src/infra/envs"
	infraHelper "github.com/goinfinite/ez/src/infra/helper"
	"github.com/google/uuid"
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
	storageBytesStr := quota.StorageBytes.String()
	inodesStr := strconv.FormatUint(quota.StorageInodes, 10)
	accountIdStr := accountId.String()

	shouldUpdateDiskQuota := quota.StorageBytes.Int64() > 0
	shouldUpdateInodeQuota := quota.StorageInodes > 0
	shouldRemoveQuota := !shouldUpdateDiskQuota && !shouldUpdateInodeQuota

	xfsFlags := "-x -c 'limit -u"
	if shouldUpdateDiskQuota {
		xfsFlags += " bhard=" + storageBytesStr
	}
	if shouldUpdateInodeQuota {
		xfsFlags += " ihard=" + inodesStr
	}
	if shouldRemoveQuota {
		xfsFlags = "-x -c 'limit -u bhard=0 ihard=0"
	}

	xfsFlags += " " + accountIdStr + "' " + infraEnvs.UserDataDirectory

	_, err := infraHelper.RunCmdWithSubShell("xfs_quota " + xfsFlags)
	return err
}

func (repo *AccountCmdRepo) Create(
	createDto dto.CreateAccount,
) (accountId valueObject.AccountId, err error) {
	passHash, err := bcrypt.GenerateFromPassword(
		[]byte(createDto.Password.String()), bcrypt.DefaultCost,
	)
	if err != nil {
		return accountId, errors.New("PassHashError: " + err.Error())
	}

	usernameStr := createDto.Username.String()
	homeDirectory, err := valueObject.NewUnixFilePath(
		infraEnvs.UserDataDirectory + "/" + usernameStr,
	)
	if err != nil {
		return accountId, errors.New("DefineHomeDirectoryError: " + err.Error())
	}

	_, err = infraHelper.RunCmd(
		"useradd", "-m",
		"-d", homeDirectory.String(),
		"-s", "/usr/sbin/nologin",
		"-p", string(passHash),
		usernameStr,
	)
	if err != nil {
		return accountId, errors.New("UserAddFailed: " + err.Error())
	}

	userInfo, err := user.Lookup(usernameStr)
	if err != nil {
		return accountId, errors.New("UserLookupFailed: " + err.Error())
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
		valueObject.NewAccountQuotaWithBlankValues(), homeDirectory,
		nowUnixTime, nowUnixTime,
	)

	accountModel, err := dbModel.Account{}.ToModel(accountEntity)
	if err != nil {
		return accountId, err
	}

	err = repo.persistentDbSvc.Handler.Create(&accountModel).Error
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
) (valueObject.UnixUsername, error) {
	accountQuery := NewAccountQueryRepo(repo.persistentDbSvc)
	accountEntity, err := accountQuery.ReadById(accountId)
	if err != nil {
		return "", err
	}

	return accountEntity.Username, nil
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

	accountModel := dbModel.Account{}
	accountIdUint64 := accountId.Uint64()

	relatedTables := []string{
		dbModel.AccountQuota{}.TableName(),
		dbModel.AccountQuotaUsage{}.TableName(),
	}

	for _, tableName := range relatedTables {
		err := repo.persistentDbSvc.Handler.Exec(
			"DELETE FROM "+tableName+" WHERE account_id = ?", accountIdUint64,
		).Error
		if err != nil {
			return errors.New("DeleteAccRelatedTablesDbError: " + err.Error())
		}
	}

	err = repo.persistentDbSvc.Handler.Delete(accountModel, accountIdUint64).Error
	if err != nil {
		return errors.New("DeleteAccDbError: " + err.Error())
	}

	return nil
}

func (repo *AccountCmdRepo) UpdatePassword(
	accountId valueObject.AccountId, password valueObject.Password,
) error {
	passHash, err := bcrypt.GenerateFromPassword(
		[]byte(password.String()), bcrypt.DefaultCost,
	)
	if err != nil {
		return errors.New("PassHashError: " + err.Error())
	}

	username, err := repo.getUsernameById(accountId)
	if err != nil {
		return err
	}

	_, err = infraHelper.RunCmd("usermod", "-p", string(passHash), username.String())
	if err != nil {
		return errors.New("UserModFailed: " + err.Error())
	}

	accountModel := dbModel.Account{ID: accountId.Uint64()}
	return repo.persistentDbSvc.Handler.
		Model(&accountModel).
		Update("updated_at", time.Now()).
		Error
}

func (repo *AccountCmdRepo) UpdateApiKey(
	accountId valueObject.AccountId,
) (tokenValue valueObject.AccessTokenValue, err error) {
	uuidStr := uuid.New().String()
	secretKey := os.Getenv("ACC_API_KEY_SECRET")
	apiKeyPlainText := accountId.String() + ":" + uuidStr

	encryptedApiKey, err := infraHelper.EncryptStr(secretKey, apiKeyPlainText)
	if err != nil {
		return tokenValue, err
	}

	apiKey, err := valueObject.NewAccessTokenValue(encryptedApiKey)
	if err != nil {
		return tokenValue, err
	}

	uuidHash := infraHelper.GenStrongHash(uuidStr)

	accountModel := dbModel.Account{ID: accountId.Uint64()}
	updateResult := repo.persistentDbSvc.Handler.
		Model(&accountModel).
		Update("key_hash", uuidHash)
	if updateResult.Error != nil {
		return tokenValue, err
	}

	return apiKey, nil
}

func (repo *AccountCmdRepo) UpdateQuota(
	accountId valueObject.AccountId,
	quota valueObject.AccountQuota,
) error {
	err := repo.updateFilesystemQuota(accountId, quota)
	if err != nil {
		return err
	}

	updateMap := map[string]interface{}{}

	if quota.Millicores.Uint() > 0 {
		updateMap["millicores"] = quota.Millicores.Uint()
	}

	if quota.MemoryBytes.Int64() > 0 {
		updateMap["memory_bytes"] = uint64(quota.MemoryBytes.Int64())
	}

	if quota.StorageBytes.Int64() > 0 {
		updateMap["storage_bytes"] = uint64(quota.StorageBytes.Int64())
	}

	if quota.StorageInodes > 0 {
		updateMap["storage_inodes"] = quota.StorageInodes
	}

	if quota.StoragePerformanceUnits.Uint() > 0 {
		updateMap["storage_performance_units"] = quota.StoragePerformanceUnits.Uint()
	}

	return repo.persistentDbSvc.Handler.
		Model(&dbModel.AccountQuota{}).
		Where("account_id = ?", accountId.Uint64()).
		Updates(updateMap).Error
}

func (repo *AccountCmdRepo) getStorageUsage(
	accountId valueObject.AccountId,
) (quotaUsage valueObject.AccountQuota, err error) {
	xfsReportUsage, err := infraHelper.RunCmdWithSubShell(
		"xfs_quota -x -c 'report -nbiN' " + infraEnvs.UserDataDirectory + " | awk '/#" +
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

	readContainersPaginationDto := dto.Pagination{
		PageNumber:   0,
		ItemsPerPage: 1000,
	}
	readContainersRequestDto := dto.ReadContainersRequest{
		Pagination:         readContainersPaginationDto,
		ContainerAccountId: []valueObject.AccountId{accountId},
	}

	readContainersResponseDto, err := containerQueryRepo.Read(readContainersRequestDto)
	if err != nil {
		return err
	}

	profileQueryRepo := NewContainerProfileQueryRepo(repo.persistentDbSvc)
	profileIdProfileEntityMap := map[valueObject.ContainerProfileId]entity.ContainerProfile{}
	for _, container := range readContainersResponseDto.Containers {
		if _, exists := profileIdProfileEntityMap[container.ProfileId]; exists {
			continue
		}

		profileEntity, err := profileQueryRepo.ReadById(container.ProfileId)
		if err != nil {
			slog.Debug(
				"ReadProfileByIdError",
				slog.Uint64("profileId", container.ProfileId.Uint64()),
				slog.Any("error", err),
			)
			continue
		}

		profileIdProfileEntityMap[container.ProfileId] = profileEntity
	}

	millicoresUsage := uint(0)
	memoryBytesUsage := int64(0)
	storagePerformanceUnitsUsage := uint(0)

	for _, container := range readContainersResponseDto.Containers {
		profileEntity, exists := profileIdProfileEntityMap[container.ProfileId]
		if !exists {
			slog.Debug(
				"ProfileNotFoundForContainer",
				slog.String("containerId", container.Id.String()),
			)
			continue
		}

		containerMillicores := profileEntity.BaseSpecs.Millicores.Uint()
		containerMemoryBytes := profileEntity.BaseSpecs.MemoryBytes.Int64()
		storagePerformanceUnits := profileEntity.BaseSpecs.StoragePerformanceUnits.Uint()

		millicoresUsage += containerMillicores
		memoryBytesUsage += containerMemoryBytes
		storagePerformanceUnitsUsage += storagePerformanceUnits
	}

	millicores, err := valueObject.NewMillicores(millicoresUsage)
	if err != nil {
		return err
	}

	memoryBytes, err := valueObject.NewByte(memoryBytesUsage)
	if err != nil {
		return err
	}

	storagePerformanceUnits, _ := valueObject.NewStoragePerformanceUnits(
		storagePerformanceUnitsUsage,
	)

	updateMap := map[string]interface{}{
		"millicores":                millicores.Uint(),
		"memory_bytes":              uint64(memoryBytes.Int64()),
		"storage_bytes":             uint64(storageUsage.StorageBytes.Int64()),
		"storage_inodes":            storageUsage.StorageInodes,
		"storage_performance_units": storagePerformanceUnits.Uint(),
	}

	return repo.persistentDbSvc.Handler.
		Model(&dbModel.AccountQuotaUsage{}).
		Where("account_id = ?", accountId.Uint64()).
		Updates(updateMap).Error
}
