package valueObject

import (
	"errors"
	"strings"

	voHelper "github.com/goinfinite/ez/src/domain/valueObject/helper"
)

type BackupDestinationRemoteHostType string

const (
	BackupDestinationRemoteHostTypeFtp  BackupDestinationRemoteHostType = "ftp"
	BackupDestinationRemoteHostTypeSftp BackupDestinationRemoteHostType = "sftp"
)

func NewBackupDestinationRemoteHostType(value interface{}) (
	hostType BackupDestinationRemoteHostType, err error,
) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return hostType, errors.New("BackupDestinationRemoteHostTypeMustBeString")
	}
	stringValue = strings.ToLower(stringValue)

	stringValueVo := BackupDestinationRemoteHostType(stringValue)
	switch stringValueVo {
	case BackupDestinationRemoteHostTypeFtp, BackupDestinationRemoteHostTypeSftp:
		return stringValueVo, nil
	default:
		return hostType, errors.New("InvalidBackupDestinationRemoteHostType")
	}
}

func (vo BackupDestinationRemoteHostType) String() string {
	return string(vo)
}
