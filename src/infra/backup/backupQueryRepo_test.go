package backupInfra

import (
	"testing"

	testHelpers "github.com/goinfinite/ez/src/devUtils"
	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/useCase"
)

func TestBackupQueryRepo(t *testing.T) {
	testHelpers.LoadEnvVars()
	persistentDbSvc := testHelpers.GetPersistentDbSvc()
	backupQueryRepo := NewBackupQueryRepo(persistentDbSvc)

	t.Run("ReadDestination", func(t *testing.T) {
		readDto := dto.ReadBackupDestinationsRequest{
			Pagination: useCase.BackupDestinationsDefaultPagination,
		}

		responseDto, err := backupQueryRepo.ReadDestination(readDto)
		if err != nil {
			t.Errorf("ReadDestinationError: %v", err)
			return
		}

		if len(responseDto.Destinations) == 0 {
			t.Errorf("NoItemsFound")
		}
	})

	t.Run("ReadJob", func(t *testing.T) {
		readDto := dto.ReadBackupJobsRequest{
			Pagination: useCase.BackupJobsDefaultPagination,
		}

		responseDto, err := backupQueryRepo.ReadJob(readDto)
		if err != nil {
			t.Errorf("ReadJobError: %v", err)
			return
		}

		if len(responseDto.Jobs) == 0 {
			t.Errorf("NoItemsFound")
		}
	})
}
