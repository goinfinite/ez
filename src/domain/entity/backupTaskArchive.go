package entity

import (
	"github.com/goinfinite/ez/src/domain/valueObject"
)

type BackupTaskArchive struct {
	ArchiveId       valueObject.BackupTaskArchiveId `json:"archiveId"`
	AccountId       valueObject.AccountId           `json:"accountId"`
	AccountUsername valueObject.UnixUsername        `json:"accountUsername"`
	TaskId          valueObject.BackupTaskId        `json:"taskId"`
	UnixFilePath    valueObject.UnixFilePath        `json:"unixFilePath"`
	SizeBytes       valueObject.Byte                `json:"sizeBytes"`
	DownloadUrl     *valueObject.Url                `json:"downloadUrl"`
	CreatedAt       valueObject.UnixTime            `json:"createdAt"`
}

func NewBackupTaskArchive(
	archiveId valueObject.BackupTaskArchiveId,
	accountId valueObject.AccountId,
	accountUsername valueObject.UnixUsername,
	taskId valueObject.BackupTaskId,
	unixFilePath valueObject.UnixFilePath,
	sizeBytes valueObject.Byte,
	downloadUrl *valueObject.Url,
	createdAt valueObject.UnixTime,
) BackupTaskArchive {
	return BackupTaskArchive{
		ArchiveId:       archiveId,
		AccountId:       accountId,
		AccountUsername: accountUsername,
		TaskId:          taskId,
		UnixFilePath:    unixFilePath,
		SizeBytes:       sizeBytes,
		DownloadUrl:     downloadUrl,
		CreatedAt:       createdAt,
	}
}
