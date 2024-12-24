package entity

import (
	"github.com/goinfinite/ez/src/domain/valueObject"
)

type IBackupDestination interface {
}

type BackupDestinationBase struct {
	DestinationId                     valueObject.BackupDestinationId           `json:"destinationId"`
	AccountId                         valueObject.AccountId                     `json:"accountId"`
	DestinationName                   valueObject.BackupDestinationName         `json:"destinationName"`
	DestinationDescription            *valueObject.BackupDestinationDescription `json:"destinationDescription"`
	DestinationType                   valueObject.BackupDestinationType         `json:"destinationType"`
	DestinationPath                   valueObject.UnixFilePath                  `json:"destinationPath"`
	MinLocalStorageFreePercent        *uint8                                    `json:"minLocalStorageFreePercent,omitempty"`
	MaxDestinationStorageUsagePercent *uint8                                    `json:"maxDestinationStorageUsagePercent,omitempty"`
	TotalSpaceUsageBytes              *valueObject.Byte                         `json:"totalSpaceUsageBytes"`
	TotalSpaceUsagePercent            *uint8                                    `json:"totalSpaceUsagePercent,omitempty"`
	CreatedAt                         valueObject.UnixTime                      `json:"createdAt"`
	UpdatedAt                         valueObject.UnixTime                      `json:"updatedAt"`
}

type BackupDestinationRemoteBase struct {
	BackupDestinationBase
	MaxConcurrentConnections    *uint16 `json:"maxConcurrentConnections,omitempty"`
	DownloadBytesSecRateLimit   *uint64 `json:"downloadBytesSecRateLimit,omitempty"`
	UploadBytesSecRateLimit     *uint64 `json:"uploadBytesSecRateLimit,omitempty"`
	SkipCertificateVerification *bool   `json:"skipCertificateVerification,omitempty"`
}

type BackupDestinationLocal struct {
	BackupDestinationBase
}

type BackupDestinationObjectStorage struct {
	BackupDestinationRemoteBase
	ObjectStorageProvider                *valueObject.ObjectStorageProvider                `json:"objectStorageProvider,omitempty"`
	ObjectStorageProviderRegion          *valueObject.ObjectStorageProviderRegion          `json:"objectStorageProviderRegion,omitempty"`
	ObjectStorageProviderAccessKeyId     *valueObject.ObjectStorageProviderAccessKeyId     `json:"objectStorageProviderAccessKeyId,omitempty"`
	ObjectStorageProviderSecretAccessKey *valueObject.ObjectStorageProviderSecretAccessKey `json:"-"`
	ObjectStorageEndpointUrl             *valueObject.Url                                  `json:"objectStorageEndpointUrl,omitempty"`
	ObjectStorageBucketName              *valueObject.ObjectStorageBucketName              `json:"objectStorageBucketName,omitempty"`
}

type BackupDestinationRemoteHost struct {
	BackupDestinationRemoteBase
	RemoteHostType                  *valueObject.BackupDestinationRemoteHostType `json:"remoteHostType,omitempty"`
	RemoteHostname                  *valueObject.NetworkHost                     `json:"remoteHostname,omitempty"`
	RemoteHostNetworkPort           *valueObject.NetworkPort                     `json:"remoteHostNetworkPort,omitempty"`
	RemoteHostUsername              *valueObject.UnixUsername                    `json:"remoteHostUsername,omitempty"`
	RemoteHostPassword              *valueObject.Password                        `json:"-"`
	RemoteHostPrivateKeyFilePath    *valueObject.UnixFilePath                    `json:"remoteHostPrivateKeyFilePath,omitempty"`
	RemoteHostConnectionTimeoutSecs *uint16                                      `json:"remoteHostConnectionTimeoutSecs,omitempty"`
	RemoteHostConnectionRetrySecs   *uint16                                      `json:"remoteHostConnectionRetrySecs,omitempty"`
}
