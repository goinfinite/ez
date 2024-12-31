package backupInfra

import (
	"testing"

	testHelpers "github.com/goinfinite/ez/src/devUtils"
	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/valueObject"
)

func TestBackupCmdRepo(t *testing.T) {
	testHelpers.LoadEnvVars()
	persistentDbSvc := testHelpers.GetPersistentDbSvc()
	backupCmdRepo := NewBackupCmdRepo(persistentDbSvc)

	t.Run("CreateBackupDestination", func(t *testing.T) {
		accountId, _ := valueObject.NewAccountId(1000)
		destinationName, _ := valueObject.NewBackupDestinationName("test")

		createDto := dto.CreateBackupDestination{
			AccountId:       accountId,
			DestinationName: destinationName,
			DestinationType: valueObject.BackupDestinationTypeObjectStorage,
		}

		_, err := backupCmdRepo.CreateDestination(createDto)
		if err != nil {
			t.Errorf("ExpectedNoErrorButGot: %v", err)
		}
	})

	t.Run("CreateBackupJob", func(t *testing.T) {
		accountId, _ := valueObject.NewAccountId(1000)
		destinationId, _ := valueObject.NewBackupDestinationId(1)
		backupSchedule, _ := valueObject.NewCronSchedule("@daily")

		createDto := dto.CreateBackupJob{
			AccountId:      accountId,
			DestinationIds: []valueObject.BackupDestinationId{destinationId},
			BackupSchedule: backupSchedule,
		}

		_, err := backupCmdRepo.CreateJob(createDto)
		if err != nil {
			t.Errorf("ExpectedNoErrorButGot: %v", err)
		}
	})
}
