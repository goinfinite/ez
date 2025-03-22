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
	trailDbSvc := testHelpers.GetTrailDbSvc()
	backupQueryRepo := NewBackupQueryRepo(persistentDbSvc)
	backupCmdRepo := NewBackupCmdRepo(persistentDbSvc, trailDbSvc)

	t.Run("CreateBackupDestination", func(t *testing.T) {
		accountId, _ := valueObject.NewAccountId(1000)
		destinationName, _ := valueObject.NewBackupDestinationName("test")

		createDto := dto.CreateBackupDestinationRequest{
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

	t.Run("DeleteBackupDestination", func(t *testing.T) {
		accountId, _ := valueObject.NewAccountId(1000)
		destinationId, _ := valueObject.NewBackupDestinationId(1)

		deleteDto := dto.DeleteBackupDestination{
			AccountId:     accountId,
			DestinationId: destinationId,
		}

		err := backupCmdRepo.DeleteDestination(deleteDto)
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

	t.Run("UpdateBackupJob", func(t *testing.T) {
		accountId, _ := valueObject.NewAccountId(1000)
		jobId, _ := valueObject.NewBackupJobId(1)
		newBackupSchedule, _ := valueObject.NewCronSchedule("@hourly")
		firstDestinationId, _ := valueObject.NewBackupDestinationId(1)
		secondDestinationId, _ := valueObject.NewBackupDestinationId(2)

		updateDto := dto.UpdateBackupJob{
			JobId:          jobId,
			AccountId:      accountId,
			DestinationIds: []valueObject.BackupDestinationId{firstDestinationId, secondDestinationId},
			BackupSchedule: &newBackupSchedule,
		}

		err := backupCmdRepo.UpdateJob(updateDto)
		if err != nil {
			t.Errorf("ExpectedNoErrorButGot: %v", err)
		}
	})

	t.Run("RunBackupJob", func(t *testing.T) {
		accountId, _ := valueObject.NewAccountId(1000)
		jobId, _ := valueObject.NewBackupJobId(1)

		runDto := dto.RunBackupJob{
			JobId:     jobId,
			AccountId: accountId,
		}

		_, err := backupCmdRepo.RunJob(runDto)
		if err != nil {
			t.Errorf("ExpectedNoErrorButGot: %v", err)
		}
	})

	t.Run("DeleteBackupJob", func(t *testing.T) {
		accountId, _ := valueObject.NewAccountId(1000)
		jobId, _ := valueObject.NewBackupJobId(1)

		deleteDto := dto.DeleteBackupJob{
			JobId:     jobId,
			AccountId: accountId,
		}

		err := backupCmdRepo.DeleteJob(deleteDto)
		if err != nil {
			t.Errorf("ExpectedNoErrorButGot: %v", err)
		}
	})

	t.Run("RestoreBackupTask", func(t *testing.T) {
		taskEntity, err := backupQueryRepo.ReadFirstTask(dto.ReadBackupTasksRequest{})
		if err != nil {
			t.Errorf("ExpectedNoErrorButGot: %v", err)
			return
		}

		requestRestoreDto := dto.RestoreBackupTaskRequest{
			TaskId: &taskEntity.TaskId,
		}

		_, err = backupCmdRepo.RestoreTask(requestRestoreDto)
		if err != nil {
			t.Errorf("ExpectedNoErrorButGot: %v", err)
		}
	})

	t.Run("UpdateBackupTask", func(t *testing.T) {
		taskEntity, err := backupQueryRepo.ReadFirstTask(dto.ReadBackupTasksRequest{})
		if err != nil {
			t.Errorf("ExpectedNoErrorButGot: %v", err)
			return
		}

		taskStatus := valueObject.BackupTaskStatusCancelled
		err = backupCmdRepo.UpdateTask(dto.UpdateBackupTask{
			TaskId:     taskEntity.TaskId,
			TaskStatus: &taskStatus,
		})
		if err != nil {
			t.Errorf("ExpectedNoErrorButGot: %v", err)
		}
	})

	t.Run("DeleteBackupTask", func(t *testing.T) {
		taskEntity, err := backupQueryRepo.ReadFirstTask(dto.ReadBackupTasksRequest{})
		if err != nil {
			t.Errorf("ExpectedNoErrorButGot: %v", err)
			return
		}

		deleteDo := dto.DeleteBackupTask{
			TaskId:             taskEntity.TaskId,
			ShouldDiscardFiles: true,
		}

		err = backupCmdRepo.DeleteTask(deleteDo)
		if err != nil {
			t.Errorf("ExpectedNoErrorButGot: %v", err)
		}
	})

	t.Run("CreateBackupTaskArchive", func(t *testing.T) {
		taskEntity, err := backupQueryRepo.ReadFirstTask(dto.ReadBackupTasksRequest{})
		if err != nil {
			t.Errorf("ExpectedNoErrorButGot: %v", err)
			return
		}

		timeoutSecs := uint32(60)
		createDto := dto.CreateBackupTaskArchive{
			TaskId:                    taskEntity.TaskId,
			TimeoutSecs:               &timeoutSecs,
			ContainerAccountIds:       []valueObject.AccountId{},
			ContainerIds:              []valueObject.ContainerId{},
			ExceptContainerAccountIds: []valueObject.AccountId{},
			ExceptContainerIds:        []valueObject.ContainerId{},
		}

		_, err = backupCmdRepo.CreateTaskArchive(createDto)
		if err != nil {
			t.Errorf("ExpectedNoErrorButGot: %v", err)
		}
	})

	t.Run("DeleteBackupTaskArchive", func(t *testing.T) {
		archiveEntity, err := backupQueryRepo.ReadFirstTaskArchive(dto.ReadBackupTaskArchivesRequest{})
		if err != nil {
			t.Errorf("ExpectedNoErrorButGot: %v", err)
			return
		}

		deleteDto := dto.DeleteBackupTaskArchive{
			ArchiveId: archiveEntity.ArchiveId,
		}

		err = backupCmdRepo.DeleteTaskArchive(deleteDto)
		if err != nil {
			t.Errorf("ExpectedNoErrorButGot: %v", err)
		}
	})
}
