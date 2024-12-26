package service

import (
	"errors"

	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/useCase"
	"github.com/goinfinite/ez/src/domain/valueObject"
	voHelper "github.com/goinfinite/ez/src/domain/valueObject/helper"
	"github.com/goinfinite/ez/src/infra"
	backupInfra "github.com/goinfinite/ez/src/infra/backup"
	"github.com/goinfinite/ez/src/infra/db"
	serviceHelper "github.com/goinfinite/ez/src/presentation/service/helper"
)

type BackupService struct {
	persistentDbSvc       *db.PersistentDatabaseService
	backupQueryRepo       *backupInfra.BackupQueryRepo
	activityRecordCmdRepo *infra.ActivityRecordCmdRepo
}

func NewBackupService(
	persistentDbSvc *db.PersistentDatabaseService,
	trailDbSvc *db.TrailDatabaseService,
) *BackupService {
	return &BackupService{
		persistentDbSvc:       persistentDbSvc,
		backupQueryRepo:       backupInfra.NewBackupQueryRepo(persistentDbSvc),
		activityRecordCmdRepo: infra.NewActivityRecordCmdRepo(trailDbSvc),
	}
}

func (service *BackupService) ReadDestination(input map[string]interface{}) ServiceOutput {
	var destinationIdPtr *valueObject.BackupDestinationId
	if input["destinationId"] != nil {
		destinationId, err := valueObject.NewBackupDestinationId(input["destinationId"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		destinationIdPtr = &destinationId
	}

	var accountIdPtr *valueObject.AccountId
	if input["accountId"] != nil {
		accountId, err := valueObject.NewAccountId(input["accountId"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		accountIdPtr = &accountId
	}

	var destinationNamePtr *valueObject.BackupDestinationName
	if input["destinationName"] != nil {
		destinationName, err := valueObject.NewBackupDestinationName(input["destinationName"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		destinationNamePtr = &destinationName
	}

	var destinationTypePtr *valueObject.BackupDestinationType
	if input["destinationType"] != nil {
		destinationType, err := valueObject.NewBackupDestinationType(input["destinationType"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		destinationTypePtr = &destinationType
	}

	var objectStorageProviderPtr *valueObject.ObjectStorageProvider
	if input["objectStorageProvider"] != nil {
		objectStorageProvider, err := valueObject.NewObjectStorageProvider(input["objectStorageProvider"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		objectStorageProviderPtr = &objectStorageProvider
	}

	var remoteHostnamePtr *valueObject.Fqdn
	if input["remoteHostname"] != nil {
		remoteHostname, err := valueObject.NewFqdn(input["remoteHostname"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		remoteHostnamePtr = &remoteHostname
	}

	var remoteHostTypePtr *valueObject.BackupDestinationRemoteHostType
	if input["remoteHostType"] != nil {
		remoteHostType, err := valueObject.NewBackupDestinationRemoteHostType(input["remoteHostType"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		remoteHostTypePtr = &remoteHostType
	}

	timeParamNames := []string{"createdBeforeAt", "createdAfterAt"}
	timeParamPtrs := serviceHelper.TimeParamsParser(timeParamNames, input)

	requestPagination, err := serviceHelper.PaginationParser(
		input, useCase.BackupDestinationsDefaultPagination,
	)
	if err != nil {
		return NewServiceOutput(UserError, err)
	}

	readDto := dto.ReadBackupDestinationsRequest{
		Pagination:            requestPagination,
		DestinationId:         destinationIdPtr,
		AccountId:             accountIdPtr,
		DestinationName:       destinationNamePtr,
		DestinationType:       destinationTypePtr,
		ObjectStorageProvider: objectStorageProviderPtr,
		RemoteHostType:        remoteHostTypePtr,
		RemoteHostname:        remoteHostnamePtr,
		CreatedBeforeAt:       timeParamPtrs["createdBeforeAt"],
		CreatedAfterAt:        timeParamPtrs["createdAfterAt"],
	}

	responseDto, err := useCase.ReadBackupDestinations(service.backupQueryRepo, readDto)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, responseDto)
}

func (service *BackupService) ReadJob(input map[string]interface{}) ServiceOutput {
	var jobIdPtr *valueObject.BackupJobId
	if input["jobId"] != nil {
		jobId, err := valueObject.NewBackupJobId(input["jobId"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		jobIdPtr = &jobId
	}

	var jobStatusPtr *bool
	if input["jobStatus"] != nil {
		jobStatus, err := voHelper.InterfaceToBool(input["jobStatus"])
		if err != nil {
			return NewServiceOutput(UserError, errors.New("InvalidJobStatus"))
		}
		jobStatusPtr = &jobStatus
	}

	var accountIdPtr *valueObject.AccountId
	if input["accountId"] != nil {
		accountId, err := valueObject.NewAccountId(input["accountId"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		accountIdPtr = &accountId
	}

	var destinationIdPtr *valueObject.BackupDestinationId
	if input["destinationId"] != nil {
		destinationId, err := valueObject.NewBackupDestinationId(input["destinationId"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		destinationIdPtr = &destinationId
	}

	var retentionStrategyPtr *valueObject.BackupRetentionStrategy
	if input["retentionStrategy"] != nil {
		retentionStrategy, err := valueObject.NewBackupRetentionStrategy(input["retentionStrategy"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		retentionStrategyPtr = &retentionStrategy
	}

	var archiveCompressionFormatPtr *valueObject.CompressionFormat
	if input["archiveCompressionFormat"] != nil {
		archiveCompressionFormat, err := valueObject.NewCompressionFormat(input["archiveCompressionFormat"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		archiveCompressionFormatPtr = &archiveCompressionFormat
	}

	var lastRunStatusPtr *valueObject.BackupTaskStatus
	if input["lastRunStatus"] != nil {
		lastRunStatus, err := valueObject.NewBackupTaskStatus(input["lastRunStatus"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		lastRunStatusPtr = &lastRunStatus
	}

	timeParamNames := []string{
		"lastRunBeforeAt", "lastRunAfterAt", "nextRunBeforeAt", "nextRunAfterAt",
		"createdBeforeAt", "createdAfterAt",
	}
	timeParamPtrs := serviceHelper.TimeParamsParser(timeParamNames, input)

	requestPagination, err := serviceHelper.PaginationParser(
		input, useCase.BackupJobsDefaultPagination,
	)
	if err != nil {
		return NewServiceOutput(UserError, err)
	}

	readDto := dto.ReadBackupJobsRequest{
		Pagination:               requestPagination,
		JobId:                    jobIdPtr,
		JobStatus:                jobStatusPtr,
		AccountId:                accountIdPtr,
		DestinationId:            destinationIdPtr,
		RetentionStrategy:        retentionStrategyPtr,
		ArchiveCompressionFormat: archiveCompressionFormatPtr,
		LastRunStatus:            lastRunStatusPtr,
		LastRunBeforeAt:          timeParamPtrs["lastRunBeforeAt"],
		LastRunAfterAt:           timeParamPtrs["lastRunAfterAt"],
		NextRunBeforeAt:          timeParamPtrs["nextRunBeforeAt"],
		NextRunAfterAt:           timeParamPtrs["nextRunAfterAt"],
		CreatedBeforeAt:          timeParamPtrs["createdBeforeAt"],
		CreatedAfterAt:           timeParamPtrs["createdAfterAt"],
	}

	responseDto, err := useCase.ReadBackupJobs(service.backupQueryRepo, readDto)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, responseDto)
}

func (service *BackupService) ReadTask(input map[string]interface{}) ServiceOutput {
	var taskIdPtr *valueObject.BackupTaskId
	if input["taskId"] != nil {
		taskId, err := valueObject.NewBackupTaskId(input["taskId"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		taskIdPtr = &taskId
	}

	var accountIdPtr *valueObject.AccountId
	if input["accountId"] != nil {
		accountId, err := valueObject.NewAccountId(input["accountId"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		accountIdPtr = &accountId
	}

	var jobIdPtr *valueObject.BackupJobId
	if input["jobId"] != nil {
		jobId, err := valueObject.NewBackupJobId(input["jobId"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		jobIdPtr = &jobId
	}

	var destinationIdPtr *valueObject.BackupDestinationId
	if input["destinationId"] != nil {
		destinationId, err := valueObject.NewBackupDestinationId(input["destinationId"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		destinationIdPtr = &destinationId
	}

	var taskStatusPtr *valueObject.BackupTaskStatus
	if input["taskStatus"] != nil {
		taskStatus, err := valueObject.NewBackupTaskStatus(input["taskStatus"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		taskStatusPtr = &taskStatus
	}

	var retentionStrategyPtr *valueObject.BackupRetentionStrategy
	if input["retentionStrategy"] != nil {
		retentionStrategy, err := valueObject.NewBackupRetentionStrategy(input["retentionStrategy"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		retentionStrategyPtr = &retentionStrategy
	}

	var containerIdPtr *valueObject.ContainerId
	if input["containerId"] != nil {
		containerId, err := valueObject.NewContainerId(input["containerId"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		containerIdPtr = &containerId
	}

	timeParamNames := []string{
		"startedBeforeAt", "startedAfterAt", "finishedBeforeAt", "finishedAfterAt",
		"createdBeforeAt", "createdAfterAt",
	}
	timeParamPtrs := serviceHelper.TimeParamsParser(timeParamNames, input)

	requestPagination, err := serviceHelper.PaginationParser(
		input, useCase.BackupTasksDefaultPagination,
	)
	if err != nil {
		return NewServiceOutput(UserError, err)
	}

	readDto := dto.ReadBackupTasksRequest{
		Pagination:        requestPagination,
		TaskId:            taskIdPtr,
		AccountId:         accountIdPtr,
		JobId:             jobIdPtr,
		DestinationId:     destinationIdPtr,
		TaskStatus:        taskStatusPtr,
		RetentionStrategy: retentionStrategyPtr,
		ContainerId:       containerIdPtr,
		StartedBeforeAt:   timeParamPtrs["startedBeforeAt"],
		StartedAfterAt:    timeParamPtrs["startedAfterAt"],
		FinishedBeforeAt:  timeParamPtrs["finishedBeforeAt"],
		FinishedAfterAt:   timeParamPtrs["finishedAfterAt"],
		CreatedBeforeAt:   timeParamPtrs["createdBeforeAt"],
		CreatedAfterAt:    timeParamPtrs["createdAfterAt"],
	}

	responseDto, err := useCase.ReadBackupTasks(service.backupQueryRepo, readDto)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, responseDto)
}
