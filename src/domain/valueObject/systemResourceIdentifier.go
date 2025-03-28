package valueObject

import (
	"errors"
	"log/slog"
	"regexp"

	voHelper "github.com/goinfinite/ez/src/domain/valueObject/helper"
)

const systemResourceIdentifierRegex string = `^sri://(?P<accountId>[\d]{1,64}):(?P<resourceType>[a-zA-Z0-9][\w\-]{1,128})\/(?P<resourceId>[a-zA-Z0-9][\w\.\-]{0,256}|\*)$`

type SystemResourceIdentifier string

func NewSystemResourceIdentifier(
	value interface{},
) (sri SystemResourceIdentifier, err error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return sri, errors.New("SystemResourceIdentifierMustBeString")
	}

	re := regexp.MustCompile(systemResourceIdentifierRegex)
	if !re.MatchString(stringValue) {
		return sri, errors.New("InvalidSystemResourceIdentifier")
	}

	return SystemResourceIdentifier(stringValue), nil
}

// Note: this panic is solely to warn about the misuse of the VO, specifically for developer
// auditing, and has nothing to do with user input. This is not a standard and should not be
// followed for the development of other VOs.
func NewSystemResourceIdentifierIgnoreError(value interface{}) SystemResourceIdentifier {
	sri, err := NewSystemResourceIdentifier(value)
	if err != nil {
		panicMessage := "UnexpectedSystemResourceIdentifierCreationError"
		slog.Debug(panicMessage, slog.Any("value", value), slog.Any("error", err))
		panic(panicMessage)
	}

	return sri
}

func NewAccountSri(accountId AccountId) SystemResourceIdentifier {
	return NewSystemResourceIdentifierIgnoreError(
		"sri://0:account/" + accountId.String(),
	)
}

func NewContainerSri(accountId AccountId, containerId ContainerId) SystemResourceIdentifier {
	return NewSystemResourceIdentifierIgnoreError(
		"sri://" + accountId.String() + ":container/" + containerId.String(),
	)
}

func NewContainerProfileSri(
	accountId AccountId,
	containerProfileId ContainerProfileId,
) SystemResourceIdentifier {
	return NewSystemResourceIdentifierIgnoreError(
		"sri://" + accountId.String() + ":containerProfile/" + containerProfileId.String(),
	)
}

func NewContainerImageSri(
	accountId AccountId,
	containerImageId ContainerImageId,
) SystemResourceIdentifier {
	return NewSystemResourceIdentifierIgnoreError(
		"sri://" + accountId.String() + ":containerImage/" + containerImageId.String(),
	)
}

func NewMappingSri(accountId AccountId, mappingId MappingId) SystemResourceIdentifier {
	return NewSystemResourceIdentifierIgnoreError(
		"sri://" + accountId.String() + ":mapping/" + mappingId.String(),
	)
}

func NewMappingTargetSri(
	accountId AccountId,
	mappingTargetId MappingTargetId,
) SystemResourceIdentifier {
	return NewSystemResourceIdentifierIgnoreError(
		"sri://" + accountId.String() + ":mappingTarget/" + mappingTargetId.String(),
	)
}

func NewBackupDestinationSri(
	accountId AccountId,
	backupDestinationId BackupDestinationId,
) SystemResourceIdentifier {
	return NewSystemResourceIdentifierIgnoreError(
		"sri://" + accountId.String() + ":backupDestination/" + backupDestinationId.String(),
	)
}

func NewBackupJobSri(accountId AccountId, backupJobId BackupJobId) SystemResourceIdentifier {
	return NewSystemResourceIdentifierIgnoreError(
		"sri://" + accountId.String() + ":backupJob/" + backupJobId.String(),
	)
}

func NewBackupTaskSri(accountId AccountId, backupTaskId BackupTaskId) SystemResourceIdentifier {
	return NewSystemResourceIdentifierIgnoreError(
		"sri://" + accountId.String() + ":backupTask/" + backupTaskId.String(),
	)
}

func NewBackupTaskArchiveSri(
	accountId AccountId,
	backupTaskArchiveId BackupTaskArchiveId,
) SystemResourceIdentifier {
	return NewSystemResourceIdentifierIgnoreError(
		"sri://" + accountId.String() + ":backupTaskArchive/" + backupTaskArchiveId.String(),
	)
}

func (vo SystemResourceIdentifier) String() string {
	return string(vo)
}
