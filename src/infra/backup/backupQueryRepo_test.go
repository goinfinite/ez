package backupInfra

import (
	"testing"

	testHelpers "github.com/goinfinite/ez/src/devUtils"
	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/useCase"
	"github.com/goinfinite/ez/src/domain/valueObject"
)

func TestBackupQueryRepo(t *testing.T) {
	testHelpers.LoadEnvVars()
	persistentDbSvc := testHelpers.GetPersistentDbSvc()
	backupQueryRepo := NewBackupQueryRepo(persistentDbSvc)

	t.Run("ReadDestination", func(t *testing.T) {
		readDto := dto.ReadBackupDestinationsRequest{
			Pagination: useCase.BackupDestinationsDefaultPagination,
		}

		responseDto, err := backupQueryRepo.ReadDestination(readDto, true)
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

	t.Run("ReadTask", func(t *testing.T) {
		containerId, _ := valueObject.NewContainerId("58837bc95af5")

		readDto := dto.ReadBackupTasksRequest{
			Pagination:  useCase.BackupTasksDefaultPagination,
			ContainerId: &containerId,
		}

		responseDto, err := backupQueryRepo.ReadTask(readDto)
		if err != nil {
			t.Errorf("ReadTaskError: %v", err)
			return
		}

		if len(responseDto.Tasks) == 0 {
			t.Errorf("NoItemsFound")
		}
	})
}
