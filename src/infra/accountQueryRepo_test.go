package infra

import (
	"os"
	"testing"

	testHelpers "github.com/goinfinite/ez/src/devUtils"
	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/useCase"
	"github.com/goinfinite/ez/src/domain/valueObject"
)

func TestAccountQueryRepo(t *testing.T) {
	testHelpers.LoadEnvVars()
	persistentDbSvc := testHelpers.GetPersistentDbSvc()
	accountQueryRepo := NewAccountQueryRepo(persistentDbSvc)

	t.Run("GetValidAccounts", func(t *testing.T) {
		requestDto := dto.ReadAccountsRequest{
			Pagination: useCase.AccountsDefaultPagination,
		}
		_, err := accountQueryRepo.Read(requestDto)
		if err != nil {
			t.Error("UnexpectedError")
		}
	})

	t.Run("GetValidAccountByUsername", func(t *testing.T) {
		username, _ := valueObject.NewUnixUsername(os.Getenv("DUMMY_USER_NAME"))

		_, err := accountQueryRepo.ReadFirst(dto.ReadAccountsRequest{
			AccountUsername: &username,
		})
		if err != nil {
			t.Error("UnexpectedError")
		}
	})

	t.Run("GetValidAccountById", func(t *testing.T) {
		accountId, _ := valueObject.NewAccountId(os.Getenv("DUMMY_USER_ID"))

		_, err := accountQueryRepo.ReadFirst(dto.ReadAccountsRequest{
			AccountId: &accountId,
		})
		if err != nil {
			t.Error("UnexpectedError")
		}
	})

	t.Run("GetInvalidAccount", func(t *testing.T) {
		username, _ := valueObject.NewUnixUsername("invalid")

		_, err := accountQueryRepo.ReadFirst(dto.ReadAccountsRequest{
			AccountUsername: &username,
		})
		if err == nil {
			t.Error("ExpectingError")
		}
	})
}
