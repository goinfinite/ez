package valueObject

import (
	"errors"
	"strings"

	voHelper "github.com/goinfinite/ez/src/domain/valueObject/helper"
)

type BackupDestinationType string

const (
	BackupDestinationTypeLocal         BackupDestinationType = "local"
	BackupDestinationTypeRemoteHost    BackupDestinationType = "remote-host"
	BackupDestinationTypeObjectStorage BackupDestinationType = "object-storage"
)

func NewBackupDestinationType(value interface{}) (
	destinationType BackupDestinationType, err error,
) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return destinationType, errors.New("BackupDestinationTypeMustBeString")
	}
	stringValue = strings.ToLower(stringValue)

	stringValueVo := BackupDestinationType(stringValue)
	switch stringValueVo {
	case BackupDestinationTypeLocal, BackupDestinationTypeRemoteHost,
		BackupDestinationTypeObjectStorage:
		return stringValueVo, nil
	default:
		return destinationType, errors.New("InvalidBackupDestinationType")
	}
}

func (vo BackupDestinationType) String() string {
	return string(vo)
}
