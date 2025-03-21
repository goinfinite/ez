package entity

import (
	"github.com/goinfinite/ez/src/domain/valueObject"
)

type IBackupDestination interface {
	ReadDestinationId() valueObject.BackupDestinationId
	ReadAccountId() valueObject.AccountId
	ReadAccountUsername() valueObject.UnixUsername
	ReadDestinationName() valueObject.BackupDestinationName
	ReadDestinationDescription() *valueObject.BackupDestinationDescription
	ReadDestinationType() valueObject.BackupDestinationType
	ReadDestinationPath() valueObject.UnixFilePath
	ReadMinLocalStorageFreePercent() *uint8
	ReadMaxDestinationStorageUsagePercent() *uint8
	ReadEncryptionKey() valueObject.Password
	ReadTasksCount() uint16
	ReadTotalSpaceUsageBytes() valueObject.Byte
	ReadCreatedAt() valueObject.UnixTime
	ReadUpdatedAt() valueObject.UnixTime
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

func (entity BackupDestinationBase) ReadDestinationId() valueObject.BackupDestinationId {
	return entity.DestinationId
}

func (entity BackupDestinationBase) ReadAccountId() valueObject.AccountId {
	return entity.AccountId
}

func (entity BackupDestinationBase) ReadAccountUsername() valueObject.UnixUsername {
	return entity.AccountUsername
}

func (entity BackupDestinationBase) ReadDestinationName() valueObject.BackupDestinationName {
	return entity.DestinationName
}

func (entity BackupDestinationBase) ReadDestinationDescription() *valueObject.BackupDestinationDescription {
	return entity.DestinationDescription
}

func (entity BackupDestinationBase) ReadDestinationType() valueObject.BackupDestinationType {
	return entity.DestinationType
}

func (entity BackupDestinationBase) ReadDestinationPath() valueObject.UnixFilePath {
	return entity.DestinationPath
}

func (entity BackupDestinationBase) ReadMinLocalStorageFreePercent() *uint8 {
	return entity.MinLocalStorageFreePercent
}

func (entity BackupDestinationBase) ReadMaxDestinationStorageUsagePercent() *uint8 {
	return entity.MaxDestinationStorageUsagePercent
}

func (entity BackupDestinationBase) ReadEncryptionKey() valueObject.Password {
	return entity.EncryptionKey
}

func (entity BackupDestinationBase) ReadTasksCount() uint16 {
	return entity.TasksCount
}

func (entity BackupDestinationBase) ReadTotalSpaceUsageBytes() valueObject.Byte {
	return entity.TotalSpaceUsageBytes
}

func (entity BackupDestinationBase) ReadCreatedAt() valueObject.UnixTime {
	return entity.CreatedAt
}

func (entity BackupDestinationBase) ReadUpdatedAt() valueObject.UnixTime {
	return entity.UpdatedAt
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
