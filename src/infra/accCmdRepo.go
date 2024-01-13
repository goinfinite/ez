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

type AccCmdRepo struct {
	dbSvc *db.DatabaseService
}

func NewAccCmdRepo(dbSvc *db.DatabaseService) *AccCmdRepo {
	return &AccCmdRepo{dbSvc: dbSvc}
}

func (repo AccCmdRepo) updateFilesystemQuota(
	accId valueObject.AccountId,
	quota valueObject.AccountQuota,
) error {
	diskBytesStr := quota.DiskBytes.String()
	inodesStr := quota.Inodes.String()
	accIdStr := accId.String()

	shouldUpdateDiskQuota := quota.DiskBytes.Get() > 0
	shouldUpdateInodeQuota := quota.Inodes.Get() > 0
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

func (repo AccCmdRepo) Add(addAccount dto.AddAccount) error {
	passHash, err := bcrypt.GenerateFromPassword(
		[]byte(addAccount.Password.String()),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}

	addAccountCmd := exec.Command(
		"useradd",
		"-m",
		"-d", "/var/data/"+addAccount.Username.String(),
		"-s", "/usr/sbin/nologin",
		"-p", string(passHash),
		addAccount.Username.String(),
	)

	err = addAccountCmd.Run()
	if err != nil {
		return err
	}

	userInfo, err := user.Lookup(addAccount.Username.String())
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

	nowUnixTime := valueObject.UnixTime(time.Now().Unix())
	accEntity := entity.NewAccount(
		accId,
		gid,
		addAccount.Username,
		*addAccount.Quota,
		valueObject.NewAccountQuotaWithBlankValues(),
		nowUnixTime,
		nowUnixTime,
	)

	accModel, err := dbModel.Account{}.ToModel(accEntity)
	if err != nil {
		return err
	}

	err = repo.dbSvc.Orm.Create(&accModel).Error
	if err != nil {
		return err
	}

	err = repo.updateFilesystemQuota(accId, *addAccount.Quota)
	if err != nil {
		return err
	}

	return nil
}

func (repo AccCmdRepo) getUsernameById(
	accId valueObject.AccountId,
) (valueObject.Username, error) {
	accQuery := NewAccQueryRepo(repo.dbSvc)
	accDetails, err := accQuery.GetById(accId)
	if err != nil {
		return "", err
	}

	return accDetails.Username, nil
}

func (repo AccCmdRepo) Delete(accId valueObject.AccountId) error {
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
		return err
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
	modelId := accId.Get()

	relatedTables := []string{
		dbModel.AccountQuota{}.TableName(),
		dbModel.AccountQuotaUsage{}.TableName(),
	}

	for _, tableName := range relatedTables {
		err := repo.dbSvc.Orm.Exec(
			"DELETE FROM "+tableName+" WHERE account_id = ?", modelId,
		).Error
		if err != nil {
			return errors.New("DeleteAccRelatedTablesDbError")
		}
	}

	err = repo.dbSvc.Orm.Delete(model, modelId).Error
	if err != nil {
		return errors.New("DeleteAccDbError")
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

	accModel := dbModel.Account{ID: uint(accId.Get())}
	updateResult := repo.dbSvc.Orm.Model(&accModel).
		Update("updated_at", time.Now())

	return updateResult.Error
}

func (repo AccCmdRepo) UpdateApiKey(
	accId valueObject.AccountId,
) (valueObject.AccessTokenStr, error) {
	uuid := uuid.New()
	secretKey := os.Getenv("ACC_API_KEY_SECRET")
	apiKeyPlainText := accId.String() + ":" + uuid.String()

	encryptedApiKey, err := infraHelper.EncryptStr(secretKey, apiKeyPlainText)
	if err != nil {
		return "", err
	}

	apiKey, err := valueObject.NewAccessTokenStr(encryptedApiKey)
	if err != nil {
		return "", err
	}

	hash := sha3.New256()
	hash.Write([]byte(uuid.String()))
	uuidHash := hex.EncodeToString(hash.Sum(nil))

	accModel := dbModel.Account{ID: uint(accId.Get())}
	updateResult := repo.dbSvc.Orm.Model(&accModel).
		Update("key_hash", uuidHash)
	if updateResult.Error != nil {
		return "", err
	}

	return apiKey, nil
}

func (repo AccCmdRepo) updateQuotaTable(
	tableName string,
	accId valueObject.AccountId,
	quota valueObject.AccountQuota,
) error {
	updateMap := map[string]interface{}{}

	if quota.CpuCores.Get() >= 0 {
		updateMap["cpu_cores"] = quota.CpuCores.Get()
	}

	if quota.MemoryBytes.Get() >= 0 {
		updateMap["memory_bytes"] = uint64(quota.MemoryBytes.Get())
	}

	if quota.DiskBytes.Get() >= 0 {
		updateMap["disk_bytes"] = uint64(quota.DiskBytes.Get())
	}

	if quota.Inodes.Get() >= 0 {
		updateMap["inodes"] = quota.Inodes.Get()
	}

	err := repo.dbSvc.Orm.Table(tableName).
		Where("account_id = ?", uint(accId.Get())).
		Updates(updateMap).Error
	if err != nil {
		return err
	}

	return nil
}

func (repo AccCmdRepo) UpdateQuota(
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

func (repo AccCmdRepo) getStorageUsage(
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

func (repo AccCmdRepo) UpdateQuotaUsage(accId valueObject.AccountId) error {
	storageUsage, err := repo.getStorageUsage(accId)
	if err != nil {
		return err
	}

	containerQueryRepo := NewContainerQueryRepo(repo.dbSvc)
	containers, err := containerQueryRepo.GetByAccId(accId)
	if err != nil {
		return err
	}
	cpuCoresUsage := float64(0)
	memoryBytesUsage := int64(0)

	profileQueryRepo := NewContainerProfileQueryRepo(repo.dbSvc)

	for _, container := range containers {
		containerProfile, err := profileQueryRepo.GetById(
			container.ProfileId,
		)
		if err != nil {
			continue
		}

		containerCpuCores := containerProfile.BaseSpecs.CpuCores.Get()
		containerMemoryBytes := containerProfile.BaseSpecs.MemoryBytes.Get()

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
