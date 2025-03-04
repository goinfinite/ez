package dto

import (
	"github.com/goinfinite/ez/src/domain/entity"
	"github.com/goinfinite/ez/src/domain/valueObject"
)

type ReadBackupDestinationsRequest struct {
	Pagination            Pagination                                   `json:"pagination"`
	DestinationId         *valueObject.BackupDestinationId             `json:"destinationId"`
	AccountId             *valueObject.AccountId                       `json:"accountId"`
	DestinationName       *valueObject.BackupDestinationName           `json:"destinationName"`
	DestinationType       *valueObject.BackupDestinationType           `json:"destinationType"`
	ObjectStorageProvider *valueObject.ObjectStorageProvider           `json:"objectStorageProvider"`
	RemoteHostType        *valueObject.BackupDestinationRemoteHostType `json:"remoteHostType"`
	RemoteHostname        *valueObject.Fqdn                            `json:"remoteHostname"`
	CreatedBeforeAt       *valueObject.UnixTime                        `json:"createdBeforeAt"`
	CreatedAfterAt        *valueObject.UnixTime                        `json:"createdAfterAt"`
}

type ReadBackupDestinationsResponse struct {
	Pagination   Pagination                  `json:"pagination"`
	Destinations []entity.IBackupDestination `json:"destinations"`
}
