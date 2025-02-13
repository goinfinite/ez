package dto

import (
	"mime/multipart"

	"github.com/goinfinite/ez/src/domain/valueObject"
)

type ImportContainerImageArchive struct {
	AccountId         valueObject.AccountId     `json:"accountId"`
	ArchiveFile       *multipart.FileHeader     `json:"archiveFile"`
	ArchiveFilePath   *valueObject.UnixFilePath `json:"-"`
	OperatorAccountId valueObject.AccountId     `json:"-"`
	OperatorIpAddress valueObject.IpAddress     `json:"-"`
}

func NewImportContainerImageArchive(
	accountId valueObject.AccountId,
	archiveFile *multipart.FileHeader,
	archiveFilePath *valueObject.UnixFilePath,
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
) ImportContainerImageArchive {
	return ImportContainerImageArchive{
		AccountId:         accountId,
		ArchiveFile:       archiveFile,
		ArchiveFilePath:   archiveFilePath,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
	}
}
