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

func (service *BackupService) ReadDestinationRequestFactory(
	serviceInput map[string]interface{},
) (readRequestDto dto.ReadBackupDestinationsRequest, err error) {
	var destinationIdPtr *valueObject.BackupDestinationId
	if serviceInput["destinationId"] != nil {
		destinationId, err := valueObject.NewBackupDestinationId(serviceInput["destinationId"])
		if err != nil {
			return readRequestDto, err
		}
		destinationIdPtr = &destinationId
	}

	var accountIdPtr *valueObject.AccountId
	if serviceInput["accountId"] != nil {
		accountId, err := valueObject.NewAccountId(serviceInput["accountId"])
		if err != nil {
			return readRequestDto, err
		}
		accountIdPtr = &accountId
	}

	var destinationNamePtr *valueObject.BackupDestinationName
	if serviceInput["destinationName"] != nil {
		destinationName, err := valueObject.NewBackupDestinationName(serviceInput["destinationName"])
		if err != nil {
			return readRequestDto, err
		}
		destinationNamePtr = &destinationName
	}

	var destinationTypePtr *valueObject.BackupDestinationType
	if serviceInput["destinationType"] != nil {
		destinationType, err := valueObject.NewBackupDestinationType(serviceInput["destinationType"])
		if err != nil {
			return readRequestDto, err
		}
		destinationTypePtr = &destinationType
	}

	var objectStorageProviderPtr *valueObject.ObjectStorageProvider
	if serviceInput["objectStorageProvider"] != nil {
		objectStorageProvider, err := valueObject.NewObjectStorageProvider(serviceInput["objectStorageProvider"])
		if err != nil {
			return readRequestDto, err
		}
		objectStorageProviderPtr = &objectStorageProvider
	}

	var remoteHostnamePtr *valueObject.Fqdn
	if serviceInput["remoteHostname"] != nil {
		remoteHostname, err := valueObject.NewFqdn(serviceInput["remoteHostname"])
		if err != nil {
			return readRequestDto, err
		}
		remoteHostnamePtr = &remoteHostname
	}

	var remoteHostTypePtr *valueObject.BackupDestinationRemoteHostType
	if serviceInput["remoteHostType"] != nil {
		remoteHostType, err := valueObject.NewBackupDestinationRemoteHostType(serviceInput["remoteHostType"])
		if err != nil {
			return readRequestDto, err
		}
		remoteHostTypePtr = &remoteHostType
	}

	timeParamNames := []string{"createdBeforeAt", "createdAfterAt"}
	timeParamPtrs := serviceHelper.TimeParamsParser(timeParamNames, serviceInput)

	requestPagination, err := serviceHelper.PaginationParser(
		serviceInput, useCase.BackupDestinationsDefaultPagination,
	)
	if err != nil {
		return readRequestDto, err
	}

	return dto.ReadBackupDestinationsRequest{
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
	}, nil
}

func (service *BackupService) ReadDestination(input map[string]interface{}) ServiceOutput {
	readRequestDto, err := service.ReadDestinationRequestFactory(input)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	responseDto, err := useCase.ReadBackupDestinations(service.backupQueryRepo, readRequestDto)
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
			return NewServiceOutput(UserError, "InvalidMinLocalStorageFreePercent")
		}
		minLocalStorageFreePercentPtr = &minLocalStorageFreePercent
	}

	var maxDestinationStorageUsagePercentPtr *uint8
	if input["maxDestinationStorageUsagePercent"] != nil {
		maxDestinationStorageUsagePercent, err := voHelper.InterfaceToUint8(
			input["maxDestinationStorageUsagePercent"],
		)
		if err != nil {
			return NewServiceOutput(UserError, "InvalidMaxDestinationStorageUsagePercent")
		}
		maxDestinationStorageUsagePercentPtr = &maxDestinationStorageUsagePercent
	}

	var maxConcurrentConnectionsPtr *uint16
	if input["maxConcurrentConnections"] != nil {
		maxConcurrentConnections, err := voHelper.InterfaceToUint16(
			input["maxConcurrentConnections"],
		)
		if err != nil {
			return NewServiceOutput(UserError, "InvalidMaxConcurrentConnections")
		}
		maxConcurrentConnectionsPtr = &maxConcurrentConnections
	}

	var downloadBytesSecRateLimitPtr *valueObject.Byte
	if input["downloadBytesSecRateLimit"] != nil {
		downloadBytesSecRateLimit, err := valueObject.NewByte(
			input["downloadBytesSecRateLimit"],
		)
		if err != nil {
			return NewServiceOutput(UserError, "InvalidDownloadBytesSecRateLimit")
		}
		downloadBytesSecRateLimitPtr = &downloadBytesSecRateLimit
	}

	var uploadBytesSecRateLimitPtr *valueObject.Byte
	if input["uploadBytesSecRateLimit"] != nil {
		uploadBytesSecRateLimit, err := valueObject.NewByte(
			input["uploadBytesSecRateLimit"],
		)
		if err != nil {
			return NewServiceOutput(UserError, "InvalidUploadBytesSecRateLimit")
		}
		uploadBytesSecRateLimitPtr = &uploadBytesSecRateLimit
	}

	var skipCertificateVerificationPtr *bool
	if input["skipCertificateVerification"] != nil {
		skipCertificateVerification, err := voHelper.InterfaceToBool(
			input["skipCertificateVerification"],
		)
		if err != nil {
			return NewServiceOutput(UserError, "InvalidSkipCertificateVerification")
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
			return NewServiceOutput(UserError, "InvalidObjectStorageEndpointUrl")
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
			return NewServiceOutput(UserError, "InvalidRemoteHostPrivateKeyFilePath")
		}
		remoteHostPrivateKeyFilePathPtr = &remoteHostPrivateKeyFilePath
	}

	var remoteHostConnectionTimeoutSecsPtr *valueObject.TimeDuration
	if input["remoteHostConnectionTimeoutSecs"] != nil {
		remoteHostConnectionTimeoutSecs, err := valueObject.NewTimeDuration(
			input["remoteHostConnectionTimeoutSecs"],
		)
		if err != nil {
			return NewServiceOutput(UserError, "InvalidRemoteHostConnectionTimeoutSecs")
		}
		remoteHostConnectionTimeoutSecsPtr = &remoteHostConnectionTimeoutSecs
	}

	var remoteHostConnectionRetrySecsPtr *valueObject.TimeDuration
	if input["remoteHostConnectionRetrySecs"] != nil {
		remoteHostConnectionRetrySecs, err := valueObject.NewTimeDuration(
			input["remoteHostConnectionRetrySecs"],
		)
		if err != nil {
			return NewServiceOutput(UserError, "InvalidRemoteHostConnectionRetrySecs")
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
			return NewServiceOutput(UserError, "InvalidMinLocalStorageFreePercent")
		}
		minLocalStorageFreePercentPtr = &minLocalStorageFreePercent
	}

	var maxDestinationStorageUsagePercentPtr *uint8
	if input["maxDestinationStorageUsagePercent"] != nil {
		maxDestinationStorageUsagePercent, err := voHelper.InterfaceToUint8(
			input["maxDestinationStorageUsagePercent"],
		)
		if err != nil {
			return NewServiceOutput(UserError, "InvalidMaxDestinationStorageUsagePercent")
		}
		maxDestinationStorageUsagePercentPtr = &maxDestinationStorageUsagePercent
	}

	var maxConcurrentConnectionsPtr *uint16
	if input["maxConcurrentConnections"] != nil {
		maxConcurrentConnections, err := voHelper.InterfaceToUint16(
			input["maxConcurrentConnections"],
		)
		if err != nil {
			return NewServiceOutput(UserError, "InvalidMaxConcurrentConnections")
		}
		maxConcurrentConnectionsPtr = &maxConcurrentConnections
	}

	var downloadBytesSecRateLimitPtr *valueObject.Byte
	if input["downloadBytesSecRateLimit"] != nil {
		downloadBytesSecRateLimit, err := valueObject.NewByte(
			input["downloadBytesSecRateLimit"],
		)
		if err != nil {
			return NewServiceOutput(UserError, "InvalidDownloadBytesSecRateLimit")
		}
		downloadBytesSecRateLimitPtr = &downloadBytesSecRateLimit
	}

	var uploadBytesSecRateLimitPtr *valueObject.Byte
	if input["uploadBytesSecRateLimit"] != nil {
		uploadBytesSecRateLimit, err := valueObject.NewByte(
			input["uploadBytesSecRateLimit"],
		)
		if err != nil {
			return NewServiceOutput(UserError, "InvalidUploadBytesSecRateLimit")
		}
		uploadBytesSecRateLimitPtr = &uploadBytesSecRateLimit
	}

	var skipCertificateVerificationPtr *bool
	if input["skipCertificateVerification"] != nil {
		skipCertificateVerification, err := voHelper.InterfaceToBool(
			input["skipCertificateVerification"],
		)
		if err != nil {
			return NewServiceOutput(UserError, "InvalidSkipCertificateVerification")
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
			return NewServiceOutput(UserError, "InvalidObjectStorageEndpointUrl")
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
			return NewServiceOutput(UserError, "InvalidRemoteHostPrivateKeyFilePath")
		}
		remoteHostPrivateKeyFilePathPtr = &remoteHostPrivateKeyFilePath
	}

	var remoteHostConnectionTimeoutSecsPtr *valueObject.TimeDuration
	if input["remoteHostConnectionTimeoutSecs"] != nil {
		remoteHostConnectionTimeoutSecs, err := valueObject.NewTimeDuration(
			input["remoteHostConnectionTimeoutSecs"],
		)
		if err != nil {
			return NewServiceOutput(UserError, "InvalidRemoteHostConnectionTimeoutSecs")
		}
		remoteHostConnectionTimeoutSecsPtr = &remoteHostConnectionTimeoutSecs
	}

	var remoteHostConnectionRetrySecsPtr *valueObject.TimeDuration
	if input["remoteHostConnectionRetrySecs"] != nil {
		remoteHostConnectionRetrySecs, err := valueObject.NewTimeDuration(
			input["remoteHostConnectionRetrySecs"],
		)
		if err != nil {
			return NewServiceOutput(UserError, "InvalidRemoteHostConnectionRetrySecs")
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

func (service *BackupService) ReadJobRequestFactory(
	serviceInput map[string]interface{},
) (readRequestDto dto.ReadBackupJobsRequest, err error) {
	var jobIdPtr *valueObject.BackupJobId
	if serviceInput["jobId"] != nil {
		jobId, err := valueObject.NewBackupJobId(serviceInput["jobId"])
		if err != nil {
			return readRequestDto, err
		}
		jobIdPtr = &jobId
	}

	var jobStatusPtr *bool
	if serviceInput["jobStatus"] != nil {
		jobStatus, err := voHelper.InterfaceToBool(serviceInput["jobStatus"])
		if err != nil {
			return readRequestDto, errors.New("InvalidJobStatus")
		}
		jobStatusPtr = &jobStatus
	}

	var accountIdPtr *valueObject.AccountId
	if serviceInput["accountId"] != nil {
		accountId, err := valueObject.NewAccountId(serviceInput["accountId"])
		if err != nil {
			return readRequestDto, err
		}
		accountIdPtr = &accountId
	}

	var destinationIdPtr *valueObject.BackupDestinationId
	if serviceInput["destinationId"] != nil {
		destinationId, err := valueObject.NewBackupDestinationId(serviceInput["destinationId"])
		if err != nil {
			return readRequestDto, err
		}
		destinationIdPtr = &destinationId
	}

	var retentionStrategyPtr *valueObject.BackupRetentionStrategy
	if serviceInput["retentionStrategy"] != nil {
		retentionStrategy, err := valueObject.NewBackupRetentionStrategy(serviceInput["retentionStrategy"])
		if err != nil {
			return readRequestDto, err
		}
		retentionStrategyPtr = &retentionStrategy
	}

	var archiveCompressionFormatPtr *valueObject.CompressionFormat
	if serviceInput["archiveCompressionFormat"] != nil {
		archiveCompressionFormat, err := valueObject.NewCompressionFormat(serviceInput["archiveCompressionFormat"])
		if err != nil {
			return readRequestDto, err
		}
		archiveCompressionFormatPtr = &archiveCompressionFormat
	}

	var lastRunStatusPtr *valueObject.BackupTaskStatus
	if serviceInput["lastRunStatus"] != nil {
		lastRunStatus, err := valueObject.NewBackupTaskStatus(serviceInput["lastRunStatus"])
		if err != nil {
			return readRequestDto, errors.New("InvalidLastRunStatus")
		}
		lastRunStatusPtr = &lastRunStatus
	}

	timeParamNames := []string{
		"lastRunBeforeAt", "lastRunAfterAt", "nextRunBeforeAt", "nextRunAfterAt",
		"createdBeforeAt", "createdAfterAt",
	}
	timeParamPtrs := serviceHelper.TimeParamsParser(timeParamNames, serviceInput)

	requestPagination, err := serviceHelper.PaginationParser(
		serviceInput, useCase.BackupJobsDefaultPagination,
	)
	if err != nil {
		return readRequestDto, err
	}

	return dto.ReadBackupJobsRequest{
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
	}, nil
}

func (service *BackupService) ReadJob(input map[string]interface{}) ServiceOutput {
	readRequestDto, err := service.ReadJobRequestFactory(input)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	responseDto, err := useCase.ReadBackupJobs(service.backupQueryRepo, readRequestDto)
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
		return NewServiceOutput(UserError, "InvalidDestinationIds")
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
			return NewServiceOutput(UserError, "InvalidTimeoutSecs")
		}
		timeoutSecsPtr = &timeoutSecs
	}

	var maxTaskRetentionCountPtr *uint16
	if input["maxTaskRetentionCount"] != nil {
		maxTaskRetentionCount, err := voHelper.InterfaceToUint16(
			input["maxTaskRetentionCount"],
		)
		if err != nil {
			return NewServiceOutput(UserError, "InvalidMaxTaskRetentionCount")
		}
		maxTaskRetentionCountPtr = &maxTaskRetentionCount
	}

	var maxTaskRetentionDaysPtr *uint16
	if input["maxTaskRetentionDays"] != nil {
		maxTaskRetentionDays, err := voHelper.InterfaceToUint16(input["maxTaskRetentionDays"])
		if err != nil {
			return NewServiceOutput(UserError, "InvalidMaxTaskRetentionDays")
		}
		maxTaskRetentionDaysPtr = &maxTaskRetentionDays
	}

	var maxConcurrentCpuCoresPtr *uint16
	if input["maxConcurrentCpuCores"] != nil {
		maxConcurrentCpuCores, err := voHelper.InterfaceToUint16(input["maxConcurrentCpuCores"])
		if err != nil {
			return NewServiceOutput(UserError, "InvalidMaxConcurrentCpuCores")
		}
		maxConcurrentCpuCoresPtr = &maxConcurrentCpuCores
	}

	containerAccountIds := []valueObject.AccountId{}
	if input["containerAccountIds"] != nil {
		containerAccountIds, assertOk = input["containerAccountIds"].([]valueObject.AccountId)
		if !assertOk {
			return NewServiceOutput(UserError, "InvalidContainerAccountIds")
		}
	}

	containerIds := []valueObject.ContainerId{}
	if input["containerIds"] != nil {
		containerIds, assertOk = input["containerIds"].([]valueObject.ContainerId)
		if !assertOk {
			return NewServiceOutput(UserError, "InvalidContainerIds")
		}
	}

	exceptContainerAccountIds := []valueObject.AccountId{}
	if input["exceptContainerAccountIds"] != nil {
		exceptContainerAccountIds, assertOk = input["exceptContainerAccountIds"].([]valueObject.AccountId)
		if !assertOk {
			return NewServiceOutput(UserError, "InvalidExceptContainerAccountIds")
		}
	}

	exceptContainerIds := []valueObject.ContainerId{}
	if input["exceptContainerIds"] != nil {
		exceptContainerIds, assertOk = input["exceptContainerIds"].([]valueObject.ContainerId)
		if !assertOk {
			return NewServiceOutput(UserError, "InvalidExceptContainerIds")
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
			return NewServiceOutput(UserError, "InvalidJobStatus")
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
			return NewServiceOutput(UserError, "InvalidDestinationIds")
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
			return NewServiceOutput(UserError, "InvalidTimeoutSecs")
		}
		timeoutSecsPtr = &timeoutSecs
	}

	var maxTaskRetentionCountPtr *uint16
	if input["maxTaskRetentionCount"] != nil {
		maxTaskRetentionCount, err := voHelper.InterfaceToUint16(
			input["maxTaskRetentionCount"],
		)
		if err != nil {
			return NewServiceOutput(UserError, "InvalidMaxTaskRetentionCount")
		}
		maxTaskRetentionCountPtr = &maxTaskRetentionCount
	}

	var maxTaskRetentionDaysPtr *uint16
	if input["maxTaskRetentionDays"] != nil {
		maxTaskRetentionDays, err := voHelper.InterfaceToUint16(input["maxTaskRetentionDays"])
		if err != nil {
			return NewServiceOutput(UserError, "InvalidMaxTaskRetentionDays")
		}
		maxTaskRetentionDaysPtr = &maxTaskRetentionDays
	}

	var maxConcurrentCpuCoresPtr *uint16
	if input["maxConcurrentCpuCores"] != nil {
		maxConcurrentCpuCores, err := voHelper.InterfaceToUint16(input["maxConcurrentCpuCores"])
		if err != nil {
			return NewServiceOutput(UserError, "InvalidMaxConcurrentCpuCores")
		}
		maxConcurrentCpuCoresPtr = &maxConcurrentCpuCores
	}

	var containerAccountIds []valueObject.AccountId
	if input["containerAccountIds"] != nil {
		containerAccountIds, assertOk = input["containerAccountIds"].([]valueObject.AccountId)
		if !assertOk {
			return NewServiceOutput(UserError, "InvalidContainerAccountIds")
		}
	}

	var containerIds []valueObject.ContainerId
	if input["containerIds"] != nil {
		containerIds, assertOk = input["containerIds"].([]valueObject.ContainerId)
		if !assertOk {
			return NewServiceOutput(UserError, "InvalidContainerIds")
		}
	}

	var exceptContainerAccountIds []valueObject.AccountId
	if input["exceptContainerAccountIds"] != nil {
		exceptContainerAccountIds, assertOk = input["exceptContainerAccountIds"].([]valueObject.AccountId)
		if !assertOk {
			return NewServiceOutput(UserError, "InvalidExceptContainerAccountIds")
		}
	}

	var exceptContainerIds []valueObject.ContainerId
	if input["exceptContainerIds"] != nil {
		exceptContainerIds, assertOk = input["exceptContainerIds"].([]valueObject.ContainerId)
		if !assertOk {
			return NewServiceOutput(UserError, "InvalidExceptContainerIds")
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

func (service *BackupService) RunJob(
	input map[string]interface{},
	shouldSchedule bool,
) ServiceOutput {
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

	if shouldSchedule {
		cliCmd := infraEnvs.InfiniteEzBinary + " backup job run"
		cliParams := []string{
			"--account-id", accountId.String(),
			"--job-id", jobId.String(),
		}

		cliCmd = cliCmd + " " + strings.Join(cliParams, " ")

		scheduledTaskCmdRepo := infra.NewScheduledTaskCmdRepo(service.persistentDbSvc)
		taskName, _ := valueObject.NewScheduledTaskName("RunBackupJob")
		taskCmd, _ := valueObject.NewUnixCommand(cliCmd)
		taskTag, _ := valueObject.NewScheduledTaskTag("backup")
		taskTags := []valueObject.ScheduledTaskTag{taskTag}
		taskTimeoutSecs := uint32(useCase.BackupJobDefaultTimeoutSecs)

		scheduledTaskCreateDto := dto.NewCreateScheduledTask(
			taskName, taskCmd, taskTags, &taskTimeoutSecs, nil,
		)

		err = useCase.CreateScheduledTask(scheduledTaskCmdRepo, scheduledTaskCreateDto)
		if err != nil {
			return NewServiceOutput(InfraError, err.Error())
		}

		return NewServiceOutput(Created, "BackupJobRunScheduled")
	}

	runDto := dto.NewRunBackupJob(jobId, accountId, operatorAccountId, operatorIpAddress)

	err = useCase.RunBackupJob(
		service.backupQueryRepo, service.backupCmdRepo,
		service.activityRecordCmdRepo, runDto,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Created, "BackupJobRunCompleted")
}

func (service *BackupService) ReadTaskRequestFactory(
	serviceInput map[string]interface{},
) (readRequestDto dto.ReadBackupTasksRequest, err error) {
	var taskIdPtr *valueObject.BackupTaskId
	if serviceInput["taskId"] != nil {
		taskId, err := valueObject.NewBackupTaskId(serviceInput["taskId"])
		if err != nil {
			return readRequestDto, err
		}
		taskIdPtr = &taskId
	}

	var accountIdPtr *valueObject.AccountId
	if serviceInput["accountId"] != nil {
		accountId, err := valueObject.NewAccountId(serviceInput["accountId"])
		if err != nil {
			return readRequestDto, err
		}
		accountIdPtr = &accountId
	}

	var jobIdPtr *valueObject.BackupJobId
	if serviceInput["jobId"] != nil {
		jobId, err := valueObject.NewBackupJobId(serviceInput["jobId"])
		if err != nil {
			return readRequestDto, err
		}
		jobIdPtr = &jobId
	}

	var destinationIdPtr *valueObject.BackupDestinationId
	if serviceInput["destinationId"] != nil {
		destinationId, err := valueObject.NewBackupDestinationId(serviceInput["destinationId"])
		if err != nil {
			return readRequestDto, err
		}
		destinationIdPtr = &destinationId
	}

	var taskStatusPtr *valueObject.BackupTaskStatus
	if serviceInput["taskStatus"] != nil {
		taskStatus, err := valueObject.NewBackupTaskStatus(serviceInput["taskStatus"])
		if err != nil {
			return readRequestDto, err
		}
		taskStatusPtr = &taskStatus
	}

	var retentionStrategyPtr *valueObject.BackupRetentionStrategy
	if serviceInput["retentionStrategy"] != nil {
		retentionStrategy, err := valueObject.NewBackupRetentionStrategy(
			serviceInput["retentionStrategy"],
		)
		if err != nil {
			return readRequestDto, err
		}
		retentionStrategyPtr = &retentionStrategy
	}

	var containerIdPtr *valueObject.ContainerId
	if serviceInput["containerId"] != nil {
		containerId, err := valueObject.NewContainerId(serviceInput["containerId"])
		if err != nil {
			return readRequestDto, err
		}
		containerIdPtr = &containerId
	}

	timeParamNames := []string{
		"startedBeforeAt", "startedAfterAt", "finishedBeforeAt", "finishedAfterAt",
		"createdBeforeAt", "createdAfterAt",
	}
	timeParamPtrs := serviceHelper.TimeParamsParser(timeParamNames, serviceInput)

	requestPagination, err := serviceHelper.PaginationParser(
		serviceInput, useCase.BackupTasksDefaultPagination,
	)
	if err != nil {
		return readRequestDto, err
	}

	return dto.ReadBackupTasksRequest{
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
	}, nil
}

func (service *BackupService) ReadTask(input map[string]interface{}) ServiceOutput {
	readRequestDto, err := service.ReadTaskRequestFactory(input)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	responseDto, err := useCase.ReadBackupTasks(service.backupQueryRepo, readRequestDto)
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

	timeoutSecs := useCase.RestoreBackupTaskDefaultTimeoutSecs
	if input["timeoutSecs"] != nil {
		timeoutSecs, err = voHelper.InterfaceToUint32(input["timeoutSecs"])
		if err != nil {
			return NewServiceOutput(UserError, "InvalidTimeoutSecs")
		}
	}

	var assertOk bool
	var containerAccountIds []valueObject.AccountId
	if input["containerAccountIds"] != nil {
		containerAccountIds, assertOk = input["containerAccountIds"].([]valueObject.AccountId)
		if !assertOk {
			return NewServiceOutput(UserError, "InvalidAccountIds")
		}
	}

	var containerIds []valueObject.ContainerId
	if input["containerIds"] != nil {
		containerIds, assertOk = input["containerIds"].([]valueObject.ContainerId)
		if !assertOk {
			return NewServiceOutput(UserError, "InvalidContainerIds")
		}
	}

	var exceptContainerAccountIds []valueObject.AccountId
	if input["exceptContainerAccountIds"] != nil {
		exceptContainerAccountIds, assertOk = input["exceptContainerAccountIds"].([]valueObject.AccountId)
		if !assertOk {
			return NewServiceOutput(UserError, "InvalidExceptContainerAccountIds")
		}
	}

	var exceptContainerIds []valueObject.ContainerId
	if input["exceptContainerIds"] != nil {
		exceptContainerIds, assertOk = input["exceptContainerIds"].([]valueObject.ContainerId)
		if !assertOk {
			return NewServiceOutput(UserError, "InvalidExceptContainerIds")
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
			"--timeout-secs", strconv.Itoa(int(timeoutSecs)),
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
			taskName, taskCmd, taskTags, &timeoutSecs, nil,
		)

		err = useCase.CreateScheduledTask(scheduledTaskCmdRepo, scheduledTaskCreateDto)
		if err != nil {
			return NewServiceOutput(InfraError, err.Error())
		}

		return NewServiceOutput(Created, "BackupTaskRestoreScheduled")
	}

	restoreRequestDto := dto.RestoreBackupTaskRequest{
		TaskId:                          taskIdPtr,
		ArchiveId:                       archiveIdPtr,
		ShouldReplaceExistingContainers: &shouldReplaceExistingContainers,
		ShouldRestoreMappings:           &shouldRestoreMappings,
		TimeoutSecs:                     &timeoutSecs,
		ContainerAccountIds:             containerAccountIds,
		ContainerIds:                    containerIds,
		ExceptContainerAccountIds:       exceptContainerAccountIds,
		ExceptContainerIds:              exceptContainerIds,
		OperatorAccountId:               operatorAccountId,
		OperatorIpAddress:               operatorIpAddress,
	}

	restoreResponseDto, err := useCase.RestoreBackupTask(
		service.backupQueryRepo, service.backupCmdRepo,
		service.activityRecordCmdRepo, restoreRequestDto,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	responseStatusEnum := Success
	if len(restoreResponseDto.FailedContainerImageIds) > 0 {
		responseStatusEnum = MultiStatus
	}
	if len(restoreResponseDto.SuccessfulContainerIds) == 0 {
		responseStatusEnum = InfraError
	}

	return NewServiceOutput(responseStatusEnum, restoreResponseDto)
}

func (service *BackupService) UpdateTask(input map[string]interface{}) ServiceOutput {
	requiredParams := []string{"taskId"}
	err := serviceHelper.RequiredParamsInspector(input, requiredParams)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	taskId, err := valueObject.NewBackupTaskId(input["taskId"])
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

	var taskStatusPtr *valueObject.BackupTaskStatus
	if input["taskStatus"] != nil {
		taskStatus, err := valueObject.NewBackupTaskStatus(input["taskStatus"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		taskStatusPtr = &taskStatus
	}

	updateDto := dto.NewUpdateBackupTask(
		taskId, taskStatusPtr, operatorAccountId, operatorIpAddress,
	)

	err = useCase.UpdateBackupTask(
		service.backupQueryRepo, service.backupCmdRepo,
		service.activityRecordCmdRepo, updateDto,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, "BackupTaskUpdated")
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

func (service *BackupService) ReadTaskArchiveRequestFactory(
	serviceInput map[string]interface{},
) (readRequestDto dto.ReadBackupTaskArchivesRequest, err error) {
	var archiveIdPtr *valueObject.BackupTaskArchiveId
	if serviceInput["archiveId"] != nil {
		archiveId, err := valueObject.NewBackupTaskArchiveId(serviceInput["archiveId"])
		if err != nil {
			return readRequestDto, err
		}
		archiveIdPtr = &archiveId
	}

	var accountIdPtr *valueObject.AccountId
	if serviceInput["accountId"] != nil {
		accountId, err := valueObject.NewAccountId(serviceInput["accountId"])
		if err != nil {
			return readRequestDto, err
		}
		accountIdPtr = &accountId
	}

	var taskIdPtr *valueObject.BackupTaskId
	if serviceInput["taskId"] != nil {
		taskId, err := valueObject.NewBackupTaskId(serviceInput["taskId"])
		if err != nil {
			return readRequestDto, err
		}
		taskIdPtr = &taskId
	}

	timeParamNames := []string{"createdBeforeAt", "createdAfterAt"}
	timeParamPtrs := serviceHelper.TimeParamsParser(timeParamNames, serviceInput)

	requestPagination, err := serviceHelper.PaginationParser(
		serviceInput, useCase.BackupTaskArchivesDefaultPagination,
	)
	if err != nil {
		return readRequestDto, err
	}

	return dto.ReadBackupTaskArchivesRequest{
		Pagination:      requestPagination,
		ArchiveId:       archiveIdPtr,
		AccountId:       accountIdPtr,
		TaskId:          taskIdPtr,
		CreatedBeforeAt: timeParamPtrs["createdBeforeAt"],
		CreatedAfterAt:  timeParamPtrs["createdAfterAt"],
	}, nil
}

func (service *BackupService) ReadTaskArchive(
	input map[string]interface{},
	requestHostname *string,
) ServiceOutput {
	readRequestDto, err := service.ReadTaskArchiveRequestFactory(input)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	responseDto, err := useCase.ReadBackupTaskArchives(service.backupQueryRepo, readRequestDto)
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

	timeoutSecs := useCase.CreateBackupTaskArchiveDefaultTimeoutSecs
	if input["timeoutSecs"] != nil {
		timeoutSecs, err = voHelper.InterfaceToUint32(input["timeoutSecs"])
		if err != nil {
			return NewServiceOutput(UserError, "InvalidTimeoutSecs")
		}
	}

	var assertOk bool
	var containerAccountIds []valueObject.AccountId
	if input["containerAccountIds"] != nil {
		containerAccountIds, assertOk = input["containerAccountIds"].([]valueObject.AccountId)
		if !assertOk {
			return NewServiceOutput(UserError, "InvalidAccountIds")
		}
	}

	var containerIds []valueObject.ContainerId
	if input["containerIds"] != nil {
		containerIds, assertOk = input["containerIds"].([]valueObject.ContainerId)
		if !assertOk {
			return NewServiceOutput(UserError, "InvalidContainerIds")
		}
	}

	var exceptContainerAccountIds []valueObject.AccountId
	if input["exceptContainerAccountIds"] != nil {
		exceptContainerAccountIds, assertOk = input["exceptContainerAccountIds"].([]valueObject.AccountId)
		if !assertOk {
			return NewServiceOutput(UserError, "InvalidExceptContainerAccountIds")
		}
	}

	var exceptContainerIds []valueObject.ContainerId
	if input["exceptContainerIds"] != nil {
		exceptContainerIds, assertOk = input["exceptContainerIds"].([]valueObject.ContainerId)
		if !assertOk {
			return NewServiceOutput(UserError, "InvalidExceptContainerIds")
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
			"--timeout-secs", strconv.Itoa(int(timeoutSecs)),
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
			taskName, taskCmd, taskTags, &timeoutSecs, nil,
		)

		err = useCase.CreateScheduledTask(scheduledTaskCmdRepo, scheduledTaskCreateDto)
		if err != nil {
			return NewServiceOutput(InfraError, err.Error())
		}

		return NewServiceOutput(Created, "CreateBackupTaskArchiveScheduled")
	}

	createDto := dto.NewCreateBackupTaskArchive(
		taskId, &timeoutSecs, containerAccountIds, containerIds, exceptContainerAccountIds,
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
