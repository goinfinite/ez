package infra

import (
	"os"
	"testing"

	testHelpers "github.com/speedianet/control/src/devUtils"
	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/valueObject"
	infraHelper "github.com/speedianet/control/src/infra/helper"
)

func addDummyUser() error {
	username, _ := valueObject.NewUsername(os.Getenv("DUMMY_USER_NAME"))
	password, _ := valueObject.NewPassword(os.Getenv("DUMMY_USER_PASS"))
	quota := valueObject.NewAccountQuotaWithDefaultValues()
	ipAddress, _ := valueObject.NewIpAddress("127.0.0.1")
	createDto := dto.NewCreateAccount(username, password, &quota, &ipAddress)

	accountCmdRepo := NewAccountCmdRepo(testHelpers.GetPersistentDbSvc())
	err := accountCmdRepo.Create(createDto)
	if err != nil {
		return err
	}

	return nil
}

func deleteDummyUser() error {
	accountId := valueObject.NewAccountIdPanic(os.Getenv("DUMMY_USER_ID"))

	accountCmdRepo := NewAccountCmdRepo(testHelpers.GetPersistentDbSvc())
	err := accountCmdRepo.Delete(accountId)
	if err != nil {
		return err
	}

	return nil
}

func resetDummyUser() {
	_ = addDummyUser()
	_ = deleteDummyUser()
	_ = addDummyUser()
}

func TestAccountCmdRepo(t *testing.T) {
	testHelpers.LoadEnvVars()
	persistentDbSvc := testHelpers.GetPersistentDbSvc()
	accountQueryRepo := NewAccountQueryRepo(persistentDbSvc)
	accountCmdRepo := NewAccountCmdRepo(persistentDbSvc)

	t.Run("AddValidAccount", func(t *testing.T) {
		err := addDummyUser()
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("AddInvalidAccount", func(t *testing.T) {
		username, _ := valueObject.NewUsername("root")
		password, _ := valueObject.NewPassword("invalid")
		quota := valueObject.NewAccountQuotaWithDefaultValues()
		ipAddress, _ := valueObject.NewIpAddress("127.0.0.1")
		createDto := dto.NewCreateAccount(username, password, &quota, &ipAddress)

		err := accountCmdRepo.Create(createDto)
		if err == nil {
			t.Error("AccountShouldNotBeAdded")
		}
	})

	t.Run("DeleteValidAccount", func(t *testing.T) {
		err := deleteDummyUser()
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("UpdatePasswordValidAccount", func(t *testing.T) {
		resetDummyUser()

		accountId := valueObject.NewAccountIdPanic(os.Getenv("DUMMY_USER_ID"))
		newPassword, _ := valueObject.NewPassword("newPassword")

		err := accountCmdRepo.UpdatePassword(accountId, newPassword)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("UpdateApiKeyValidAccount", func(t *testing.T) {
		resetDummyUser()

		accountId := valueObject.NewAccountIdPanic(os.Getenv("DUMMY_USER_ID"))

		_, err := accountCmdRepo.UpdateApiKey(accountId)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("UpdateQuotaValidAccount", func(t *testing.T) {
		resetDummyUser()

		accountId := valueObject.NewAccountIdPanic(os.Getenv("DUMMY_USER_ID"))
		quota := valueObject.NewAccountQuotaWithDefaultValues()
		quota.CpuCores = valueObject.NewCpuCoresCountPanic(1)
		quota.DiskBytes = valueObject.NewBytePanic(1073741824)

		err := accountCmdRepo.UpdateQuota(accountId, quota)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("UpdateQuotasUsage", func(t *testing.T) {
		resetDummyUser()

		testFilePath := "/var/data/" + os.Getenv("DUMMY_USER_NAME") + "/test.file"

		_, err := infraHelper.RunCmd("fallocate", "-l", "100M", testFilePath)
		if err != nil {
			t.Error(err)
		}

		accId := valueObject.NewAccountIdPanic(os.Getenv("DUMMY_USER_ID"))
		os.Chown(testFilePath, int(accId.Get()), int(accId.Get()))

		err = accountCmdRepo.UpdateQuotaUsage(accId)
		if err != nil {
			t.Error(err)
		}

		accEntity, err := accountQueryRepo.GetById(accId)
		if err != nil {
			t.Error(err)
		}
		if accEntity.QuotaUsage.DiskBytes.Get() < 100000000 {
			t.Error("QuotaUsageNotUpdated")
		}

		_, err = infraHelper.RunCmd("rm", "-f", testFilePath)
		if err != nil {
			t.Error(err)
		}
	})
}
