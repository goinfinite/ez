package infra

import (
	"os"
	"testing"

	testHelpers "github.com/goinfinite/ez/src/devUtils"
	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/valueObject"
	infraHelper "github.com/goinfinite/ez/src/infra/helper"
)

func addDummyUser() error {
	username, _ := valueObject.NewUsername(os.Getenv("DUMMY_USER_NAME"))
	password, _ := valueObject.NewPassword(os.Getenv("DUMMY_USER_PASS"))
	quota := valueObject.NewAccountQuotaWithDefaultValues()
	ipAddress := valueObject.NewLocalhostIpAddress()
	operatorAccountId, _ := valueObject.NewAccountId(0)
	createDto := dto.NewCreateAccount(
		username, password, &quota, operatorAccountId, ipAddress,
	)

	accountCmdRepo := NewAccountCmdRepo(testHelpers.GetPersistentDbSvc())
	_, err := accountCmdRepo.Create(createDto)
	if err != nil {
		return err
	}

	return nil
}

func deleteDummyUser() error {
	accountId, _ := valueObject.NewAccountId(os.Getenv("DUMMY_USER_ID"))
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
	accountId, _ := valueObject.NewAccountId(os.Getenv("DUMMY_USER_ID"))

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
		ipAddress := valueObject.NewLocalhostIpAddress()
		operatorAccountId, _ := valueObject.NewAccountId(0)
		createDto := dto.NewCreateAccount(
			username, password, &quota, operatorAccountId, ipAddress,
		)

		_, err := accountCmdRepo.Create(createDto)
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

		newPassword, _ := valueObject.NewPassword("newPassword")

		err := accountCmdRepo.UpdatePassword(accountId, newPassword)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("UpdateApiKeyValidAccount", func(t *testing.T) {
		resetDummyUser()

		_, err := accountCmdRepo.UpdateApiKey(accountId)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("UpdateQuotaValidAccount", func(t *testing.T) {
		resetDummyUser()

		quota := valueObject.NewAccountQuotaWithDefaultValues()
		quota.Millicores, _ = valueObject.NewMillicores(1000)
		quota.StorageBytes, _ = valueObject.NewByte(1073741824)

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

		os.Chown(testFilePath, int(accountId.Uint64()), int(accountId.Uint64()))

		err = accountCmdRepo.UpdateQuotaUsage(accountId)
		if err != nil {
			t.Error(err)
		}

		accountEntity, err := accountQueryRepo.ReadById(accountId)
		if err != nil {
			t.Error(err)
		}
		if accountEntity.QuotaUsage.StorageBytes.Int64() < 100000000 {
			t.Error("QuotaUsageNotUpdated")
		}

		_, err = infraHelper.RunCmd("rm", "-f", testFilePath)
		if err != nil {
			t.Error(err)
		}
	})
}
