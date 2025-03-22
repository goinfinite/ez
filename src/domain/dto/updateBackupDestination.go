package dto

import (
	"github.com/goinfinite/ez/src/domain/valueObject"
)

type UpdateBackupDestination struct {
	DestinationId                        valueObject.BackupDestinationId                   `json:"destinationId"`
	AccountId                            valueObject.AccountId                             `json:"accountId"`
	DestinationName                      *valueObject.BackupDestinationName                `json:"destinationName"`
	DestinationDescription               *valueObject.BackupDestinationDescription         `json:"destinationDescription"`
	DestinationPath                      *valueObject.UnixFilePath                         `json:"destinationPath,omitempty"`
	MinLocalStorageFreePercent           *uint8                                            `json:"minLocalStorageFreePercent,omitempty"`
	MaxDestinationStorageUsagePercent    *uint8                                            `json:"maxDestinationStorageUsagePercent,omitempty"`
	MaxConcurrentConnections             *uint16                                           `json:"maxConcurrentConnections,omitempty"`
	DownloadBytesSecRateLimit            *valueObject.Byte                                 `json:"downloadBytesSecRateLimit,omitempty"`
	UploadBytesSecRateLimit              *valueObject.Byte                                 `json:"uploadBytesSecRateLimit,omitempty"`
	SkipCertificateVerification          *bool                                             `json:"skipCertificateVerification,omitempty"`
	ObjectStorageProvider                *valueObject.ObjectStorageProvider                `json:"objectStorageProvider,omitempty"`
	ObjectStorageProviderRegion          *valueObject.ObjectStorageProviderRegion          `json:"objectStorageProviderRegion,omitempty"`
	ObjectStorageProviderAccessKeyId     *valueObject.ObjectStorageProviderAccessKeyId     `json:"objectStorageProviderAccessKeyId,omitempty"`
	ObjectStorageProviderSecretAccessKey *valueObject.ObjectStorageProviderSecretAccessKey `json:"objectStorageProviderSecretAccessKey,omitempty"`
	ObjectStorageEndpointUrl             *valueObject.Url                                  `json:"objectStorageEndpointUrl,omitempty"`
	ObjectStorageBucketName              *valueObject.ObjectStorageBucketName              `json:"objectStorageBucketName,omitempty"`
	RemoteHostType                       *valueObject.BackupDestinationRemoteHostType      `json:"remoteHostType,omitempty"`
	RemoteHostname                       *valueObject.NetworkHost                          `json:"remoteHostname,omitempty"`
	RemoteHostNetworkPort                *valueObject.NetworkPort                          `json:"remoteHostNetworkPort,omitempty"`
	RemoteHostUsername                   *valueObject.UnixUsername                         `json:"remoteHostUsername,omitempty"`
	RemoteHostPassword                   *valueObject.Password                             `json:"remoteHostPassword,omitempty"`
	RemoteHostPrivateKeyFilePath         *valueObject.UnixFilePath                         `json:"remoteHostPrivateKeyFilePath,omitempty"`
	RemoteHostConnectionTimeoutSecs      *valueObject.TimeDuration                         `json:"remoteHostConnectionTimeoutSecs,omitempty"`
	RemoteHostConnectionRetrySecs        *valueObject.TimeDuration                         `json:"remoteHostConnectionRetrySecs,omitempty"`
	TasksCount                           *uint16                                           `json:"-"`
	TotalSpaceUsageBytes                 *valueObject.Byte                                 `json:"-"`
	OperatorAccountId                    valueObject.AccountId                             `json:"-"`
	OperatorIpAddress                    valueObject.IpAddress                             `json:"-"`
}
