package service

import (
	"errors"
	"log/slog"
	"strconv"
	"strings"

	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/useCase"
	"github.com/goinfinite/ez/src/domain/valueObject"
	voHelper "github.com/goinfinite/ez/src/domain/valueObject/helper"
	"github.com/goinfinite/ez/src/infra"
	backupInfra "github.com/goinfinite/ez/src/infra/backup"
	"github.com/goinfinite/ez/src/infra/db"
	infraEnvs "github.com/goinfinite/ez/src/infra/envs"
	infraHelper "github.com/goinfinite/ez/src/infra/helper"
	serviceHelper "github.com/goinfinite/ez/src/presentation/service/helper"
)

type BackupService struct {
	persistentDbSvc       *db.PersistentDatabaseService
	backupQueryRepo       *backupInfra.BackupQueryRepo
	backupCmdRepo         *backupInfra.BackupCmdRepo
	activityRecordCmdRepo *infra.ActivityRecordCmdRepo
}

func NewBackupService(
	persistentDbSvc *db.PersistentDatabaseService,
	trailDbSvc *db.TrailDatabaseService,
) *BackupService {
	return &BackupService{
		persistentDbSvc:       persistentDbSvc,
		backupQueryRepo:       backupInfra.NewBackupQueryRepo(persistentDbSvc),
		backupCmdRepo:         backupInfra.NewBackupCmdRepo(persistentDbSvc, trailDbSvc),
		activityRecordCmdRepo: infra.NewActivityRecordCmdRepo(trailDbSvc),
	}
}

func (service *BackupService) ReadDestination(input map[string]interface{}) ServiceOutput {
	var destinationIdPtr *valueObject.BackupDestinationId
	if input["destinationId"] != nil {
		destinationId, err := valueObject.NewBackupDestinationId(input["destinationId"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		destinationIdPtr = &destinationId
	}

	var accountIdPtr *valueObject.AccountId
	if input["accountId"] != nil {
		accountId, err := valueObject.NewAccountId(input["accountId"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		accountIdPtr = &accountId
	}

	var destinationNamePtr *valueObject.BackupDestinationName
	if input["destinationName"] != nil {
		destinationName, err := valueObject.NewBackupDestinationName(input["destinationName"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		destinationNamePtr = &destinationName
	}

	var destinationTypePtr *valueObject.BackupDestinationType
	if input["destinationType"] != nil {
		destinationType, err := valueObject.NewBackupDestinationType(input["destinationType"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		destinationTypePtr = &destinationType
	}

	var objectStorageProviderPtr *valueObject.ObjectStorageProvider
	if input["objectStorageProvider"] != nil {
		objectStorageProvider, err := valueObject.NewObjectStorageProvider(input["objectStorageProvider"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		objectStorageProviderPtr = &objectStorageProvider
	}

	var remoteHostnamePtr *valueObject.Fqdn
	if input["remoteHostname"] != nil {
		remoteHostname, err := valueObject.NewFqdn(input["remoteHostname"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		remoteHostnamePtr = &remoteHostname
	}

	var remoteHostTypePtr *valueObject.BackupDestinationRemoteHostType
	if input["remoteHostType"] != nil {
		remoteHostType, err := valueObject.NewBackupDestinationRemoteHostType(input["remoteHostType"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		remoteHostTypePtr = &remoteHostType
	}

	timeParamNames := []string{"createdBeforeAt", "createdAfterAt"}
	timeParamPtrs := serviceHelper.TimeParamsParser(timeParamNames, input)

	requestPagination, err := serviceHelper.PaginationParser(
		input, useCase.BackupDestinationsDefaultPagination,
	)
	if err != nil {
		return NewServiceOutput(UserError, err)
	}

	requestDto := dto.ReadBackupDestinationsRequest{
		Pagination:            requestPagination,
		DestinationId:         destinationIdPtr,
		AccountId:             accountIdPtr,
		DestinationName:       destinationNamePtr,
		DestinationType:       destinationTypePtr,
		ObjectStorageProvider: objectStorageProviderPtr,
		RemoteHostType:        remoteHostTypePtr,
		RemoteHostname:        remoteHostnamePtr,
		CreatedBeforeAt:       timeParamPtrs["createdBeforeAt"],
		CreatedAfterAt:        timeParamPtrs["createdAfterAt"],
	}

	responseDto, err := useCase.ReadBackupDestinations(service.backupQueryRepo, requestDto)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, responseDto)
}

func (service *BackupService) CreateDestination(
	input map[string]interface{},
) ServiceOutput {
	requiredParams := []string{"accountId", "destinationName", "destinationType"}

	err := serviceHelper.RequiredParamsInspector(input, requiredParams)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	accountId, err := valueObject.NewAccountId(input["accountId"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	destinationName, err := valueObject.NewBackupDestinationName(input["destinationName"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	var destinationDescriptionPtr *valueObject.BackupDestinationDescription
	if input["destinationDescription"] != nil {
		destinationDescription, err := valueObject.NewBackupDestinationDescription(
			input["destinationDescription"],
		)
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		destinationDescriptionPtr = &destinationDescription
	}

	destinationType, err := valueObject.NewBackupDestinationType(input["destinationType"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	var destinationPathPtr *valueObject.UnixFilePath
	if input["destinationPath"] != nil {
		destinationPath, err := valueObject.NewUnixFilePath(input["destinationPath"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		destinationPathPtr = &destinationPath
	}

	var minLocalStorageFreePercentPtr *uint8
	if input["minLocalStorageFreePercent"] != nil {
		minLocalStorageFreePercent, err := voHelper.InterfaceToUint8(
			input["minLocalStorageFreePercent"],
		)
		if err != nil {
			return NewServiceOutput(UserError, errors.New("InvalidMinLocalStorageFreePercent"))
		}
		minLocalStorageFreePercentPtr = &minLocalStorageFreePercent
	}

	var maxDestinationStorageUsagePercentPtr *uint8
	if input["maxDestinationStorageUsagePercent"] != nil {
		maxDestinationStorageUsagePercent, err := voHelper.InterfaceToUint8(
			input["maxDestinationStorageUsagePercent"],
		)
		if err != nil {
			return NewServiceOutput(UserError, errors.New("InvalidMaxDestinationStorageUsagePercent"))
		}
		maxDestinationStorageUsagePercentPtr = &maxDestinationStorageUsagePercent
	}

	var maxConcurrentConnectionsPtr *uint16
	if input["maxConcurrentConnections"] != nil {
		maxConcurrentConnections, err := voHelper.InterfaceToUint16(
			input["maxConcurrentConnections"],
		)
		if err != nil {
			return NewServiceOutput(UserError, errors.New("InvalidMaxConcurrentConnections"))
		}
		maxConcurrentConnectionsPtr = &maxConcurrentConnections
	}

	var downloadBytesSecRateLimitPtr *valueObject.Byte
	if input["downloadBytesSecRateLimit"] != nil {
		downloadBytesSecRateLimit, err := valueObject.NewByte(
			input["downloadBytesSecRateLimit"],
		)
		if err != nil {
			return NewServiceOutput(UserError, errors.New("InvalidDownloadBytesSecRateLimit"))
		}
		downloadBytesSecRateLimitPtr = &downloadBytesSecRateLimit
	}

	var uploadBytesSecRateLimitPtr *valueObject.Byte
	if input["uploadBytesSecRateLimit"] != nil {
		uploadBytesSecRateLimit, err := valueObject.NewByte(
			input["uploadBytesSecRateLimit"],
		)
		if err != nil {
			return NewServiceOutput(UserError, errors.New("InvalidUploadBytesSecRateLimit"))
		}
		uploadBytesSecRateLimitPtr = &uploadBytesSecRateLimit
	}

	var skipCertificateVerificationPtr *bool
	if input["skipCertificateVerification"] != nil {
		skipCertificateVerification, err := voHelper.InterfaceToBool(
			input["skipCertificateVerification"],
		)
		if err != nil {
			return NewServiceOutput(UserError, errors.New("InvalidSkipCertificateVerification"))
		}
		skipCertificateVerificationPtr = &skipCertificateVerification
	}

	var objectStorageProviderPtr *valueObject.ObjectStorageProvider
	if input["objectStorageProvider"] != nil {
		objectStorageProvider, err := valueObject.NewObjectStorageProvider(
			input["objectStorageProvider"],
		)
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		objectStorageProviderPtr = &objectStorageProvider
	}

	var objectStorageProviderRegionPtr *valueObject.ObjectStorageProviderRegion
	if input["objectStorageProviderRegion"] != nil {
		objectStorageProviderRegion, err := valueObject.NewObjectStorageProviderRegion(
			input["objectStorageProviderRegion"],
		)
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		objectStorageProviderRegionPtr = &objectStorageProviderRegion
	}

	var objectStorageProviderAccessKeyIdPtr *valueObject.ObjectStorageProviderAccessKeyId
	if input["objectStorageProviderAccessKeyId"] != nil {
		objectStorageProviderAccessKeyId, err := valueObject.NewObjectStorageProviderAccessKeyId(
			input["objectStorageProviderAccessKeyId"],
		)
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		objectStorageProviderAccessKeyIdPtr = &objectStorageProviderAccessKeyId
	}

	var objectStorageProviderSecretAccessKeyPtr *valueObject.ObjectStorageProviderSecretAccessKey
	if input["objectStorageProviderSecretAccessKey"] != nil {
		objectStorageProviderSecretAccessKey, err := valueObject.NewObjectStorageProviderSecretAccessKey(
			input["objectStorageProviderSecretAccessKey"],
		)
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		objectStorageProviderSecretAccessKeyPtr = &objectStorageProviderSecretAccessKey
	}

	var objectStorageEndpointUrlPtr *valueObject.Url
	if input["objectStorageEndpointUrl"] != nil {
		objectStorageEndpointUrl, err := valueObject.NewUrl(input["objectStorageEndpointUrl"])
		if err != nil {
			return NewServiceOutput(UserError, errors.New("InvalidObjectStorageEndpointUrl"))
		}
		objectStorageEndpointUrlPtr = &objectStorageEndpointUrl
	}

	var objectStorageBucketNamePtr *valueObject.ObjectStorageBucketName
	if input["objectStorageBucketName"] != nil {
		objectStorageBucketName, err := valueObject.NewObjectStorageBucketName(
			input["objectStorageBucketName"],
		)
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		objectStorageBucketNamePtr = &objectStorageBucketName
	}

	var remoteHostTypePtr *valueObject.BackupDestinationRemoteHostType
	if input["remoteHostType"] != nil {
		remoteHostType, err := valueObject.NewBackupDestinationRemoteHostType(
			input["remoteHostType"],
		)
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		remoteHostTypePtr = &remoteHostType
	}

	var remoteHostnamePtr *valueObject.NetworkHost
	if input["remoteHostname"] != nil {
		remoteHostname, err := valueObject.NewNetworkHost(input["remoteHostname"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		remoteHostnamePtr = &remoteHostname
	}

	var remoteHostNetworkPortPtr *valueObject.NetworkPort
	if input["remoteHostNetworkPort"] != nil {
		remoteHostNetworkPort, err := valueObject.NewNetworkPort(
			input["remoteHostNetworkPort"],
		)
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		remoteHostNetworkPortPtr = &remoteHostNetworkPort
	}

	var remoteHostUsernamePtr *valueObject.UnixUsername
	if input["remoteHostUsername"] != nil {
		remoteHostUsername, err := valueObject.NewUnixUsername(input["remoteHostUsername"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		remoteHostUsernamePtr = &remoteHostUsername
	}

	var remoteHostPasswordPtr *valueObject.Password
	if input["remoteHostPassword"] != nil {
		remoteHostPassword, err := valueObject.NewPassword(input["remoteHostPassword"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		remoteHostPasswordPtr = &remoteHostPassword
	}

	var remoteHostPrivateKeyFilePathPtr *valueObject.UnixFilePath
	if input["remoteHostPrivateKeyFilePath"] != nil {
		remoteHostPrivateKeyFilePath, err := valueObject.NewUnixFilePath(
			input["remoteHostPrivateKeyFilePath"],
		)
		if err != nil {
			return NewServiceOutput(UserError, errors.New("InvalidRemoteHostPrivateKeyFilePath"))
		}
		remoteHostPrivateKeyFilePathPtr = &remoteHostPrivateKeyFilePath
	}

	var remoteHostConnectionTimeoutSecsPtr *valueObject.TimeDuration
	if input["remoteHostConnectionTimeoutSecs"] != nil {
		remoteHostConnectionTimeoutSecs, err := valueObject.NewTimeDuration(
			input["remoteHostConnectionTimeoutSecs"],
		)
		if err != nil {
			return NewServiceOutput(UserError, errors.New("InvalidRemoteHostConnectionTimeoutSecs"))
		}
		remoteHostConnectionTimeoutSecsPtr = &remoteHostConnectionTimeoutSecs
	}

	var remoteHostConnectionRetrySecsPtr *valueObject.TimeDuration
	if input["remoteHostConnectionRetrySecs"] != nil {
		remoteHostConnectionRetrySecs, err := valueObject.NewTimeDuration(
			input["remoteHostConnectionRetrySecs"],
		)
		if err != nil {
			return NewServiceOutput(UserError, errors.New("InvalidRemoteHostConnectionRetrySecs"))
		}
		remoteHostConnectionRetrySecsPtr = &remoteHostConnectionRetrySecs
	}

	operatorAccountId := LocalOperatorAccountId
	if input["operatorAccountId"] != nil {
		operatorAccountId, err = valueObject.NewAccountId(input["operatorAccountId"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	operatorIpAddress := LocalOperatorIpAddress
	if input["operatorIpAddress"] != nil {
		operatorIpAddress, err = valueObject.NewIpAddress(input["operatorIpAddress"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	createDto := dto.CreateBackupDestinationRequest{
		AccountId:                            accountId,
		DestinationName:                      destinationName,
		DestinationDescription:               destinationDescriptionPtr,
		DestinationType:                      destinationType,
		DestinationPath:                      destinationPathPtr,
		MinLocalStorageFreePercent:           minLocalStorageFreePercentPtr,
		MaxDestinationStorageUsagePercent:    maxDestinationStorageUsagePercentPtr,
		MaxConcurrentConnections:             maxConcurrentConnectionsPtr,
		DownloadBytesSecRateLimit:            downloadBytesSecRateLimitPtr,
		UploadBytesSecRateLimit:              uploadBytesSecRateLimitPtr,
		SkipCertificateVerification:          skipCertificateVerificationPtr,
		ObjectStorageProvider:                objectStorageProviderPtr,
		ObjectStorageProviderRegion:          objectStorageProviderRegionPtr,
		ObjectStorageProviderAccessKeyId:     objectStorageProviderAccessKeyIdPtr,
		ObjectStorageProviderSecretAccessKey: objectStorageProviderSecretAccessKeyPtr,
		ObjectStorageEndpointUrl:             objectStorageEndpointUrlPtr,
		ObjectStorageBucketName:              objectStorageBucketNamePtr,
		RemoteHostType:                       remoteHostTypePtr,
		RemoteHostname:                       remoteHostnamePtr,
		RemoteHostNetworkPort:                remoteHostNetworkPortPtr,
		RemoteHostUsername:                   remoteHostUsernamePtr,
		RemoteHostPassword:                   remoteHostPasswordPtr,
		RemoteHostPrivateKeyFilePath:         remoteHostPrivateKeyFilePathPtr,
		RemoteHostConnectionTimeoutSecs:      remoteHostConnectionTimeoutSecsPtr,
		RemoteHostConnectionRetrySecs:        remoteHostConnectionRetrySecsPtr,
		OperatorAccountId:                    operatorAccountId,
		OperatorIpAddress:                    operatorIpAddress,
	}

	responseDto, err := useCase.CreateBackupDestination(
		service.backupCmdRepo, service.activityRecordCmdRepo, createDto,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Created, responseDto)
}

func (service *BackupService) UpdateDestination(input map[string]interface{}) ServiceOutput {
	requiredParams := []string{"destinationId", "accountId"}

	err := serviceHelper.RequiredParamsInspector(input, requiredParams)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	destinationId, err := valueObject.NewBackupDestinationId(input["destinationId"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	accountId, err := valueObject.NewAccountId(input["accountId"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	var destinationNamePtr *valueObject.BackupDestinationName
	if input["destinationName"] != nil {
		destinationName, err := valueObject.NewBackupDestinationName(input["destinationName"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		destinationNamePtr = &destinationName
	}

	var destinationDescriptionPtr *valueObject.BackupDestinationDescription
	if input["destinationDescription"] != nil {
		destinationDescription, err := valueObject.NewBackupDestinationDescription(
			input["destinationDescription"],
		)
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		destinationDescriptionPtr = &destinationDescription
	}

	var destinationTypePtr *valueObject.BackupDestinationType
	if input["destinationType"] != nil {
		destinationType, err := valueObject.NewBackupDestinationType(input["destinationType"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		destinationTypePtr = &destinationType
	}

	var destinationPathPtr *valueObject.UnixFilePath
	if input["destinationPath"] != nil {
		destinationPath, err := valueObject.NewUnixFilePath(input["destinationPath"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		destinationPathPtr = &destinationPath
	}

	var minLocalStorageFreePercentPtr *uint8
	if input["minLocalStorageFreePercent"] != nil {
		minLocalStorageFreePercent, err := voHelper.InterfaceToUint8(
			input["minLocalStorageFreePercent"],
		)
		if err != nil {
			return NewServiceOutput(UserError, errors.New("InvalidMinLocalStorageFreePercent"))
		}
		minLocalStorageFreePercentPtr = &minLocalStorageFreePercent
	}

	var maxDestinationStorageUsagePercentPtr *uint8
	if input["maxDestinationStorageUsagePercent"] != nil {
		maxDestinationStorageUsagePercent, err := voHelper.InterfaceToUint8(
			input["maxDestinationStorageUsagePercent"],
		)
		if err != nil {
			return NewServiceOutput(UserError, errors.New("InvalidMaxDestinationStorageUsagePercent"))
		}
		maxDestinationStorageUsagePercentPtr = &maxDestinationStorageUsagePercent
	}

	var maxConcurrentConnectionsPtr *uint16
	if input["maxConcurrentConnections"] != nil {
		maxConcurrentConnections, err := voHelper.InterfaceToUint16(
			input["maxConcurrentConnections"],
		)
		if err != nil {
			return NewServiceOutput(UserError, errors.New("InvalidMaxConcurrentConnections"))
		}
		maxConcurrentConnectionsPtr = &maxConcurrentConnections
	}

	var downloadBytesSecRateLimitPtr *valueObject.Byte
	if input["downloadBytesSecRateLimit"] != nil {
		downloadBytesSecRateLimit, err := valueObject.NewByte(
			input["downloadBytesSecRateLimit"],
		)
		if err != nil {
			return NewServiceOutput(UserError, errors.New("InvalidDownloadBytesSecRateLimit"))
		}
		downloadBytesSecRateLimitPtr = &downloadBytesSecRateLimit
	}

	var uploadBytesSecRateLimitPtr *valueObject.Byte
	if input["uploadBytesSecRateLimit"] != nil {
		uploadBytesSecRateLimit, err := valueObject.NewByte(
			input["uploadBytesSecRateLimit"],
		)
		if err != nil {
			return NewServiceOutput(UserError, errors.New("InvalidUploadBytesSecRateLimit"))
		}
		uploadBytesSecRateLimitPtr = &uploadBytesSecRateLimit
	}

	var skipCertificateVerificationPtr *bool
	if input["skipCertificateVerification"] != nil {
		skipCertificateVerification, err := voHelper.InterfaceToBool(
			input["skipCertificateVerification"],
		)
		if err != nil {
			return NewServiceOutput(UserError, errors.New("InvalidSkipCertificateVerification"))
		}
		skipCertificateVerificationPtr = &skipCertificateVerification
	}

	var objectStorageProviderPtr *valueObject.ObjectStorageProvider
	if input["objectStorageProvider"] != nil {
		objectStorageProvider, err := valueObject.NewObjectStorageProvider(
			input["objectStorageProvider"],
		)
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		objectStorageProviderPtr = &objectStorageProvider
	}

	var objectStorageProviderRegionPtr *valueObject.ObjectStorageProviderRegion
	if input["objectStorageProviderRegion"] != nil {
		objectStorageProviderRegion, err := valueObject.NewObjectStorageProviderRegion(
			input["objectStorageProviderRegion"],
		)
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		objectStorageProviderRegionPtr = &objectStorageProviderRegion
	}

	var objectStorageProviderAccessKeyIdPtr *valueObject.ObjectStorageProviderAccessKeyId
	if input["objectStorageProviderAccessKeyId"] != nil {
		objectStorageProviderAccessKeyId, err := valueObject.NewObjectStorageProviderAccessKeyId(
			input["objectStorageProviderAccessKeyId"],
		)
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		objectStorageProviderAccessKeyIdPtr = &objectStorageProviderAccessKeyId
	}

	var objectStorageProviderSecretAccessKeyPtr *valueObject.ObjectStorageProviderSecretAccessKey
	if input["objectStorageProviderSecretAccessKey"] != nil {
		objectStorageProviderSecretAccessKey, err := valueObject.NewObjectStorageProviderSecretAccessKey(
			input["objectStorageProviderSecretAccessKey"],
		)
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		objectStorageProviderSecretAccessKeyPtr = &objectStorageProviderSecretAccessKey
	}

	var objectStorageEndpointUrlPtr *valueObject.Url
	if input["objectStorageEndpointUrl"] != nil {
		objectStorageEndpointUrl, err := valueObject.NewUrl(input["objectStorageEndpointUrl"])
		if err != nil {
			return NewServiceOutput(UserError, errors.New("InvalidObjectStorageEndpointUrl"))
		}
		objectStorageEndpointUrlPtr = &objectStorageEndpointUrl
	}

	var objectStorageBucketNamePtr *valueObject.ObjectStorageBucketName
	if input["objectStorageBucketName"] != nil {
		objectStorageBucketName, err := valueObject.NewObjectStorageBucketName(
			input["objectStorageBucketName"],
		)
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		objectStorageBucketNamePtr = &objectStorageBucketName
	}

	var remoteHostTypePtr *valueObject.BackupDestinationRemoteHostType
	if input["remoteHostType"] != nil {
		remoteHostType, err := valueObject.NewBackupDestinationRemoteHostType(
			input["remoteHostType"],
		)
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		remoteHostTypePtr = &remoteHostType
	}

	var remoteHostnamePtr *valueObject.NetworkHost
	if input["remoteHostname"] != nil {
		remoteHostname, err := valueObject.NewNetworkHost(input["remoteHostname"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		remoteHostnamePtr = &remoteHostname
	}

	var remoteHostNetworkPortPtr *valueObject.NetworkPort
	if input["remoteHostNetworkPort"] != nil {
		remoteHostNetworkPort, err := valueObject.NewNetworkPort(
			input["remoteHostNetworkPort"],
		)
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		remoteHostNetworkPortPtr = &remoteHostNetworkPort
	}

	var remoteHostUsernamePtr *valueObject.UnixUsername
	if input["remoteHostUsername"] != nil {
		remoteHostUsername, err := valueObject.NewUnixUsername(input["remoteHostUsername"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		remoteHostUsernamePtr = &remoteHostUsername
	}

	var remoteHostPasswordPtr *valueObject.Password
	if input["remoteHostPassword"] != nil {
		remoteHostPassword, err := valueObject.NewPassword(input["remoteHostPassword"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		remoteHostPasswordPtr = &remoteHostPassword
	}

	var remoteHostPrivateKeyFilePathPtr *valueObject.UnixFilePath
	if input["remoteHostPrivateKeyFilePath"] != nil {
		remoteHostPrivateKeyFilePath, err := valueObject.NewUnixFilePath(
			input["remoteHostPrivateKeyFilePath"],
		)
		if err != nil {
			return NewServiceOutput(UserError, errors.New("InvalidRemoteHostPrivateKeyFilePath"))
		}
		remoteHostPrivateKeyFilePathPtr = &remoteHostPrivateKeyFilePath
	}

	var remoteHostConnectionTimeoutSecsPtr *valueObject.TimeDuration
	if input["remoteHostConnectionTimeoutSecs"] != nil {
		remoteHostConnectionTimeoutSecs, err := valueObject.NewTimeDuration(
			input["remoteHostConnectionTimeoutSecs"],
		)
		if err != nil {
			return NewServiceOutput(UserError, errors.New("InvalidRemoteHostConnectionTimeoutSecs"))
		}
		remoteHostConnectionTimeoutSecsPtr = &remoteHostConnectionTimeoutSecs
	}

	var remoteHostConnectionRetrySecsPtr *valueObject.TimeDuration
	if input["remoteHostConnectionRetrySecs"] != nil {
		remoteHostConnectionRetrySecs, err := valueObject.NewTimeDuration(
			input["remoteHostConnectionRetrySecs"],
		)
		if err != nil {
			return NewServiceOutput(UserError, errors.New("InvalidRemoteHostConnectionRetrySecs"))
		}
		remoteHostConnectionRetrySecsPtr = &remoteHostConnectionRetrySecs
	}

	operatorAccountId := LocalOperatorAccountId
	if input["operatorAccountId"] != nil {
		operatorAccountId, err = valueObject.NewAccountId(input["operatorAccountId"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	operatorIpAddress := LocalOperatorIpAddress
	if input["operatorIpAddress"] != nil {
		operatorIpAddress, err = valueObject.NewIpAddress(input["operatorIpAddress"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	updateDto := dto.UpdateBackupDestination{
		DestinationId:                        destinationId,
		AccountId:                            accountId,
		DestinationName:                      destinationNamePtr,
		DestinationDescription:               destinationDescriptionPtr,
		DestinationType:                      destinationTypePtr,
		DestinationPath:                      destinationPathPtr,
		MinLocalStorageFreePercent:           minLocalStorageFreePercentPtr,
		MaxDestinationStorageUsagePercent:    maxDestinationStorageUsagePercentPtr,
		MaxConcurrentConnections:             maxConcurrentConnectionsPtr,
		DownloadBytesSecRateLimit:            downloadBytesSecRateLimitPtr,
		UploadBytesSecRateLimit:              uploadBytesSecRateLimitPtr,
		SkipCertificateVerification:          skipCertificateVerificationPtr,
		ObjectStorageProvider:                objectStorageProviderPtr,
		ObjectStorageProviderRegion:          objectStorageProviderRegionPtr,
		ObjectStorageProviderAccessKeyId:     objectStorageProviderAccessKeyIdPtr,
		ObjectStorageProviderSecretAccessKey: objectStorageProviderSecretAccessKeyPtr,
		ObjectStorageEndpointUrl:             objectStorageEndpointUrlPtr,
		ObjectStorageBucketName:              objectStorageBucketNamePtr,
		RemoteHostType:                       remoteHostTypePtr,
		RemoteHostname:                       remoteHostnamePtr,
		RemoteHostNetworkPort:                remoteHostNetworkPortPtr,
		RemoteHostUsername:                   remoteHostUsernamePtr,
		RemoteHostPassword:                   remoteHostPasswordPtr,
		RemoteHostPrivateKeyFilePath:         remoteHostPrivateKeyFilePathPtr,
		RemoteHostConnectionTimeoutSecs:      remoteHostConnectionTimeoutSecsPtr,
		RemoteHostConnectionRetrySecs:        remoteHostConnectionRetrySecsPtr,
		OperatorAccountId:                    operatorAccountId,
		OperatorIpAddress:                    operatorIpAddress,
	}

	err = useCase.UpdateBackupDestination(
		service.backupQueryRepo, service.backupCmdRepo,
		service.activityRecordCmdRepo, updateDto,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, "BackupDestinationUpdated")
}

func (service *BackupService) DeleteDestination(input map[string]interface{}) ServiceOutput {
	requiredParams := []string{"destinationId", "accountId"}

	err := serviceHelper.RequiredParamsInspector(input, requiredParams)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	destinationId, err := valueObject.NewBackupDestinationId(input["destinationId"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	accountId, err := valueObject.NewAccountId(input["accountId"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	operatorAccountId := LocalOperatorAccountId
	if input["operatorAccountId"] != nil {
		operatorAccountId, err = valueObject.NewAccountId(input["operatorAccountId"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	operatorIpAddress := LocalOperatorIpAddress
	if input["operatorIpAddress"] != nil {
		operatorIpAddress, err = valueObject.NewIpAddress(input["operatorIpAddress"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	deleteDto := dto.NewDeleteBackupDestination(
		destinationId, accountId, operatorAccountId, operatorIpAddress,
	)

	err = useCase.DeleteBackupDestination(
		service.backupQueryRepo, service.backupCmdRepo,
		service.activityRecordCmdRepo, deleteDto,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, "BackupDestinationDeleted")
}

func (service *BackupService) ReadJob(input map[string]interface{}) ServiceOutput {
	var jobIdPtr *valueObject.BackupJobId
	if input["jobId"] != nil {
		jobId, err := valueObject.NewBackupJobId(input["jobId"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		jobIdPtr = &jobId
	}

	var jobStatusPtr *bool
	if input["jobStatus"] != nil {
		jobStatus, err := voHelper.InterfaceToBool(input["jobStatus"])
		if err != nil {
			return NewServiceOutput(UserError, errors.New("InvalidJobStatus"))
		}
		jobStatusPtr = &jobStatus
	}

	var accountIdPtr *valueObject.AccountId
	if input["accountId"] != nil {
		accountId, err := valueObject.NewAccountId(input["accountId"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		accountIdPtr = &accountId
	}

	var destinationIdPtr *valueObject.BackupDestinationId
	if input["destinationId"] != nil {
		destinationId, err := valueObject.NewBackupDestinationId(input["destinationId"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		destinationIdPtr = &destinationId
	}

	var retentionStrategyPtr *valueObject.BackupRetentionStrategy
	if input["retentionStrategy"] != nil {
		retentionStrategy, err := valueObject.NewBackupRetentionStrategy(input["retentionStrategy"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		retentionStrategyPtr = &retentionStrategy
	}

	var archiveCompressionFormatPtr *valueObject.CompressionFormat
	if input["archiveCompressionFormat"] != nil {
		archiveCompressionFormat, err := valueObject.NewCompressionFormat(input["archiveCompressionFormat"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		archiveCompressionFormatPtr = &archiveCompressionFormat
	}

	var lastRunStatusPtr *valueObject.BackupTaskStatus
	if input["lastRunStatus"] != nil {
		lastRunStatus, err := valueObject.NewBackupTaskStatus(input["lastRunStatus"])
		if err != nil {
			return NewServiceOutput(UserError, errors.New("InvalidLastRunStatus"))
		}
		lastRunStatusPtr = &lastRunStatus
	}

	timeParamNames := []string{
		"lastRunBeforeAt", "lastRunAfterAt", "nextRunBeforeAt", "nextRunAfterAt",
		"createdBeforeAt", "createdAfterAt",
	}
	timeParamPtrs := serviceHelper.TimeParamsParser(timeParamNames, input)

	requestPagination, err := serviceHelper.PaginationParser(
		input, useCase.BackupJobsDefaultPagination,
	)
	if err != nil {
		return NewServiceOutput(UserError, err)
	}

	requestDto := dto.ReadBackupJobsRequest{
		Pagination:               requestPagination,
		JobId:                    jobIdPtr,
		JobStatus:                jobStatusPtr,
		AccountId:                accountIdPtr,
		DestinationId:            destinationIdPtr,
		RetentionStrategy:        retentionStrategyPtr,
		ArchiveCompressionFormat: archiveCompressionFormatPtr,
		LastRunStatus:            lastRunStatusPtr,
		LastRunBeforeAt:          timeParamPtrs["lastRunBeforeAt"],
		LastRunAfterAt:           timeParamPtrs["lastRunAfterAt"],
		NextRunBeforeAt:          timeParamPtrs["nextRunBeforeAt"],
		NextRunAfterAt:           timeParamPtrs["nextRunAfterAt"],
		CreatedBeforeAt:          timeParamPtrs["createdBeforeAt"],
		CreatedAfterAt:           timeParamPtrs["createdAfterAt"],
	}

	responseDto, err := useCase.ReadBackupJobs(service.backupQueryRepo, requestDto)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, responseDto)
}

func (service *BackupService) CreateJob(input map[string]interface{}) ServiceOutput {
	requiredParams := []string{"accountId", "destinationIds", "backupSchedule"}

	err := serviceHelper.RequiredParamsInspector(input, requiredParams)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	accountId, err := valueObject.NewAccountId(input["accountId"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	var jobDescriptionPtr *valueObject.BackupJobDescription
	if input["jobDescription"] != nil {
		jobDescription, err := valueObject.NewBackupJobDescription(input["jobDescription"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		jobDescriptionPtr = &jobDescription
	}

	destinationIds, assertOk := input["destinationIds"].([]valueObject.BackupDestinationId)
	if !assertOk {
		return NewServiceOutput(UserError, errors.New("InvalidDestinationIds"))
	}

	backupSchedule, err := valueObject.NewCronSchedule(input["backupSchedule"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	var retentionStrategyPtr *valueObject.BackupRetentionStrategy
	if input["retentionStrategy"] != nil {
		retentionStrategy, err := valueObject.NewBackupRetentionStrategy(
			input["retentionStrategy"],
		)
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		retentionStrategyPtr = &retentionStrategy
	}

	var archiveCompressionFormatPtr *valueObject.CompressionFormat
	if input["archiveCompressionFormat"] != nil {
		archiveCompressionFormat, err := valueObject.NewCompressionFormat(
			input["archiveCompressionFormat"],
		)
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		archiveCompressionFormatPtr = &archiveCompressionFormat
	}

	var timeoutSecsPtr *valueObject.TimeDuration
	if input["timeoutSecs"] != nil {
		timeoutSecs, err := valueObject.NewTimeDuration(input["timeoutSecs"])
		if err != nil {
			return NewServiceOutput(UserError, errors.New("InvalidTimeoutSecs"))
		}
		timeoutSecsPtr = &timeoutSecs
	}

	var maxTaskRetentionCountPtr *uint16
	if input["maxTaskRetentionCount"] != nil {
		maxTaskRetentionCount, err := voHelper.InterfaceToUint16(
			input["maxTaskRetentionCount"],
		)
		if err != nil {
			return NewServiceOutput(UserError, errors.New("InvalidMaxTaskRetentionCount"))
		}
		maxTaskRetentionCountPtr = &maxTaskRetentionCount
	}

	var maxTaskRetentionDaysPtr *uint16
	if input["maxTaskRetentionDays"] != nil {
		maxTaskRetentionDays, err := voHelper.InterfaceToUint16(input["maxTaskRetentionDays"])
		if err != nil {
			return NewServiceOutput(UserError, errors.New("InvalidMaxTaskRetentionDays"))
		}
		maxTaskRetentionDaysPtr = &maxTaskRetentionDays
	}

	var maxConcurrentCpuCoresPtr *uint16
	if input["maxConcurrentCpuCores"] != nil {
		maxConcurrentCpuCores, err := voHelper.InterfaceToUint16(input["maxConcurrentCpuCores"])
		if err != nil {
			return NewServiceOutput(UserError, errors.New("InvalidMaxConcurrentCpuCores"))
		}
		maxConcurrentCpuCoresPtr = &maxConcurrentCpuCores
	}

	containerAccountIds := []valueObject.AccountId{}
	if input["containerAccountIds"] != nil {
		containerAccountIds, assertOk = input["containerAccountIds"].([]valueObject.AccountId)
		if !assertOk {
			return NewServiceOutput(UserError, errors.New("InvalidContainerAccountIds"))
		}
	}

	containerIds := []valueObject.ContainerId{}
	if input["containerIds"] != nil {
		containerIds, assertOk = input["containerIds"].([]valueObject.ContainerId)
		if !assertOk {
			return NewServiceOutput(UserError, errors.New("InvalidContainerIds"))
		}
	}

	exceptContainerAccountIds := []valueObject.AccountId{}
	if input["exceptContainerAccountIds"] != nil {
		exceptContainerAccountIds, assertOk = input["exceptContainerAccountIds"].([]valueObject.AccountId)
		if !assertOk {
			return NewServiceOutput(UserError, errors.New("InvalidExceptContainerAccountIds"))
		}
	}

	exceptContainerIds := []valueObject.ContainerId{}
	if input["exceptContainerIds"] != nil {
		exceptContainerIds, assertOk = input["exceptContainerIds"].([]valueObject.ContainerId)
		if !assertOk {
			return NewServiceOutput(UserError, errors.New("InvalidExceptContainerIds"))
		}
	}

	operatorAccountId := LocalOperatorAccountId
	if input["operatorAccountId"] != nil {
		operatorAccountId, err = valueObject.NewAccountId(input["operatorAccountId"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	operatorIpAddress := LocalOperatorIpAddress
	if input["operatorIpAddress"] != nil {
		operatorIpAddress, err = valueObject.NewIpAddress(input["operatorIpAddress"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	createDto := dto.NewCreateBackupJob(
		accountId, jobDescriptionPtr, destinationIds, retentionStrategyPtr, backupSchedule,
		archiveCompressionFormatPtr, timeoutSecsPtr, maxTaskRetentionCountPtr,
		maxTaskRetentionDaysPtr, maxConcurrentCpuCoresPtr, containerAccountIds,
		containerIds, exceptContainerAccountIds, exceptContainerIds,
		operatorAccountId, operatorIpAddress,
	)

	err = useCase.CreateBackupJob(
		service.backupCmdRepo, service.activityRecordCmdRepo, createDto,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Created, "BackupJobCreated")
}

func (service *BackupService) UpdateJob(input map[string]interface{}) ServiceOutput {
	requiredParams := []string{"jobId", "accountId"}

	err := serviceHelper.RequiredParamsInspector(input, requiredParams)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	jobId, err := valueObject.NewBackupJobId(input["jobId"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	accountId, err := valueObject.NewAccountId(input["accountId"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	var jobStatusPtr *bool
	if input["jobStatus"] != nil {
		jobStatus, err := voHelper.InterfaceToBool(input["jobStatus"])
		if err != nil {
			return NewServiceOutput(UserError, errors.New("InvalidJobStatus"))
		}
		jobStatusPtr = &jobStatus
	}

	var jobDescriptionPtr *valueObject.BackupJobDescription
	if input["jobDescription"] != nil {
		jobDescription, err := valueObject.NewBackupJobDescription(input["jobDescription"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		jobDescriptionPtr = &jobDescription
	}

	var destinationIds []valueObject.BackupDestinationId
	var assertOk bool
	if input["destinationIds"] != nil {
		destinationIds, assertOk = input["destinationIds"].([]valueObject.BackupDestinationId)
		if !assertOk {
			return NewServiceOutput(UserError, errors.New("InvalidDestinationIds"))
		}
	}

	var backupSchedulePtr *valueObject.CronSchedule
	if input["backupSchedule"] != nil {
		backupSchedule, err := valueObject.NewCronSchedule(input["backupSchedule"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		backupSchedulePtr = &backupSchedule
	}

	var timeoutSecsPtr *valueObject.TimeDuration
	if input["timeoutSecs"] != nil {
		timeoutSecs, err := valueObject.NewTimeDuration(input["timeoutSecs"])
		if err != nil {
			return NewServiceOutput(UserError, errors.New("InvalidTimeoutSecs"))
		}
		timeoutSecsPtr = &timeoutSecs
	}

	var maxTaskRetentionCountPtr *uint16
	if input["maxTaskRetentionCount"] != nil {
		maxTaskRetentionCount, err := voHelper.InterfaceToUint16(
			input["maxTaskRetentionCount"],
		)
		if err != nil {
			return NewServiceOutput(UserError, errors.New("InvalidMaxTaskRetentionCount"))
		}
		maxTaskRetentionCountPtr = &maxTaskRetentionCount
	}

	var maxTaskRetentionDaysPtr *uint16
	if input["maxTaskRetentionDays"] != nil {
		maxTaskRetentionDays, err := voHelper.InterfaceToUint16(input["maxTaskRetentionDays"])
		if err != nil {
			return NewServiceOutput(UserError, errors.New("InvalidMaxTaskRetentionDays"))
		}
		maxTaskRetentionDaysPtr = &maxTaskRetentionDays
	}

	var maxConcurrentCpuCoresPtr *uint16
	if input["maxConcurrentCpuCores"] != nil {
		maxConcurrentCpuCores, err := voHelper.InterfaceToUint16(input["maxConcurrentCpuCores"])
		if err != nil {
			return NewServiceOutput(UserError, errors.New("InvalidMaxConcurrentCpuCores"))
		}
		maxConcurrentCpuCoresPtr = &maxConcurrentCpuCores
	}

	var containerAccountIds []valueObject.AccountId
	if input["containerAccountIds"] != nil {
		containerAccountIds, assertOk = input["containerAccountIds"].([]valueObject.AccountId)
		if !assertOk {
			return NewServiceOutput(UserError, errors.New("InvalidContainerAccountIds"))
		}
	}

	var containerIds []valueObject.ContainerId
	if input["containerIds"] != nil {
		containerIds, assertOk = input["containerIds"].([]valueObject.ContainerId)
		if !assertOk {
			return NewServiceOutput(UserError, errors.New("InvalidContainerIds"))
		}
	}

	var exceptContainerAccountIds []valueObject.AccountId
	if input["exceptContainerAccountIds"] != nil {
		exceptContainerAccountIds, assertOk = input["exceptContainerAccountIds"].([]valueObject.AccountId)
		if !assertOk {
			return NewServiceOutput(UserError, errors.New("InvalidExceptContainerAccountIds"))
		}
	}

	var exceptContainerIds []valueObject.ContainerId
	if input["exceptContainerIds"] != nil {
		exceptContainerIds, assertOk = input["exceptContainerIds"].([]valueObject.ContainerId)
		if !assertOk {
			return NewServiceOutput(UserError, errors.New("InvalidExceptContainerIds"))
		}
	}

	operatorAccountId := LocalOperatorAccountId
	if input["operatorAccountId"] != nil {
		operatorAccountId, err = valueObject.NewAccountId(input["operatorAccountId"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	operatorIpAddress := LocalOperatorIpAddress
	if input["operatorIpAddress"] != nil {
		operatorIpAddress, err = valueObject.NewIpAddress(input["operatorIpAddress"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	updateDto := dto.UpdateBackupJob{
		JobId:                     jobId,
		AccountId:                 accountId,
		JobStatus:                 jobStatusPtr,
		JobDescription:            jobDescriptionPtr,
		DestinationIds:            destinationIds,
		BackupSchedule:            backupSchedulePtr,
		TimeoutSecs:               timeoutSecsPtr,
		MaxTaskRetentionCount:     maxTaskRetentionCountPtr,
		MaxTaskRetentionDays:      maxTaskRetentionDaysPtr,
		MaxConcurrentCpuCores:     maxConcurrentCpuCoresPtr,
		ContainerAccountIds:       containerAccountIds,
		ContainerIds:              containerIds,
		ExceptContainerAccountIds: exceptContainerAccountIds,
		ExceptContainerIds:        exceptContainerIds,
		OperatorAccountId:         operatorAccountId,
		OperatorIpAddress:         operatorIpAddress,
	}

	err = useCase.UpdateBackupJob(
		service.backupQueryRepo, service.backupCmdRepo,
		service.activityRecordCmdRepo, updateDto,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, "BackupJobUpdated")
}

func (service *BackupService) DeleteJob(input map[string]interface{}) ServiceOutput {
	requiredParams := []string{"jobId", "accountId"}

	err := serviceHelper.RequiredParamsInspector(input, requiredParams)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	jobId, err := valueObject.NewBackupJobId(input["jobId"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	accountId, err := valueObject.NewAccountId(input["accountId"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	operatorAccountId := LocalOperatorAccountId
	if input["operatorAccountId"] != nil {
		operatorAccountId, err = valueObject.NewAccountId(input["operatorAccountId"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	operatorIpAddress := LocalOperatorIpAddress
	if input["operatorIpAddress"] != nil {
		operatorIpAddress, err = valueObject.NewIpAddress(input["operatorIpAddress"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	deleteDto := dto.NewDeleteBackupJob(jobId, accountId, operatorAccountId, operatorIpAddress)

	err = useCase.DeleteBackupJob(
		service.backupQueryRepo, service.backupCmdRepo,
		service.activityRecordCmdRepo, deleteDto,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, "BackupJobDeleted")
}

func (service *BackupService) RunJob(input map[string]interface{}) ServiceOutput {
	requiredParams := []string{"jobId", "accountId"}

	err := serviceHelper.RequiredParamsInspector(input, requiredParams)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	jobId, err := valueObject.NewBackupJobId(input["jobId"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	accountId, err := valueObject.NewAccountId(input["accountId"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	operatorAccountId := LocalOperatorAccountId
	if input["operatorAccountId"] != nil {
		operatorAccountId, err = valueObject.NewAccountId(input["operatorAccountId"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	operatorIpAddress := LocalOperatorIpAddress
	if input["operatorIpAddress"] != nil {
		operatorIpAddress, err = valueObject.NewIpAddress(input["operatorIpAddress"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	runDto := dto.NewRunBackupJob(jobId, accountId, operatorAccountId, operatorIpAddress)

	err = useCase.RunBackupJob(
		service.backupQueryRepo, service.backupCmdRepo,
		service.activityRecordCmdRepo, runDto,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Created, "BackupTaskCreated")
}

func (service *BackupService) ReadTask(input map[string]interface{}) ServiceOutput {
	var taskIdPtr *valueObject.BackupTaskId
	if input["taskId"] != nil {
		taskId, err := valueObject.NewBackupTaskId(input["taskId"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		taskIdPtr = &taskId
	}

	var accountIdPtr *valueObject.AccountId
	if input["accountId"] != nil {
		accountId, err := valueObject.NewAccountId(input["accountId"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		accountIdPtr = &accountId
	}

	var jobIdPtr *valueObject.BackupJobId
	if input["jobId"] != nil {
		jobId, err := valueObject.NewBackupJobId(input["jobId"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		jobIdPtr = &jobId
	}

	var destinationIdPtr *valueObject.BackupDestinationId
	if input["destinationId"] != nil {
		destinationId, err := valueObject.NewBackupDestinationId(input["destinationId"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		destinationIdPtr = &destinationId
	}

	var taskStatusPtr *valueObject.BackupTaskStatus
	if input["taskStatus"] != nil {
		taskStatus, err := valueObject.NewBackupTaskStatus(input["taskStatus"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		taskStatusPtr = &taskStatus
	}

	var retentionStrategyPtr *valueObject.BackupRetentionStrategy
	if input["retentionStrategy"] != nil {
		retentionStrategy, err := valueObject.NewBackupRetentionStrategy(input["retentionStrategy"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		retentionStrategyPtr = &retentionStrategy
	}

	var containerIdPtr *valueObject.ContainerId
	if input["containerId"] != nil {
		containerId, err := valueObject.NewContainerId(input["containerId"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		containerIdPtr = &containerId
	}

	timeParamNames := []string{
		"startedBeforeAt", "startedAfterAt", "finishedBeforeAt", "finishedAfterAt",
		"createdBeforeAt", "createdAfterAt",
	}
	timeParamPtrs := serviceHelper.TimeParamsParser(timeParamNames, input)

	requestPagination, err := serviceHelper.PaginationParser(
		input, useCase.BackupTasksDefaultPagination,
	)
	if err != nil {
		return NewServiceOutput(UserError, err)
	}

	requestDto := dto.ReadBackupTasksRequest{
		Pagination:        requestPagination,
		TaskId:            taskIdPtr,
		AccountId:         accountIdPtr,
		JobId:             jobIdPtr,
		DestinationId:     destinationIdPtr,
		TaskStatus:        taskStatusPtr,
		RetentionStrategy: retentionStrategyPtr,
		ContainerId:       containerIdPtr,
		StartedBeforeAt:   timeParamPtrs["startedBeforeAt"],
		StartedAfterAt:    timeParamPtrs["startedAfterAt"],
		FinishedBeforeAt:  timeParamPtrs["finishedBeforeAt"],
		FinishedAfterAt:   timeParamPtrs["finishedAfterAt"],
		CreatedBeforeAt:   timeParamPtrs["createdBeforeAt"],
		CreatedAfterAt:    timeParamPtrs["createdAfterAt"],
	}

	responseDto, err := useCase.ReadBackupTasks(service.backupQueryRepo, requestDto)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, responseDto)
}

func (service *BackupService) RestoreTask(
	input map[string]interface{},
	shouldSchedule bool,
) ServiceOutput {
	var taskIdPtr *valueObject.BackupTaskId
	if input["taskId"] != nil {
		taskId, err := valueObject.NewBackupTaskId(input["taskId"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		taskIdPtr = &taskId
	}

	var archiveIdPtr *valueObject.BackupTaskArchiveId
	if input["archiveId"] != nil {
		archiveId, err := valueObject.NewBackupTaskArchiveId(input["archiveId"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		archiveIdPtr = &archiveId
	}

	var err error
	shouldReplaceExistingContainers := false
	if input["shouldReplaceExistingContainers"] != nil {
		shouldReplaceExistingContainers, err = voHelper.InterfaceToBool(
			input["shouldReplaceExistingContainers"],
		)
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	shouldRestoreMappings := true
	if input["shouldRestoreMappings"] != nil {
		shouldRestoreMappings, err = voHelper.InterfaceToBool(input["shouldRestoreMappings"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	timeoutSeconds := useCase.RestoreBackupTaskDefaultTimeoutSecs
	if input["timeoutSeconds"] != nil {
		timeoutSeconds, err = voHelper.InterfaceToUint32(input["timeoutSeconds"])
		if err != nil {
			return NewServiceOutput(UserError, errors.New("InvalidTimeoutSeconds"))
		}
	}

	var assertOk bool
	var containerAccountIds []valueObject.AccountId
	if input["containerAccountIds"] != nil {
		containerAccountIds, assertOk = input["containerAccountIds"].([]valueObject.AccountId)
		if !assertOk {
			return NewServiceOutput(UserError, errors.New("InvalidAccountIds"))
		}
	}

	var containerIds []valueObject.ContainerId
	if input["containerIds"] != nil {
		containerIds, assertOk = input["containerIds"].([]valueObject.ContainerId)
		if !assertOk {
			return NewServiceOutput(UserError, errors.New("InvalidContainerIds"))
		}
	}

	var exceptContainerAccountIds []valueObject.AccountId
	if input["exceptContainerAccountIds"] != nil {
		exceptContainerAccountIds, assertOk = input["exceptContainerAccountIds"].([]valueObject.AccountId)
		if !assertOk {
			return NewServiceOutput(UserError, errors.New("InvalidExceptContainerAccountIds"))
		}
	}

	var exceptContainerIds []valueObject.ContainerId
	if input["exceptContainerIds"] != nil {
		exceptContainerIds, assertOk = input["exceptContainerIds"].([]valueObject.ContainerId)
		if !assertOk {
			return NewServiceOutput(UserError, errors.New("InvalidExceptContainerIds"))
		}
	}

	operatorAccountId := LocalOperatorAccountId
	if input["operatorAccountId"] != nil {
		operatorAccountId, err = valueObject.NewAccountId(input["operatorAccountId"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	operatorIpAddress := LocalOperatorIpAddress
	if input["operatorIpAddress"] != nil {
		operatorIpAddress, err = valueObject.NewIpAddress(input["operatorIpAddress"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	if shouldSchedule {
		cliCmd := infraEnvs.InfiniteEzBinary + " backup task restore"
		cliParams := []string{
			"--should-replace-existing-containers", strconv.FormatBool(shouldReplaceExistingContainers),
			"--should-restore-mappings", strconv.FormatBool(shouldRestoreMappings),
			"--timeout-secs", strconv.Itoa(int(timeoutSeconds)),
		}
		if taskIdPtr != nil {
			cliParams = append(cliParams, "--task-id", taskIdPtr.String())
		}
		if archiveIdPtr != nil {
			cliParams = append(cliParams, "--archive-id", archiveIdPtr.String())
		}

		for _, accountId := range containerAccountIds {
			cliParams = append(cliParams, "--container-account-ids", accountId.String())
		}
		for _, containerId := range containerIds {
			cliParams = append(cliParams, "--container-ids", containerId.String())
		}
		for _, accountId := range exceptContainerAccountIds {
			cliParams = append(cliParams, "--except-container-account-ids", accountId.String())
		}
		for _, containerId := range exceptContainerIds {
			cliParams = append(cliParams, "--except-container-ids", containerId.String())
		}

		if operatorAccountId != LocalOperatorAccountId {
			cliParams = append(cliParams, "--operator-account-id", operatorAccountId.String())
		}

		cliCmd = cliCmd + " " + strings.Join(cliParams, " ")

		scheduledTaskCmdRepo := infra.NewScheduledTaskCmdRepo(service.persistentDbSvc)
		taskName, _ := valueObject.NewScheduledTaskName("RestoreBackupTask")
		taskCmd, _ := valueObject.NewUnixCommand(cliCmd)
		taskTag, _ := valueObject.NewScheduledTaskTag("backup")
		taskTags := []valueObject.ScheduledTaskTag{taskTag}

		scheduledTaskCreateDto := dto.NewCreateScheduledTask(
			taskName, taskCmd, taskTags, &timeoutSeconds, nil,
		)

		err = useCase.CreateScheduledTask(scheduledTaskCmdRepo, scheduledTaskCreateDto)
		if err != nil {
			return NewServiceOutput(InfraError, err.Error())
		}

		return NewServiceOutput(Created, "BackupTaskRestoreScheduled")
	}

	requestRestoreDto := dto.NewRestoreBackupTaskRequest(
		taskIdPtr, archiveIdPtr, &shouldReplaceExistingContainers, &shouldRestoreMappings,
		&timeoutSeconds, containerAccountIds, containerIds, exceptContainerAccountIds,
		exceptContainerIds, operatorAccountId, operatorIpAddress,
	)

	responseRestoreDto, err := useCase.RestoreBackupTask(
		service.backupQueryRepo, service.backupCmdRepo,
		service.activityRecordCmdRepo, requestRestoreDto,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	responseStatusEnum := Success
	if len(responseRestoreDto.FailedContainerImageIds) > 0 {
		responseStatusEnum = MultiStatus
	}
	if len(responseRestoreDto.SuccessfulContainerIds) == 0 {
		responseStatusEnum = InfraError
	}

	return NewServiceOutput(responseStatusEnum, responseRestoreDto)
}

func (service *BackupService) DeleteTask(input map[string]interface{}) ServiceOutput {
	requiredParams := []string{"taskId"}
	err := serviceHelper.RequiredParamsInspector(input, requiredParams)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	taskId, err := valueObject.NewBackupTaskId(input["taskId"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	shouldDiscardFiles := false
	if input["shouldDiscardFiles"] != nil {
		shouldDiscardFiles, err = voHelper.InterfaceToBool(input["shouldDiscardFiles"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	operatorAccountId := LocalOperatorAccountId
	if input["operatorAccountId"] != nil {
		operatorAccountId, err = valueObject.NewAccountId(input["operatorAccountId"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	operatorIpAddress := LocalOperatorIpAddress
	if input["operatorIpAddress"] != nil {
		operatorIpAddress, err = valueObject.NewIpAddress(input["operatorIpAddress"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	deleteDto := dto.NewDeleteBackupTask(
		taskId, shouldDiscardFiles, operatorAccountId, operatorIpAddress,
	)

	err = useCase.DeleteBackupTask(
		service.backupQueryRepo, service.backupCmdRepo,
		service.activityRecordCmdRepo, deleteDto,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, "BackupTaskDeleted")
}

func (service *BackupService) ReadTaskArchive(
	input map[string]interface{},
	requestHostname *string,
) ServiceOutput {
	var archiveIdPtr *valueObject.BackupTaskArchiveId
	if input["archiveId"] != nil {
		archiveId, err := valueObject.NewBackupTaskArchiveId(input["archiveId"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		archiveIdPtr = &archiveId
	}

	var accountIdPtr *valueObject.AccountId
	if input["accountId"] != nil {
		accountId, err := valueObject.NewAccountId(input["accountId"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		accountIdPtr = &accountId
	}

	var taskIdPtr *valueObject.BackupTaskId
	if input["taskId"] != nil {
		taskId, err := valueObject.NewBackupTaskId(input["taskId"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		taskIdPtr = &taskId
	}

	timeParamNames := []string{"createdBeforeAt", "createdAfterAt"}
	timeParamPtrs := serviceHelper.TimeParamsParser(timeParamNames, input)

	requestPagination, err := serviceHelper.PaginationParser(
		input, useCase.BackupTaskArchivesDefaultPagination,
	)
	if err != nil {
		return NewServiceOutput(UserError, err)
	}

	requestDto := dto.ReadBackupTaskArchivesRequest{
		Pagination:      requestPagination,
		ArchiveId:       archiveIdPtr,
		AccountId:       accountIdPtr,
		TaskId:          taskIdPtr,
		CreatedBeforeAt: timeParamPtrs["createdBeforeAt"],
		CreatedAfterAt:  timeParamPtrs["createdAfterAt"],
	}

	responseDto, err := useCase.ReadBackupTaskArchives(service.backupQueryRepo, requestDto)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	if requestHostname != nil {
		serverHostname, err := infraHelper.ReadServerHostname()
		if err != nil {
			return NewServiceOutput(InfraError, err.Error())
		}
		serverHostnameStr := serverHostname.String()

		for archiveFileIndex, archiveFile := range responseDto.Archives {
			rawUpdatedUrl := strings.Replace(
				archiveFile.DownloadUrl.String(), serverHostnameStr, *requestHostname, 1,
			)

			updatedUrl, err := valueObject.NewUrl(rawUpdatedUrl)
			if err != nil {
				slog.Debug(
					"UpdateDownloadUrlError",
					slog.Int("archiveFileIndex", archiveFileIndex),
					slog.String("rawUpdatedUrl", rawUpdatedUrl),
				)
				continue
			}

			responseDto.Archives[archiveFileIndex].DownloadUrl = &updatedUrl
		}
	}

	return NewServiceOutput(Success, responseDto)
}

func (service *BackupService) CreateTaskArchive(
	input map[string]interface{},
	shouldSchedule bool,
) ServiceOutput {
	requiredParams := []string{"taskId"}
	err := serviceHelper.RequiredParamsInspector(input, requiredParams)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	taskId, err := valueObject.NewBackupTaskId(input["taskId"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	timeoutSeconds := useCase.CreateBackupTaskArchiveDefaultTimeoutSecs
	if input["timeoutSeconds"] != nil {
		timeoutSeconds, err = voHelper.InterfaceToUint32(input["timeoutSeconds"])
		if err != nil {
			return NewServiceOutput(UserError, errors.New("InvalidTimeoutSeconds"))
		}
	}

	var assertOk bool
	var containerAccountIds []valueObject.AccountId
	if input["containerAccountIds"] != nil {
		containerAccountIds, assertOk = input["containerAccountIds"].([]valueObject.AccountId)
		if !assertOk {
			return NewServiceOutput(UserError, errors.New("InvalidAccountIds"))
		}
	}

	var containerIds []valueObject.ContainerId
	if input["containerIds"] != nil {
		containerIds, assertOk = input["containerIds"].([]valueObject.ContainerId)
		if !assertOk {
			return NewServiceOutput(UserError, errors.New("InvalidContainerIds"))
		}
	}

	var exceptContainerAccountIds []valueObject.AccountId
	if input["exceptContainerAccountIds"] != nil {
		exceptContainerAccountIds, assertOk = input["exceptContainerAccountIds"].([]valueObject.AccountId)
		if !assertOk {
			return NewServiceOutput(UserError, errors.New("InvalidExceptContainerAccountIds"))
		}
	}

	var exceptContainerIds []valueObject.ContainerId
	if input["exceptContainerIds"] != nil {
		exceptContainerIds, assertOk = input["exceptContainerIds"].([]valueObject.ContainerId)
		if !assertOk {
			return NewServiceOutput(UserError, errors.New("InvalidExceptContainerIds"))
		}
	}

	operatorAccountId := LocalOperatorAccountId
	if input["operatorAccountId"] != nil {
		operatorAccountId, err = valueObject.NewAccountId(input["operatorAccountId"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	operatorIpAddress := LocalOperatorIpAddress
	if input["operatorIpAddress"] != nil {
		operatorIpAddress, err = valueObject.NewIpAddress(input["operatorIpAddress"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	if shouldSchedule {
		cliCmd := infraEnvs.InfiniteEzBinary + " backup task archive create"
		cliParams := []string{
			"--task-id", taskId.String(),
			"--operator-account-id", operatorAccountId.String(),
			"--timeout-secs", strconv.Itoa(int(timeoutSeconds)),
		}
		for _, accountId := range containerAccountIds {
			cliParams = append(cliParams, "--container-account-ids", accountId.String())
		}
		for _, containerId := range containerIds {
			cliParams = append(cliParams, "--container-ids", containerId.String())
		}
		for _, accountId := range exceptContainerAccountIds {
			cliParams = append(cliParams, "--except-container-account-ids", accountId.String())
		}
		for _, containerId := range exceptContainerIds {
			cliParams = append(cliParams, "--except-container-ids", containerId.String())
		}

		cliCmd = cliCmd + " " + strings.Join(cliParams, " ")

		scheduledTaskCmdRepo := infra.NewScheduledTaskCmdRepo(service.persistentDbSvc)
		taskName, _ := valueObject.NewScheduledTaskName("CreateBackupTaskArchive")
		taskCmd, _ := valueObject.NewUnixCommand(cliCmd)
		taskTag, _ := valueObject.NewScheduledTaskTag("backup")
		taskTags := []valueObject.ScheduledTaskTag{taskTag}

		scheduledTaskCreateDto := dto.NewCreateScheduledTask(
			taskName, taskCmd, taskTags, &timeoutSeconds, nil,
		)

		err = useCase.CreateScheduledTask(scheduledTaskCmdRepo, scheduledTaskCreateDto)
		if err != nil {
			return NewServiceOutput(InfraError, err.Error())
		}

		return NewServiceOutput(Created, "CreateBackupTaskArchiveScheduled")
	}

	createDto := dto.NewCreateBackupTaskArchive(
		taskId, &timeoutSeconds, containerAccountIds, containerIds, exceptContainerAccountIds,
		exceptContainerIds, operatorAccountId, operatorIpAddress,
	)

	archiveEntity, err := useCase.CreateBackupTaskArchive(
		service.backupQueryRepo, service.backupCmdRepo,
		service.activityRecordCmdRepo, createDto,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Created, archiveEntity)
}

func (service *BackupService) DeleteTaskArchive(input map[string]interface{}) ServiceOutput {
	requiredParams := []string{"archiveId"}

	err := serviceHelper.RequiredParamsInspector(input, requiredParams)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	archiveId, err := valueObject.NewBackupTaskArchiveId(input["archiveId"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	operatorAccountId := LocalOperatorAccountId
	if input["operatorAccountId"] != nil {
		operatorAccountId, err = valueObject.NewAccountId(input["operatorAccountId"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	operatorIpAddress := LocalOperatorIpAddress
	if input["operatorIpAddress"] != nil {
		operatorIpAddress, err = valueObject.NewIpAddress(input["operatorIpAddress"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	deleteDto := dto.NewDeleteBackupTaskArchive(
		archiveId, operatorAccountId, operatorIpAddress,
	)

	err = useCase.DeleteBackupTaskArchive(
		service.backupQueryRepo, service.backupCmdRepo,
		service.activityRecordCmdRepo, deleteDto,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, "BackupTaskArchiveDeleted")
}
