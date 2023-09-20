package infra

import (
	"encoding/hex"
	"os"
	"os/exec"
	"os/user"
	"strings"
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
		return err
	}

	err = dbSvc.Create(&accModel).Error
	if err != nil {
		return err
	}

	err = repo.updateFilesystemQuota(accId, addAccount.Quota)
	if err != nil {
		return err
	}

	return nil
}

func getUsernameById(accId valueObject.AccountId) (valueObject.Username, error) {
	accQuery := AccQueryRepo{}
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

	username, err := getUsernameById(accId)
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

	dbSvc, err := db.DatabaseService()
	if err != nil {
		return err
	}

	err = dbModel.Account{ID: uint(accId.Get())}.Delete(dbSvc)
	if err != nil {
		return err
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
		return err
	}

	dbSvc, err := db.DatabaseService()
	if err != nil {
		return err
	}

	err = dbSvc.Model(&dbModel.Account{ID: uint(accId.Get())}).
		Update("updated_at", time.Now()).Error
	if err != nil {
		return err
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
		return "", err
	}

	apiKey, err := valueObject.NewAccessTokenStr(encryptedApiKey)
	if err != nil {
		return "", err
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
		return "", err
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

func (repo AccCmdRepo) UpdateQuotasUsage() error {
	xfsQuotasUsage, err := infraHelper.RunCmd(
		"bash",
		"-c",
		"xfs_quota -x -c 'report -nbiN' /var/data | awk '{print $1, $2, $7}'",
	)
	if err != nil {
		return err
	}

	if xfsQuotasUsage == "" {
		return nil
	}

	quotasUsageLines := strings.Split(xfsQuotasUsage, "\n")
	for _, quotasUsageLine := range quotasUsageLines {
		quotasUsageLine = strings.TrimSpace(quotasUsageLine)
		if quotasUsageLine == "" {
			continue
		}
		usageColumns := strings.Split(quotasUsageLine, " ")
		accIdStr := usageColumns[0]
		if accIdStr == "#0" || accIdStr == "" {
			continue
		}
		accIdStrWithoutHash := strings.Replace(accIdStr, "#", "", 1)
		accId, err := valueObject.NewAccountId(accIdStrWithoutHash)
		if err != nil {
			continue
		}

		diskUsageKilobytesStr := usageColumns[1]
		inodesUsageStr := usageColumns[2]
		if diskUsageKilobytesStr == "" || inodesUsageStr == "" {
			continue
		}

		diskUsageBytesStr := diskUsageKilobytesStr + "000"
		diskUsage, err := valueObject.NewByte(diskUsageBytesStr)
		if err != nil {
			continue
		}

		inodesUsage, err := valueObject.NewInodesCount(inodesUsageStr)
		if err != nil {
			continue
		}

		quota := valueObject.NewAccountQuota(
			valueObject.NewCpuCoresCountPanic(0),
			valueObject.NewBytePanic(0),
			diskUsage,
			inodesUsage,
		)

		err = repo.UpdateQuotaUsage(accId, quota)
		if err != nil {
			continue
		}
	}

	return nil
}
