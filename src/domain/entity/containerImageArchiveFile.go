package entity

import "github.com/speedianet/control/src/domain/valueObject"

type ContainerImageArchiveFile struct {
	ImageId      valueObject.ContainerImageId `json:"imageId"`
	AccountId    valueObject.AccountId        `json:"accountId"`
	UnixFilePath valueObject.UnixFilePath     `json:"unixFilePath"`
	DownloadUrl  valueObject.Url              `json:"downloadUrl"`
	SizeBytes    valueObject.Byte             `json:"sizeBytes"`
	CreatedAt    valueObject.UnixTime         `json:"createdAt"`
}

func NewContainerImageArchiveFile(
	imageId valueObject.ContainerImageId,
	accountId valueObject.AccountId,
	unixFilePath valueObject.UnixFilePath,
	downloadUrl valueObject.Url,
	sizeBytes valueObject.Byte,
	createdAt valueObject.UnixTime,
) ContainerImageArchiveFile {
	return ContainerImageArchiveFile{
		ImageId:      imageId,
		AccountId:    accountId,
		UnixFilePath: unixFilePath,
		DownloadUrl:  downloadUrl,
		SizeBytes:    sizeBytes,
		CreatedAt:    createdAt,
	}
}
