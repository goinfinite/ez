package useCase

import (
	"log/slog"

	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/repository"
	"github.com/goinfinite/ez/src/domain/valueObject"
)

type CreateSecurityActivityRecord struct {
	activityRecordCmdRepo repository.ActivityRecordCmdRepo
	recordLevel           valueObject.ActivityRecordLevel
}

func NewCreateSecurityActivityRecord(
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
) *CreateSecurityActivityRecord {
	recordLevel, _ := valueObject.NewActivityRecordLevel("SEC")
	return &CreateSecurityActivityRecord{
		activityRecordCmdRepo: activityRecordCmdRepo,
		recordLevel:           recordLevel,
	}
}

func (uc *CreateSecurityActivityRecord) createActivityRecord(
	createDto dto.CreateActivityRecord,
) {
	err := uc.activityRecordCmdRepo.Create(createDto)
	if err != nil {
		slog.Debug(
			"CreateSecurityActivityRecordError",
			slog.Any("createDto", createDto),
			slog.Any("error", err),
		)
	}
}

func (uc *CreateSecurityActivityRecord) CreateSessionToken(
	recordCode valueObject.ActivityRecordCode,
	createDto dto.CreateSessionToken,
) {
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel:       uc.recordLevel,
		RecordCode:        recordCode,
		RecordDetails:     createDto.Username,
		OperatorIpAddress: &createDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) CreateContainerSessionToken(
	createDto dto.CreateContainerSessionToken,
) {
	recordCode, _ := valueObject.NewActivityRecordCode("ContainerSessionTokenGenerated")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []valueObject.SystemResourceIdentifier{
			valueObject.NewContainerSri(createDto.AccountId, createDto.ContainerId),
		},
		OperatorAccountId: &createDto.OperatorAccountId,
		OperatorIpAddress: &createDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) CreateAccount(
	createDto dto.CreateAccount,
	accountId valueObject.AccountId,
) {
	recordCode, _ := valueObject.NewActivityRecordCode("AccountCreated")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []valueObject.SystemResourceIdentifier{
			valueObject.NewAccountSri(accountId),
		},
		OperatorAccountId: &createDto.OperatorAccountId,
		OperatorIpAddress: &createDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) UpdateAccount(
	updateDto dto.UpdateAccount,
) {
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		AffectedResources: []valueObject.SystemResourceIdentifier{
			valueObject.NewAccountSri(updateDto.AccountId),
		},
		RecordDetails:     updateDto,
		OperatorAccountId: &updateDto.OperatorAccountId,
		OperatorIpAddress: &updateDto.OperatorIpAddress,
	}

	codeStr := "AccountUpdated"
	if updateDto.Password != nil {
		codeStr = "AccountPasswordUpdated"
		createRecordDto.RecordDetails = nil
	}
	recordCode, _ := valueObject.NewActivityRecordCode(codeStr)
	createRecordDto.RecordCode = recordCode

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) DeleteAccount(
	deleteDto dto.DeleteAccount,
) {
	recordCode, _ := valueObject.NewActivityRecordCode("AccountDeleted")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []valueObject.SystemResourceIdentifier{
			valueObject.NewAccountSri(deleteDto.AccountId),
		},
		OperatorAccountId: &deleteDto.OperatorAccountId,
		OperatorIpAddress: &deleteDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) CreateContainer(
	createDto dto.CreateContainer,
	containerId valueObject.ContainerId,
) {
	recordCode, _ := valueObject.NewActivityRecordCode("ContainerCreated")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []valueObject.SystemResourceIdentifier{
			valueObject.NewContainerSri(createDto.AccountId, containerId),
		},
		OperatorAccountId: &createDto.OperatorAccountId,
		OperatorIpAddress: &createDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) UpdateContainer(
	updateDto dto.UpdateContainer,
) {
	recordCode, _ := valueObject.NewActivityRecordCode("ContainerUpdated")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []valueObject.SystemResourceIdentifier{
			valueObject.NewContainerSri(updateDto.AccountId, updateDto.ContainerId),
		},
		RecordDetails:     updateDto,
		OperatorAccountId: &updateDto.OperatorAccountId,
		OperatorIpAddress: &updateDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) DeleteContainer(
	deleteDto dto.DeleteContainer,
) {
	recordCode, _ := valueObject.NewActivityRecordCode("ContainerDeleted")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []valueObject.SystemResourceIdentifier{
			valueObject.NewContainerSri(deleteDto.AccountId, deleteDto.ContainerId),
		},
		OperatorAccountId: &deleteDto.OperatorAccountId,
		OperatorIpAddress: &deleteDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) CreateContainerProfile(
	createDto dto.CreateContainerProfile,
	profileId valueObject.ContainerProfileId,
) {
	recordCode, _ := valueObject.NewActivityRecordCode("ContainerProfileCreated")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []valueObject.SystemResourceIdentifier{
			valueObject.NewContainerProfileSri(createDto.AccountId, profileId),
		},
		OperatorAccountId: &createDto.OperatorAccountId,
		OperatorIpAddress: &createDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) UpdateContainerProfile(
	updateDto dto.UpdateContainerProfile,
) {
	recordCode, _ := valueObject.NewActivityRecordCode("ContainerProfileUpdated")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []valueObject.SystemResourceIdentifier{
			valueObject.NewContainerProfileSri(updateDto.AccountId, updateDto.ProfileId),
		},
		RecordDetails:     updateDto,
		OperatorAccountId: &updateDto.OperatorAccountId,
		OperatorIpAddress: &updateDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) DeleteContainerProfile(
	deleteDto dto.DeleteContainerProfile,
) {
	recordCode, _ := valueObject.NewActivityRecordCode("ContainerProfileDeleted")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []valueObject.SystemResourceIdentifier{
			valueObject.NewContainerProfileSri(deleteDto.AccountId, deleteDto.ProfileId),
		},
		OperatorAccountId: &deleteDto.OperatorAccountId,
		OperatorIpAddress: &deleteDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) CreateMapping(
	createDto dto.CreateMapping,
	mappingId valueObject.MappingId,
) {
	recordCode, _ := valueObject.NewActivityRecordCode("MappingCreated")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []valueObject.SystemResourceIdentifier{
			valueObject.NewMappingSri(createDto.AccountId, mappingId),
		},
		OperatorAccountId: &createDto.OperatorAccountId,
		OperatorIpAddress: &createDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) DeleteMapping(
	deleteDto dto.DeleteMapping,
) {
	recordCode, _ := valueObject.NewActivityRecordCode("MappingDeleted")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []valueObject.SystemResourceIdentifier{
			valueObject.NewMappingSri(deleteDto.AccountId, deleteDto.MappingId),
		},
		OperatorAccountId: &deleteDto.OperatorAccountId,
		OperatorIpAddress: &deleteDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) CreateMappingTarget(
	createDto dto.CreateMappingTarget,
	targetId valueObject.MappingTargetId,
) {
	recordCode, _ := valueObject.NewActivityRecordCode("MappingTargetCreated")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []valueObject.SystemResourceIdentifier{
			valueObject.NewMappingSri(createDto.AccountId, createDto.MappingId),
			valueObject.NewMappingTargetSri(createDto.AccountId, targetId),
		},
		OperatorAccountId: &createDto.OperatorAccountId,
		OperatorIpAddress: &createDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) DeleteMappingTarget(
	deleteDto dto.DeleteMappingTarget,
) {
	recordCode, _ := valueObject.NewActivityRecordCode("MappingTargetDeleted")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []valueObject.SystemResourceIdentifier{
			valueObject.NewMappingSri(deleteDto.AccountId, deleteDto.MappingId),
			valueObject.NewMappingTargetSri(deleteDto.AccountId, deleteDto.TargetId),
		},
		OperatorAccountId: &deleteDto.OperatorAccountId,
		OperatorIpAddress: &deleteDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) CreateContainerSnapshotImage(
	createDto dto.CreateContainerSnapshotImage,
	imageAccountId valueObject.AccountId,
	imageId valueObject.ContainerImageId,
) {
	recordCode, _ := valueObject.NewActivityRecordCode("ContainerSnapshotImageCreated")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []valueObject.SystemResourceIdentifier{
			valueObject.NewContainerSri(imageAccountId, createDto.ContainerId),
			valueObject.NewContainerImageSri(imageAccountId, imageId),
		},
		OperatorAccountId: &createDto.OperatorAccountId,
		OperatorIpAddress: &createDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) DeleteContainerImage(
	deleteDto dto.DeleteContainerImage,
) {
	recordCode, _ := valueObject.NewActivityRecordCode("ContainerImageDeleted")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []valueObject.SystemResourceIdentifier{
			valueObject.NewContainerImageSri(deleteDto.AccountId, deleteDto.ImageId),
		},
		OperatorAccountId: &deleteDto.OperatorAccountId,
		OperatorIpAddress: &deleteDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) CreateContainerImageArchiveFile(
	createDto dto.CreateContainerImageArchiveFile,
) {
	recordCode, _ := valueObject.NewActivityRecordCode("ContainerImageArchiveFileCreated")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []valueObject.SystemResourceIdentifier{
			valueObject.NewContainerImageSri(createDto.AccountId, createDto.ImageId),
		},
		OperatorAccountId: &createDto.OperatorAccountId,
		OperatorIpAddress: &createDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) ImportContainerImageArchiveFile(
	importDto dto.ImportContainerImageArchiveFile,
	imageId valueObject.ContainerImageId,
) {
	recordCode, _ := valueObject.NewActivityRecordCode("ContainerImageArchiveFileImported")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []valueObject.SystemResourceIdentifier{
			valueObject.NewContainerImageSri(importDto.AccountId, imageId),
		},
		OperatorAccountId: &importDto.OperatorAccountId,
		OperatorIpAddress: &importDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) DeleteContainerImageArchiveFile(
	deleteDto dto.DeleteContainerImageArchiveFile,
) {
	recordCode, _ := valueObject.NewActivityRecordCode("ContainerImageArchiveFileDeleted")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []valueObject.SystemResourceIdentifier{
			valueObject.NewContainerImageSri(deleteDto.AccountId, deleteDto.ImageId),
		},
		OperatorAccountId: &deleteDto.OperatorAccountId,
		OperatorIpAddress: &deleteDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) CreateBackupDestination(
	createDto dto.CreateBackupDestination,
	destinationId valueObject.BackupDestinationId,
) {
	recordCode, _ := valueObject.NewActivityRecordCode("BackupDestinationCreated")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []valueObject.SystemResourceIdentifier{
			valueObject.NewBackupDestinationSri(createDto.AccountId, destinationId),
		},
		OperatorAccountId: &createDto.OperatorAccountId,
		OperatorIpAddress: &createDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) UpdateBackupDestination(
	updateDto dto.UpdateBackupDestination,
) {
	recordCode, _ := valueObject.NewActivityRecordCode("BackupDestinationUpdated")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []valueObject.SystemResourceIdentifier{
			valueObject.NewBackupDestinationSri(updateDto.AccountId, updateDto.DestinationId),
		},
		RecordDetails:     updateDto,
		OperatorAccountId: &updateDto.OperatorAccountId,
		OperatorIpAddress: &updateDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) DeleteBackupDestination(
	deleteDto dto.DeleteBackupDestination,
) {
	recordCode, _ := valueObject.NewActivityRecordCode("BackupDestinationDeleted")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []valueObject.SystemResourceIdentifier{
			valueObject.NewBackupDestinationSri(deleteDto.AccountId, deleteDto.DestinationId),
		},
		OperatorAccountId: &deleteDto.OperatorAccountId,
		OperatorIpAddress: &deleteDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) CreateBackupJob(
	createDto dto.CreateBackupJob,
	jobId valueObject.BackupJobId,
) {
	recordCode, _ := valueObject.NewActivityRecordCode("BackupJobCreated")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []valueObject.SystemResourceIdentifier{
			valueObject.NewBackupJobSri(createDto.AccountId, jobId),
		},
		OperatorAccountId: &createDto.OperatorAccountId,
		OperatorIpAddress: &createDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) UpdateBackupJob(
	updateDto dto.UpdateBackupJob,
) {
	recordCode, _ := valueObject.NewActivityRecordCode("BackupJobUpdated")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []valueObject.SystemResourceIdentifier{
			valueObject.NewBackupJobSri(updateDto.AccountId, updateDto.JobId),
		},
		RecordDetails:     updateDto,
		OperatorAccountId: &updateDto.OperatorAccountId,
		OperatorIpAddress: &updateDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}
