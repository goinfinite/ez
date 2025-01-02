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

	t.Run("UpdateBackupDestination", func(t *testing.T) {
		accountId, _ := valueObject.NewAccountId(1000)
		destinationId, _ := valueObject.NewBackupDestinationId(1)
		newDestinationName, _ := valueObject.NewBackupDestinationName("test2")

		updateDto := dto.UpdateBackupDestination{
			AccountId:       accountId,
			DestinationId:   destinationId,
			DestinationName: &newDestinationName,
		}

		err := backupCmdRepo.UpdateDestination(updateDto)
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
