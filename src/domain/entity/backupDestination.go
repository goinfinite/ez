package entity

import (
	"github.com/goinfinite/ez/src/domain/valueObject"
)

type IBackupDestination interface {
}

type BackupDestinationBase struct {
	DestinationId                     valueObject.BackupDestinationId           `json:"destinationId"`
	AccountId                         valueObject.AccountId                     `json:"accountId"`
	AccountUsername                   valueObject.UnixUsername                  `json:"accountUsername"`
	DestinationName                   valueObject.BackupDestinationName         `json:"destinationName"`
	DestinationDescription            *valueObject.BackupDestinationDescription `json:"destinationDescription"`
	DestinationType                   valueObject.BackupDestinationType         `json:"destinationType"`
	DestinationPath                   valueObject.UnixFilePath                  `json:"destinationPath"`
	MinLocalStorageFreePercent        *uint8                                    `json:"minLocalStorageFreePercent"`
	MaxDestinationStorageUsagePercent *uint8                                    `json:"maxDestinationStorageUsagePercent"`
	EncryptionKey                     valueObject.Password                      `json:"-"`
	TasksCount                        uint16                                    `json:"tasksCount"`
	TotalSpaceUsageBytes              valueObject.Byte                          `json:"totalSpaceUsageBytes"`
	CreatedAt                         valueObject.UnixTime                      `json:"createdAt"`
	UpdatedAt                         valueObject.UnixTime                      `json:"updatedAt"`
}

type BackupDestinationRemoteBase struct {
	BackupDestinationBase
	MaxConcurrentConnections    *uint16           `json:"maxConcurrentConnections"`
	DownloadBytesSecRateLimit   *valueObject.Byte `json:"downloadBytesSecRateLimit"`
	UploadBytesSecRateLimit     *valueObject.Byte `json:"uploadBytesSecRateLimit"`
	SkipCertificateVerification *bool             `json:"skipCertificateVerification"`
}

type BackupDestinationLocal struct {
	BackupDestinationBase
}

type BackupDestinationObjectStorage struct {
	BackupDestinationRemoteBase
	ObjectStorageProvider                *valueObject.ObjectStorageProvider                `json:"objectStorageProvider"`
	ObjectStorageProviderRegion          *valueObject.ObjectStorageProviderRegion          `json:"objectStorageProviderRegion"`
	ObjectStorageProviderAccessKeyId     *valueObject.ObjectStorageProviderAccessKeyId     `json:"objectStorageProviderAccessKeyId"`
	ObjectStorageProviderSecretAccessKey *valueObject.ObjectStorageProviderSecretAccessKey `json:"-"`
	ObjectStorageEndpointUrl             *valueObject.Url                                  `json:"objectStorageEndpointUrl"`
	ObjectStorageBucketName              *valueObject.ObjectStorageBucketName              `json:"objectStorageBucketName"`
}

type BackupDestinationRemoteHost struct {
	BackupDestinationRemoteBase
	RemoteHostType                  *valueObject.BackupDestinationRemoteHostType `json:"remoteHostType"`
	RemoteHostname                  *valueObject.NetworkHost                     `json:"remoteHostname"`
	RemoteHostNetworkPort           *valueObject.NetworkPort                     `json:"remoteHostNetworkPort"`
	RemoteHostUsername              *valueObject.UnixUsername                    `json:"remoteHostUsername"`
	RemoteHostPassword              *valueObject.Password                        `json:"-"`
	RemoteHostPrivateKeyFilePath    *valueObject.UnixFilePath                    `json:"remoteHostPrivateKeyFilePath"`
	RemoteHostConnectionTimeoutSecs *valueObject.TimeDuration                    `json:"remoteHostConnectionTimeoutSecs"`
	RemoteHostConnectionRetrySecs   *valueObject.TimeDuration                    `json:"remoteHostConnectionRetrySecs"`
}
