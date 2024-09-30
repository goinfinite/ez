package dto

import (
	"mime/multipart"

	"github.com/goinfinite/ez/src/domain/valueObject"
)

type ImportContainerImageArchiveFile struct {
	AccountId         valueObject.AccountId `json:"accountId"`
	ArchiveFile       *multipart.FileHeader `json:"archiveFile"`
	OperatorAccountId valueObject.AccountId `json:"-"`
	OperatorIpAddress valueObject.IpAddress `json:"-"`
}

func NewImportContainerImageArchiveFile(
	accountId valueObject.AccountId,
	archiveFile *multipart.FileHeader,
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
) ImportContainerImageArchiveFile {
	return ImportContainerImageArchiveFile{
		AccountId:         accountId,
		ArchiveFile:       archiveFile,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
	}
}
