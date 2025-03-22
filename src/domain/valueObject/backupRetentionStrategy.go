package valueObject

import (
	"errors"
	"strings"

	voHelper "github.com/goinfinite/ez/src/domain/valueObject/helper"
)

type BackupRetentionStrategy string

const (
	BackupRetentionStrategyFull        BackupRetentionStrategy = "full"
	BackupRetentionStrategyIncremental BackupRetentionStrategy = "incremental"
)

func NewBackupRetentionStrategy(value interface{}) (
	retentionStrategy BackupRetentionStrategy, err error,
) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return retentionStrategy, errors.New("BackupRetentionStrategyMustBeString")
	}
	stringValue = strings.ToLower(stringValue)

	stringValueVo := BackupRetentionStrategy(stringValue)
	switch stringValueVo {
	case BackupRetentionStrategyFull, BackupRetentionStrategyIncremental:
		return stringValueVo, nil
	default:
		return retentionStrategy, errors.New("InvalidBackupRetentionStrategy")
	}
}

func (vo BackupRetentionStrategy) String() string {
	return string(vo)
}
