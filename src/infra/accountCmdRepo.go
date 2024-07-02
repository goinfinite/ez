package infra

import (
	"encoding/hex"
	"errors"
	"os"
	"os/exec"
	"os/user"
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
	"golang.org/x/crypto/sha3"
)

type AccountCmdRepo struct {
	persistentDbSvc *db.PersistentDatabaseService
}

func NewAccountCmdRepo(persistentDbSvc *db.PersistentDatabaseService) *AccountCmdRepo {
	return &AccountCmdRepo{persistentDbSvc: persistentDbSvc}
}

func (repo *AccountCmdRepo) updateFilesystemQuota(
	accId valueObject.AccountId,
	quota valueObject.AccountQuota,
) error {
	diskBytesStr := quota.DiskBytes.String()
	inodesStr := quota.Inodes.String()
	accIdStr := accId.String()

	shouldUpdateDiskQuota := quota.DiskBytes.Read() > 0
	shouldUpdateInodeQuota := quota.Inodes.Read() > 0
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
	xfsFlags += " " + accIdStr + "' /var/data"

	_, err := infraHelper.RunCmd("bash", "-c", "xfs_quota "+xfsFlags)
	if err != nil {
		return err
	}

	return nil
}

func (repo *AccountCmdRepo) Create(createDto dto.CreateAccount) error {
	passHash, err := bcrypt.GenerateFromPassword(
		[]byte(createDto.Password.String()),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}

	createAccountCmd := exec.Command(
		"useradd",
		"-m",
		"-d", "/var/data/"+createDto.Username.String(),
		"-s", "/usr/sbin/nologin",
		"-p", string(passHash),
		createDto.Username.String(),
	)

	err = createAccountCmd.Run()
	if err != nil {
		return err
	}

	userInfo, err := user.Lookup(createDto.Username.String())
	if err != nil {
		return err
	}
	accId, err := valueObject.NewAccountId(userInfo.Uid)
	if err != nil {
		return err
	}
	gid, err := valueObject.NewGroupId(userInfo.Gid)
	if err != nil {
		return err
	}

	nowUnixTime := valueObject.NewUnixTimeNow()
	accEntity := entity.NewAccount(
		accId,
		gid,
		createDto.Username,
		*createDto.Quota,
		valueObject.NewAccountQuotaWithBlankValues(),
		nowUnixTime,
		nowUnixTime,
	)

	accModel, err := dbModel.Account{}.ToModel(accEntity)
	if err != nil {
		return err
	}

	err = repo.persistentDbSvc.Handler.Create(&accModel).Error
	if err != nil {
		return err
	}

	err = repo.updateFilesystemQuota(accId, *createDto.Quota)
	if err != nil {
		return err
	}

	err = infraHelper.EnableLingering(accId)
	if err != nil {
		return errors.New("EnableLingeringFailed: " + err.Error())
	}

	return nil
}

func (repo *AccountCmdRepo) getUsernameById(
	accId valueObject.AccountId,
) (valueObject.Username, error) {
	accQuery := NewAccountQueryRepo(repo.persistentDbSvc)
	accDetails, err := accQuery.ReadById(accId)
	if err != nil {
		return "", err
	}

	return accDetails.Username, nil
}

func (repo *AccountCmdRepo) Delete(accId valueObject.AccountId) error {
	quota := valueObject.NewAccountQuotaWithBlankValues()
	err := repo.updateFilesystemQuota(accId, quota)
	if err != nil {
		return err
	}

	username, err := repo.getUsernameById(accId)
	if err != nil {
		return err
	}

	err = infraHelper.DisableLingering(accId)
	if err != nil {
		return errors.New("DisableLingeringFailed: " + err.Error())
	}

	_, err = infraHelper.RunCmd("pgrep", "-u", accId.String())
	if err == nil {
		_, _ = infraHelper.RunCmd("pkill", "-9", "-U", accId.String())
	}

	_, err = infraHelper.RunCmd("userdel", "-r", username.String())
	if err != nil {
		return err
	}

	model := dbModel.Account{}
	modelId := accId.Read()

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
	accId valueObject.AccountId,
	password valueObject.Password,
) error {
	passHash, err := bcrypt.GenerateFromPassword(
		[]byte(password.String()),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}

	username, err := repo.getUsernameById(accId)
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
		return err
	}

	accModel := dbModel.Account{ID: uint(accId.Read())}
	updateResult := repo.persistentDbSvc.Handler.Model(&accModel).
		Update("updated_at", time.Now())

	return updateResult.Error
}

func (repo *AccountCmdRepo) UpdateApiKey(
	accId valueObject.AccountId,
) (valueObject.AccessTokenValue, error) {
	uuid := uuid.New()
	secretKey := os.Getenv("ACC_API_KEY_SECRET")
	apiKeyPlainText := accId.String() + ":" + uuid.String()

	encryptedApiKey, err := infraHelper.EncryptStr(secretKey, apiKeyPlainText)
	if err != nil {
		return "", err
	}

	apiKey, err := valueObject.NewAccessTokenValue(encryptedApiKey)
	if err != nil {
		return "", err
	}

	hash := sha3.New256()
	hash.Write([]byte(uuid.String()))
	uuidHash := hex.EncodeToString(hash.Sum(nil))

	accModel := dbModel.Account{ID: uint(accId.Read())}
	updateResult := repo.persistentDbSvc.Handler.Model(&accModel).
		Update("key_hash", uuidHash)
	if updateResult.Error != nil {
		return "", err
	}

	return apiKey, nil
}

func (repo *AccountCmdRepo) updateQuotaTable(
	tableName string,
	accId valueObject.AccountId,
	quota valueObject.AccountQuota,
) error {
	updateMap := map[string]interface{}{}

	if quota.CpuCores.Read() >= 0 {
		updateMap["cpu_cores"] = quota.CpuCores.Read()
	}

	if quota.MemoryBytes.Read() >= 0 {
		updateMap["memory_bytes"] = uint64(quota.MemoryBytes.Read())
	}

	if quota.DiskBytes.Read() >= 0 {
		updateMap["disk_bytes"] = uint64(quota.DiskBytes.Read())
	}

	updateMap["inodes"] = quota.Inodes.Read()

	err := repo.persistentDbSvc.Handler.Table(tableName).
		Where("account_id = ?", uint(accId.Read())).
		Updates(updateMap).Error
	if err != nil {
		return err
	}

	return nil
}

func (repo *AccountCmdRepo) UpdateQuota(
	accId valueObject.AccountId,
	quota valueObject.AccountQuota,
) error {
	err := repo.updateFilesystemQuota(accId, quota)
	if err != nil {
		return err
	}

	return repo.updateQuotaTable(
		dbModel.AccountQuota{}.TableName(),
		accId,
		quota,
	)
}

func (repo *AccountCmdRepo) getStorageUsage(
	accId valueObject.AccountId,
) (valueObject.AccountQuota, error) {
	var quotaUsage valueObject.AccountQuota

	xfsReportUsage, err := infraHelper.RunCmd(
		"bash",
		"-c",
		"xfs_quota -x -c 'report -nbiN' /var/data | awk '/#"+
			accId.String()+" /{print $1, $2, $7; exit;}'",
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

	diskUsageKilobytesStr := usageColumns[1]
	inodesUsageStr := usageColumns[2]
	if diskUsageKilobytesStr == "" || inodesUsageStr == "" {
		return quotaUsage, errors.New("InvalidXfsReportUsage")
	}

	diskUsageBytesStr := diskUsageKilobytesStr + "000"
	diskUsage, err := valueObject.NewByte(diskUsageBytesStr)
	if err != nil {
		return quotaUsage, err
	}

	inodesUsage, err := valueObject.NewInodesCount(inodesUsageStr)
	if err != nil {
		return quotaUsage, err
	}

	cpuCores, _ := valueObject.NewCpuCoresCount(0)
	memoryBytes, _ := valueObject.NewByte(0)

	return valueObject.NewAccountQuota(
		cpuCores,
		memoryBytes,
		diskUsage,
		inodesUsage,
	), nil
}

func (repo *AccountCmdRepo) UpdateQuotaUsage(accId valueObject.AccountId) error {
	storageUsage, err := repo.getStorageUsage(accId)
	if err != nil {
		return err
	}

	containerQueryRepo := NewContainerQueryRepo(repo.persistentDbSvc)
	containers, err := containerQueryRepo.ReadByAccountId(accId)
	if err != nil {
		return err
	}
	cpuCoresUsage := float64(0)
	memoryBytesUsage := int64(0)

	profileQueryRepo := NewContainerProfileQueryRepo(repo.persistentDbSvc)

	for _, container := range containers {
		containerProfile, err := profileQueryRepo.ReadById(
			container.ProfileId,
		)
		if err != nil {
			continue
		}

		containerCpuCores := containerProfile.BaseSpecs.CpuCores.Read()
		containerMemoryBytes := containerProfile.BaseSpecs.MemoryBytes.Read()

		cpuCoresUsage += containerCpuCores
		memoryBytesUsage += containerMemoryBytes
	}

	cpuCores, _ := valueObject.NewCpuCoresCount(cpuCoresUsage)
	memoryBytes, _ := valueObject.NewByte(memoryBytesUsage)

	quota := valueObject.NewAccountQuota(
		cpuCores,
		memoryBytes,
		storageUsage.DiskBytes,
		storageUsage.Inodes,
	)
	return repo.updateQuotaTable(
		dbModel.AccountQuotaUsage{}.TableName(),
		accId,
		quota,
	)
}
