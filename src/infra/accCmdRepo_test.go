package infra

import (
	"os"
	"testing"

	testHelpers "github.com/speedianet/sfm/src/devUtils"
	"github.com/speedianet/sfm/src/domain/dto"
	"github.com/speedianet/sfm/src/domain/valueObject"
)

func addDummyUser() error {
	username := valueObject.NewUsernamePanic(os.Getenv("DUMMY_USER_NAME"))
	password := valueObject.NewPasswordPanic(os.Getenv("DUMMY_USER_PASS"))
	quota := valueObject.NewAccountQuotaWithDefaultValues()
	addAccount := dto.NewAddAccount(username, password, quota)

	accCmdRepo := AccCmdRepo{}
	err := accCmdRepo.Add(addAccount)
	if err != nil {
		return err
	}

	return nil
}

func deleteDummyUser() error {
	accountId := valueObject.NewAccountIdPanic(os.Getenv("DUMMY_USER_ID"))

	accCmdRepo := AccCmdRepo{}
	err := accCmdRepo.Delete(accountId)
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

func TestAccCmdRepo(t *testing.T) {
	testHelpers.LoadEnvVars()

	t.Run("AddValidAccount", func(t *testing.T) {
		err := addDummyUser()
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("AddInvalidAccount", func(t *testing.T) {
		username := valueObject.NewUsernamePanic("root")
		password := valueObject.NewPasswordPanic("invalid")
		quota := valueObject.NewAccountQuotaWithDefaultValues()
		addAccount := dto.NewAddAccount(username, password, quota)

		accCmdRepo := AccCmdRepo{}
		err := accCmdRepo.Add(addAccount)
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
		newPassword := valueObject.NewPasswordPanic("newPassword")

		accCmdRepo := AccCmdRepo{}
		err := accCmdRepo.UpdatePassword(accountId, newPassword)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("UpdateApiKeyValidAccount", func(t *testing.T) {
		resetDummyUser()

		accountId := valueObject.NewAccountIdPanic(os.Getenv("DUMMY_USER_ID"))

		accCmdRepo := AccCmdRepo{}
		_, err := accCmdRepo.UpdateApiKey(accountId)
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

		accCmdRepo := AccCmdRepo{}
		err := accCmdRepo.UpdateQuota(accountId, quota)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})
}
