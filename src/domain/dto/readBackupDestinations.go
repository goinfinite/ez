package dto

import (
	"github.com/goinfinite/ez/src/domain/entity"
	"github.com/goinfinite/ez/src/domain/valueObject"
)

type ReadBackupDestinationsRequest struct {
	Pagination            Pagination                                   `json:"pagination"`
	DestinationId         *valueObject.BackupDestinationId             `json:"destinationId,omitempty"`
	AccountId             *valueObject.AccountId                       `json:"accountId,omitempty"`
	DestinationName       *valueObject.BackupDestinationName           `json:"destinationName,omitempty"`
	DestinationType       *valueObject.BackupDestinationType           `json:"destinationType,omitempty"`
	ObjectStorageProvider *valueObject.ObjectStorageProvider           `json:"objectStorageProvider,omitempty"`
	RemoteHostType        *valueObject.BackupDestinationRemoteHostType `json:"remoteHostType,omitempty"`
	RemoteHostname        *valueObject.Fqdn                            `json:"remoteHostname,omitempty"`
	CreatedBeforeAt       *valueObject.UnixTime                        `json:"createdBeforeAt,omitempty"`
	CreatedAfterAt        *valueObject.UnixTime                        `json:"createdAfterAt,omitempty"`
}

type ReadBackupDestinationsResponse struct {
	Pagination   Pagination                  `json:"pagination"`
	Destinations []entity.IBackupDestination `json:"destinations"`
}
