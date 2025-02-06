package valueObject

import (
	"errors"
)

type BackupTaskArchiveId string

func NewBackupTaskArchiveId(value interface{}) (BackupTaskArchiveId, error) {
	archiveId, err := NewHash(value)
	if err != nil {
		return "", errors.New("InvalidBackupTaskArchiveId")
	}

	return BackupTaskArchiveId(archiveId), nil
}

func (vo BackupTaskArchiveId) String() string {
	return string(vo)
}
