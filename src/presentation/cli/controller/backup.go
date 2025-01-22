package cliController

import (
	"github.com/goinfinite/ez/src/domain/valueObject"
	"github.com/goinfinite/ez/src/infra/db"
	cliHelper "github.com/goinfinite/ez/src/presentation/cli/helper"
	"github.com/goinfinite/ez/src/presentation/service"
	sharedHelper "github.com/goinfinite/ez/src/presentation/shared/helper"
	"github.com/spf13/cobra"
)

type BackupController struct {
	backupService *service.BackupService
}

func NewBackupController(
	persistentDbSvc *db.PersistentDatabaseService,
	trailDbSvc *db.TrailDatabaseService,
) *BackupController {
	return &BackupController{
		backupService: service.NewBackupService(persistentDbSvc, trailDbSvc),
	}
}

func (controller *BackupController) ReadDestination() *cobra.Command {
	var destinationIdUint, accountIdUint uint64
	var destinationNameStr, destinationTypeStr string
	var objectStorageProviderStr, remoteHostTypeStr, remoteHostnameStr string
	var paginationPageNumberUint32 uint32
	var paginationItemsPerPageUint16 uint16
	var paginationSortByStr, paginationSortDirectionStr string
	var paginationLastSeenIdStr string

	cmd := &cobra.Command{
		Use:   "get",
		Short: "ReadBackupsDestinations",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{}

			if destinationIdUint != 0 {
				requestBody["destinationId"] = destinationIdUint
			}

			if accountIdUint != 0 {
				requestBody["accountId"] = accountIdUint
			}

			if destinationNameStr != "" {
				requestBody["destinationName"] = destinationNameStr
			}

			if destinationTypeStr != "" {
				requestBody["destinationType"] = destinationTypeStr
			}

			if objectStorageProviderStr != "" {
				requestBody["objectStorageProvider"] = objectStorageProviderStr
			}

			if remoteHostTypeStr != "" {
				requestBody["remoteHostType"] = remoteHostTypeStr
			}

			if remoteHostnameStr != "" {
				requestBody["remoteHostname"] = remoteHostnameStr
			}

			if paginationPageNumberUint32 != 0 {
				requestBody["pageNumber"] = paginationPageNumberUint32
			}
			if paginationItemsPerPageUint16 != 0 {
				requestBody["itemsPerPage"] = paginationItemsPerPageUint16
			}
			if paginationSortByStr != "" {
				requestBody["sortBy"] = paginationSortByStr
			}
			if paginationSortDirectionStr != "" {
				requestBody["sortDirection"] = paginationSortDirectionStr
			}
			if paginationLastSeenIdStr != "" {
				requestBody["lastSeenId"] = paginationLastSeenIdStr
			}

			cliHelper.ServiceResponseWrapper(
				controller.backupService.ReadDestination(requestBody),
			)
		},
	}

	cmd.Flags().Uint64VarP(
		&destinationIdUint, "destination-id", "d", 0, "BackupDestinationId",
	)
	cmd.Flags().Uint64VarP(&accountIdUint, "account-id", "a", 0, "AccountId")
	cmd.Flags().StringVarP(&destinationNameStr, "destinationName", "n", "", "BackupDestinationName")
	cmd.Flags().StringVarP(&destinationTypeStr, "destinationType", "t", "", "BackupDestinationType")
	cmd.Flags().StringVarP(
		&objectStorageProviderStr, "objectStorageProvider", "p", "", "ObjectStorageProvider",
	)
	cmd.Flags().StringVarP(
		&remoteHostTypeStr, "remoteHostType", "r", "", "RemoteHostType (sftp|ftp)",
	)
	cmd.Flags().StringVarP(&remoteHostnameStr, "remoteHostname", "H", "", "RemoteHostname")
	cmd.Flags().Uint32VarP(
		&paginationPageNumberUint32, "page-number", "o", 0, "PageNumber (Pagination)",
	)
	cmd.Flags().Uint16VarP(
		&paginationItemsPerPageUint16, "items-per-page", "j", 0, "ItemsPerPage (Pagination)",
	)
	cmd.Flags().StringVarP(
		&paginationSortByStr, "sort-by", "y", "", "SortBy (Pagination)",
	)
	cmd.Flags().StringVarP(
		&paginationSortDirectionStr, "sort-direction", "x", "", "SortDirection (Pagination)",
	)
	cmd.Flags().StringVarP(
		&paginationLastSeenIdStr, "last-seen-id", "l", "", "LastSeenId (Pagination)",
	)

	return cmd
}

func (controller *BackupController) CreateDestination() *cobra.Command {
	var accountIdUint uint64
	var minLocalStorageFreePercentUint8, maxDestinationStorageUsagePercentUint8 uint8
	var maxConcurrentConnectionsUint16 uint16
	var downloadBytesSecRateLimitUint64, uploadBytesSecRateLimitUint64 uint64
	var destinationNameStr, destinationDescriptionStr, destinationTypeStr, destinationPathStr string
	var skipCertificateVerificationBoolStr string
	var objectStorageProviderStr, objectStorageProviderRegionStr, objectStorageProviderAccessKeyIdStr string
	var objectStorageProviderSecretAccessKeyStr, objectStorageEndpointUrlStr, objectStorageBucketNameStr string
	var remoteHostTypeStr, remoteHostnameStr, remoteHostUsernameStr string
	var remoteHostPasswordStr, remoteHostPrivateKeyFilePathStr string
	var remoteHostNetworkPortUint16, remoteHostConnectionTimeoutSecsUint16 uint16
	var remoteHostConnectionRetrySecsUint16 uint16

	cmd := &cobra.Command{
		Use:   "create",
		Short: "CreateBackupDestination",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{
				"accountId":                   accountIdUint,
				"destinationName":             destinationNameStr,
				"destinationType":             destinationTypeStr,
				"skipCertificateVerification": skipCertificateVerificationBoolStr,
			}

			if destinationDescriptionStr != "" {
				requestBody["destinationDescription"] = destinationDescriptionStr
			}

			if destinationPathStr != "" {
				requestBody["destinationPath"] = destinationPathStr
			}

			if minLocalStorageFreePercentUint8 != 0 {
				requestBody["minLocalStorageFreePercent"] = minLocalStorageFreePercentUint8
			}

			if maxDestinationStorageUsagePercentUint8 != 0 {
				requestBody["maxDestinationStorageUsagePercent"] = maxDestinationStorageUsagePercentUint8
			}

			if maxConcurrentConnectionsUint16 != 0 {
				requestBody["maxConcurrentConnections"] = maxConcurrentConnectionsUint16
			}

			if downloadBytesSecRateLimitUint64 != 0 {
				requestBody["downloadBytesSecRateLimit"] = downloadBytesSecRateLimitUint64
			}

			if uploadBytesSecRateLimitUint64 != 0 {
				requestBody["uploadBytesSecRateLimit"] = uploadBytesSecRateLimitUint64
			}

			if objectStorageProviderStr != "" {
				requestBody["objectStorageProvider"] = objectStorageProviderStr
			}

			if objectStorageProviderRegionStr != "" {
				requestBody["objectStorageProviderRegion"] = objectStorageProviderRegionStr
			}

			if objectStorageProviderAccessKeyIdStr != "" {
				requestBody["objectStorageProviderAccessKeyId"] = objectStorageProviderAccessKeyIdStr
			}

			if objectStorageProviderSecretAccessKeyStr != "" {
				requestBody["objectStorageProviderSecretAccessKey"] = objectStorageProviderSecretAccessKeyStr
			}

			if objectStorageEndpointUrlStr != "" {
				requestBody["objectStorageEndpointUrl"] = objectStorageEndpointUrlStr
			}

			if objectStorageBucketNameStr != "" {
				requestBody["objectStorageBucketName"] = objectStorageBucketNameStr
			}

			if remoteHostTypeStr != "" {
				requestBody["remoteHostType"] = remoteHostTypeStr
			}

			if remoteHostnameStr != "" {
				requestBody["remoteHostname"] = remoteHostnameStr
			}

			if remoteHostNetworkPortUint16 != 0 {
				requestBody["remoteHostNetworkPort"] = remoteHostNetworkPortUint16
			}

			if remoteHostUsernameStr != "" {
				requestBody["remoteHostUsername"] = remoteHostUsernameStr
			}

			if remoteHostPasswordStr != "" {
				requestBody["remoteHostPassword"] = remoteHostPasswordStr
			}

			if remoteHostPrivateKeyFilePathStr != "" {
				requestBody["remoteHostPrivateKeyFilePath"] = remoteHostPrivateKeyFilePathStr
			}

			if remoteHostConnectionTimeoutSecsUint16 != 0 {
				requestBody["remoteHostConnectionTimeoutSecs"] = remoteHostConnectionTimeoutSecsUint16
			}

			if remoteHostConnectionRetrySecsUint16 != 0 {
				requestBody["remoteHostConnectionRetrySecs"] = remoteHostConnectionRetrySecsUint16
			}

			cliHelper.ServiceResponseWrapper(
				controller.backupService.CreateDestination(requestBody),
			)
		},
	}

	cmd.Flags().Uint64VarP(&accountIdUint, "account-id", "a", 0, "AccountId")
	cmd.MarkFlagRequired("account-id")
	cmd.Flags().StringVarP(
		&destinationNameStr, "destination-name", "n", "", "BackupDestinationName",
	)
	cmd.MarkFlagRequired("destination-name")
	cmd.Flags().StringVarP(
		&destinationDescriptionStr, "destination-description", "D", "", "BackupDestinationDescription",
	)
	cmd.Flags().StringVarP(
		&destinationTypeStr, "destination-type", "t", "", "BackupDestinationType (object-storage|remote-host|local)",
	)
	cmd.MarkFlagRequired("destination-type")
	cmd.Flags().StringVarP(
		&destinationPathStr, "destination-path", "p", "", "BackupDestinationPath",
	)
	cmd.Flags().Uint8VarP(
		&minLocalStorageFreePercentUint8, "min-local-storage-free-percent", "m", 0,
		"MinLocalStorageFreePercent",
	)
	cmd.Flags().Uint8VarP(
		&maxDestinationStorageUsagePercentUint8, "max-destination-storage-usage-percent",
		"M", 0, "MaxDestinationStorageUsagePercent",
	)
	cmd.Flags().Uint16VarP(
		&maxConcurrentConnectionsUint16, "max-concurrent-connections", "c", 0,
		"MaxConcurrentConnections",
	)
	cmd.Flags().Uint64VarP(
		&downloadBytesSecRateLimitUint64, "download-bytes-sec-rate-limit", "d", 0,
		"DownloadBytesSecRateLimit",
	)
	cmd.Flags().Uint64VarP(
		&uploadBytesSecRateLimitUint64, "upload-bytes-sec-rate-limit", "u", 0,
		"UploadBytesSecRateLimit",
	)
	cmd.Flags().StringVarP(
		&skipCertificateVerificationBoolStr, "skip-certificate-verification", "s",
		"false", "SkipCertificateVerification",
	)
	cmd.Flags().StringVarP(
		&objectStorageProviderStr, "object-storage-provider", "o", "", "ObjectStorageProvider",
	)
	cmd.Flags().StringVarP(
		&objectStorageProviderRegionStr, "object-storage-provider-region", "r",
		"", "ObjectStorageProviderRegion",
	)
	cmd.Flags().StringVarP(
		&objectStorageProviderAccessKeyIdStr, "object-storage-provider-access-key-id", "k",
		"", "ObjectStorageProviderAccessKeyId",
	)
	cmd.Flags().StringVarP(
		&objectStorageProviderSecretAccessKeyStr, "object-storage-provider-secret-access-key",
		"K", "", "ObjectStorageProviderSecretAccessKey",
	)
	cmd.Flags().StringVarP(
		&objectStorageEndpointUrlStr, "object-storage-endpoint-url", "e", "", "ObjectStorageEndpointUrl",
	)
	cmd.Flags().StringVarP(
		&objectStorageBucketNameStr, "object-storage-bucket-name", "b", "", "ObjectStorageBucketName",
	)
	cmd.Flags().StringVarP(
		&remoteHostTypeStr, "remote-host-type", "R", "", "RemoteHostType (sftp|ftp)",
	)
	cmd.Flags().StringVarP(
		&remoteHostnameStr, "remote-hostname", "H", "", "RemoteHostname",
	)
	cmd.Flags().Uint16VarP(
		&remoteHostNetworkPortUint16, "remote-host-network-port", "N", 0, "RemoteHostNetworkPort",
	)
	cmd.Flags().StringVarP(
		&remoteHostUsernameStr, "remote-host-username", "U", "", "RemoteHostUsername",
	)
	cmd.Flags().StringVarP(
		&remoteHostPasswordStr, "remote-host-password", "P", "", "RemoteHostPassword",
	)
	cmd.Flags().StringVarP(
		&remoteHostPrivateKeyFilePathStr, "remote-host-private-key-file-path", "f", "", "RemoteHostPrivateKeyFilePath",
	)
	cmd.Flags().Uint16VarP(
		&remoteHostConnectionTimeoutSecsUint16, "remote-host-connection-timeout-secs", "T", 0, "RemoteHostConnectionTimeoutSecs",
	)
	cmd.Flags().Uint16VarP(
		&remoteHostConnectionRetrySecsUint16, "remote-host-connection-retry-secs", "y", 0, "RemoteHostConnectionRetrySecs",
	)

	return cmd
}

func (controller *BackupController) UpdateDestination() *cobra.Command {
	var destinationIdUint, accountIdUint uint64
	var minLocalStorageFreePercentUint8, maxDestinationStorageUsagePercentUint8 uint8
	var maxConcurrentConnectionsUint16 uint16
	var downloadBytesSecRateLimitUint64, uploadBytesSecRateLimitUint64 uint64
	var destinationNameStr, destinationDescriptionStr, destinationTypeStr, destinationPathStr string
	var skipCertificateVerificationBoolStr string
	var objectStorageProviderStr, objectStorageProviderRegionStr, objectStorageProviderAccessKeyIdStr string
	var objectStorageProviderSecretAccessKeyStr, objectStorageEndpointUrlStr, objectStorageBucketNameStr string
	var remoteHostTypeStr, remoteHostnameStr, remoteHostUsernameStr string
	var remoteHostPasswordStr, remoteHostPrivateKeyFilePathStr string
	var remoteHostNetworkPortUint16, remoteHostConnectionTimeoutSecsUint16 uint16
	var remoteHostConnectionRetrySecsUint16 uint16

	cmd := &cobra.Command{
		Use:   "update",
		Short: "updateBackupDestination",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{
				"destinationId":               destinationIdUint,
				"accountId":                   accountIdUint,
				"skipCertificateVerification": skipCertificateVerificationBoolStr,
			}

			if destinationNameStr != "" {
				requestBody["destinationName"] = destinationNameStr
			}

			if destinationTypeStr != "" {
				requestBody["destinationType"] = destinationTypeStr
			}

			if destinationDescriptionStr != "" {
				requestBody["destinationDescription"] = destinationDescriptionStr
			}

			if destinationPathStr != "" {
				requestBody["destinationPath"] = destinationPathStr
			}

			if minLocalStorageFreePercentUint8 != 0 {
				requestBody["minLocalStorageFreePercent"] = minLocalStorageFreePercentUint8
			}

			if maxDestinationStorageUsagePercentUint8 != 0 {
				requestBody["maxDestinationStorageUsagePercent"] = maxDestinationStorageUsagePercentUint8
			}

			if maxConcurrentConnectionsUint16 != 0 {
				requestBody["maxConcurrentConnections"] = maxConcurrentConnectionsUint16
			}

			if downloadBytesSecRateLimitUint64 != 0 {
				requestBody["downloadBytesSecRateLimit"] = downloadBytesSecRateLimitUint64
			}

			if uploadBytesSecRateLimitUint64 != 0 {
				requestBody["uploadBytesSecRateLimit"] = uploadBytesSecRateLimitUint64
			}

			if objectStorageProviderStr != "" {
				requestBody["objectStorageProvider"] = objectStorageProviderStr
			}

			if objectStorageProviderRegionStr != "" {
				requestBody["objectStorageProviderRegion"] = objectStorageProviderRegionStr
			}

			if objectStorageProviderAccessKeyIdStr != "" {
				requestBody["objectStorageProviderAccessKeyId"] = objectStorageProviderAccessKeyIdStr
			}

			if objectStorageProviderSecretAccessKeyStr != "" {
				requestBody["objectStorageProviderSecretAccessKey"] = objectStorageProviderSecretAccessKeyStr
			}

			if objectStorageEndpointUrlStr != "" {
				requestBody["objectStorageEndpointUrl"] = objectStorageEndpointUrlStr
			}

			if objectStorageBucketNameStr != "" {
				requestBody["objectStorageBucketName"] = objectStorageBucketNameStr
			}

			if remoteHostTypeStr != "" {
				requestBody["remoteHostType"] = remoteHostTypeStr
			}

			if remoteHostnameStr != "" {
				requestBody["remoteHostname"] = remoteHostnameStr
			}

			if remoteHostNetworkPortUint16 != 0 {
				requestBody["remoteHostNetworkPort"] = remoteHostNetworkPortUint16
			}

			if remoteHostUsernameStr != "" {
				requestBody["remoteHostUsername"] = remoteHostUsernameStr
			}

			if remoteHostPasswordStr != "" {
				requestBody["remoteHostPassword"] = remoteHostPasswordStr
			}

			if remoteHostPrivateKeyFilePathStr != "" {
				requestBody["remoteHostPrivateKeyFilePath"] = remoteHostPrivateKeyFilePathStr
			}

			if remoteHostConnectionTimeoutSecsUint16 != 0 {
				requestBody["remoteHostConnectionTimeoutSecs"] = remoteHostConnectionTimeoutSecsUint16
			}

			if remoteHostConnectionRetrySecsUint16 != 0 {
				requestBody["remoteHostConnectionRetrySecs"] = remoteHostConnectionRetrySecsUint16
			}

			cliHelper.ServiceResponseWrapper(
				controller.backupService.UpdateDestination(requestBody),
			)
		},
	}

	cmd.Flags().Uint64VarP(&destinationIdUint, "destination-id", "d", 0, "BackupDestinationId")
	cmd.MarkFlagRequired("destination-id")
	cmd.Flags().Uint64VarP(&accountIdUint, "account-id", "a", 0, "AccountId")
	cmd.MarkFlagRequired("account-id")
	cmd.Flags().StringVarP(
		&destinationNameStr, "destination-name", "n", "", "BackupDestinationName",
	)
	cmd.Flags().StringVarP(
		&destinationDescriptionStr, "destination-description", "D", "", "BackupDestinationDescription",
	)
	cmd.Flags().StringVarP(
		&destinationTypeStr, "destination-type", "t", "", "BackupDestinationType (object-storage|remote-host|local)",
	)
	cmd.Flags().StringVarP(
		&destinationPathStr, "destination-path", "p", "", "BackupDestinationPath",
	)
	cmd.Flags().Uint8VarP(
		&minLocalStorageFreePercentUint8, "min-local-storage-free-percent", "m", 0,
		"MinLocalStorageFreePercent",
	)
	cmd.Flags().Uint8VarP(
		&maxDestinationStorageUsagePercentUint8, "max-destination-storage-usage-percent",
		"M", 0, "MaxDestinationStorageUsagePercent",
	)
	cmd.Flags().Uint16VarP(
		&maxConcurrentConnectionsUint16, "max-concurrent-connections", "c", 0,
		"MaxConcurrentConnections",
	)
	cmd.Flags().Uint64VarP(
		&downloadBytesSecRateLimitUint64, "download-bytes-sec-rate-limit", "B", 0,
		"DownloadBytesSecRateLimit",
	)
	cmd.Flags().Uint64VarP(
		&uploadBytesSecRateLimitUint64, "upload-bytes-sec-rate-limit", "U", 0,
		"UploadBytesSecRateLimit",
	)
	cmd.Flags().StringVarP(
		&skipCertificateVerificationBoolStr, "skip-certificate-verification", "s",
		"false", "SkipCertificateVerification",
	)
	cmd.Flags().StringVarP(
		&objectStorageProviderStr, "object-storage-provider", "o", "", "ObjectStorageProvider",
	)
	cmd.Flags().StringVarP(
		&objectStorageProviderRegionStr, "object-storage-provider-region", "r",
		"", "ObjectStorageProviderRegion",
	)
	cmd.Flags().StringVarP(
		&objectStorageProviderAccessKeyIdStr, "object-storage-provider-access-key-id", "k",
		"", "ObjectStorageProviderAccessKeyId",
	)
	cmd.Flags().StringVarP(
		&objectStorageProviderSecretAccessKeyStr, "object-storage-provider-secret-access-key",
		"K", "", "ObjectStorageProviderSecretAccessKey",
	)
	cmd.Flags().StringVarP(
		&objectStorageEndpointUrlStr, "object-storage-endpoint-url", "e", "", "ObjectStorageEndpointUrl",
	)
	cmd.Flags().StringVarP(
		&objectStorageBucketNameStr, "object-storage-bucket-name", "b", "", "ObjectStorageBucketName",
	)
	cmd.Flags().StringVarP(
		&remoteHostTypeStr, "remote-host-type", "R", "", "RemoteHostType (sftp|ftp)",
	)
	cmd.Flags().StringVarP(
		&remoteHostnameStr, "remote-hostname", "H", "", "RemoteHostname",
	)
	cmd.Flags().Uint16VarP(
		&remoteHostNetworkPortUint16, "remote-host-network-port", "N", 0, "RemoteHostNetworkPort",
	)
	cmd.Flags().StringVarP(
		&remoteHostUsernameStr, "remote-host-username", "u", "", "RemoteHostUsername",
	)
	cmd.Flags().StringVarP(
		&remoteHostPasswordStr, "remote-host-password", "P", "", "RemoteHostPassword",
	)
	cmd.Flags().StringVarP(
		&remoteHostPrivateKeyFilePathStr, "remote-host-private-key-file-path", "f", "", "RemoteHostPrivateKeyFilePath",
	)
	cmd.Flags().Uint16VarP(
		&remoteHostConnectionTimeoutSecsUint16, "remote-host-connection-timeout-secs", "T", 0, "RemoteHostConnectionTimeoutSecs",
	)
	cmd.Flags().Uint16VarP(
		&remoteHostConnectionRetrySecsUint16, "remote-host-connection-retry-secs", "y", 0, "RemoteHostConnectionRetrySecs",
	)

	return cmd
}

func (controller *BackupController) DeleteDestination() *cobra.Command {
	var destinationIdUint, accountIdUint uint64

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "DeleteBackupDestination",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{
				"destinationId": destinationIdUint,
				"accountId":     accountIdUint,
			}

			cliHelper.ServiceResponseWrapper(
				controller.backupService.DeleteDestination(requestBody),
			)
		},
	}

	cmd.Flags().Uint64VarP(&destinationIdUint, "destination-id", "d", 0, "BackupDestinationId")
	cmd.MarkFlagRequired("destination-id")
	cmd.Flags().Uint64VarP(&accountIdUint, "account-id", "a", 0, "AccountId")
	cmd.MarkFlagRequired("account-id")
	return cmd
}

func (controller *BackupController) ReadJob() *cobra.Command {
	var jobIdUint, accountIdUint, destinationIdUint uint64
	var retentionStrategyStr, jobStatusStr, archiveCompressionFormatStr, lastRunStatusStr string
	var paginationPageNumberUint32 uint32
	var paginationItemsPerPageUint16 uint16
	var paginationSortByStr, paginationSortDirectionStr string
	var paginationLastSeenIdStr string

	cmd := &cobra.Command{
		Use:   "get",
		Short: "ReadBackupJobs",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{}

			if jobIdUint != 0 {
				requestBody["jobId"] = jobIdUint
			}

			if jobStatusStr != "" {
				requestBody["jobStatus"] = jobStatusStr
			}

			if accountIdUint != 0 {
				requestBody["accountId"] = accountIdUint
			}

			if destinationIdUint != 0 {
				requestBody["destinationId"] = destinationIdUint
			}

			if retentionStrategyStr != "" {
				requestBody["retentionStrategy"] = retentionStrategyStr
			}

			if archiveCompressionFormatStr != "" {
				requestBody["archiveCompressionFormat"] = archiveCompressionFormatStr
			}

			if lastRunStatusStr != "" {
				requestBody["lastRunStatus"] = lastRunStatusStr
			}

			if paginationPageNumberUint32 != 0 {
				requestBody["pageNumber"] = paginationPageNumberUint32
			}
			if paginationItemsPerPageUint16 != 0 {
				requestBody["itemsPerPage"] = paginationItemsPerPageUint16
			}
			if paginationSortByStr != "" {
				requestBody["sortBy"] = paginationSortByStr
			}
			if paginationSortDirectionStr != "" {
				requestBody["sortDirection"] = paginationSortDirectionStr
			}
			if paginationLastSeenIdStr != "" {
				requestBody["lastSeenId"] = paginationLastSeenIdStr
			}

			cliHelper.ServiceResponseWrapper(
				controller.backupService.ReadJob(requestBody),
			)
		},
	}

	cmd.Flags().Uint64VarP(&jobIdUint, "job-id", "j", 0, "BackupJobId")
	cmd.Flags().StringVarP(&jobStatusStr, "job-status", "s", "", "BackupJobStatus")
	cmd.Flags().Uint64VarP(&accountIdUint, "account-id", "a", 0, "BackupAccountId")
	cmd.Flags().Uint64VarP(&destinationIdUint, "destination-id", "d", 0, "BackupDestinationId")
	cmd.Flags().StringVarP(&retentionStrategyStr, "retention-strategy", "r", "", "RetentionStrategy")
	cmd.Flags().StringVarP(
		&archiveCompressionFormatStr, "archive-compression-format", "c", "", "ArchiveCompressionFormat",
	)
	cmd.Flags().StringVarP(&lastRunStatusStr, "last-run-status", "l", "", "LastRunStatus")
	cmd.Flags().Uint32VarP(
		&paginationPageNumberUint32, "page-number", "o", 0, "PageNumber (Pagination)",
	)
	cmd.Flags().Uint16VarP(
		&paginationItemsPerPageUint16, "items-per-page", "I", 0, "ItemsPerPage (Pagination)",
	)
	cmd.Flags().StringVarP(
		&paginationSortByStr, "sort-by", "y", "", "SortBy (Pagination)",
	)
	cmd.Flags().StringVarP(
		&paginationSortDirectionStr, "sort-direction", "x", "", "SortDirection (Pagination)",
	)
	cmd.Flags().StringVarP(
		&paginationLastSeenIdStr, "last-seen-id", "L", "", "LastSeenId (Pagination)",
	)

	return cmd
}

func (controller *BackupController) CreateJob() *cobra.Command {
	var accountIdUint, timeoutSecsUint uint64
	var maxTaskRetentionCountUint, maxTaskRetentionDaysUint, maxConcurrentCpuCoresUint uint16
	var jobDescriptionStr, retentionStrategyStr, backupScheduleStr string
	var archiveCompressionFormatStr string
	var destinationIdsSlice, containerAccountIdsSlice, exceptContainerAccountIdsSlice []string
	var containerIdsSlice, exceptContainerIdsSlice []string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "CreateBackupJob",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{
				"accountId":      accountIdUint,
				"backupSchedule": backupScheduleStr,
			}

			requestBody["destinationIds"] = sharedHelper.StringSliceValueObjectParser(
				destinationIdsSlice, valueObject.NewBackupDestinationId,
			)

			if jobDescriptionStr != "" {
				requestBody["jobDescription"] = jobDescriptionStr
			}

			if retentionStrategyStr != "" {
				requestBody["retentionStrategy"] = retentionStrategyStr
			}

			if archiveCompressionFormatStr != "" {
				requestBody["archiveCompressionFormat"] = archiveCompressionFormatStr
			}

			if timeoutSecsUint != 0 {
				requestBody["timeoutSecs"] = timeoutSecsUint
			}

			if maxTaskRetentionCountUint != 0 {
				requestBody["maxTaskRetentionCount"] = maxTaskRetentionCountUint
			}

			if maxTaskRetentionDaysUint != 0 {
				requestBody["maxTaskRetentionDays"] = maxTaskRetentionDaysUint
			}

			if maxConcurrentCpuCoresUint != 0 {
				requestBody["maxConcurrentCpuCores"] = maxConcurrentCpuCoresUint
			}

			if len(containerAccountIdsSlice) > 0 {
				requestBody["containerAccountIds"] = sharedHelper.StringSliceValueObjectParser(
					containerAccountIdsSlice, valueObject.NewAccountId,
				)
			}

			if len(containerIdsSlice) > 0 {
				requestBody["containerIds"] = sharedHelper.StringSliceValueObjectParser(
					containerIdsSlice, valueObject.NewContainerId,
				)
			}

			if len(exceptContainerAccountIdsSlice) > 0 {
				requestBody["exceptContainerAccountIds"] = sharedHelper.StringSliceValueObjectParser(
					exceptContainerAccountIdsSlice, valueObject.NewAccountId,
				)
			}

			if len(exceptContainerIdsSlice) > 0 {
				requestBody["exceptContainerIds"] = sharedHelper.StringSliceValueObjectParser(
					exceptContainerIdsSlice, valueObject.NewContainerId,
				)
			}

			cliHelper.ServiceResponseWrapper(
				controller.backupService.CreateJob(requestBody),
			)
		},
	}

	cmd.Flags().Uint64VarP(&accountIdUint, "account-id", "a", 0, "AccountId")
	cmd.MarkFlagRequired("account-id")
	cmd.Flags().StringVarP(&jobDescriptionStr, "job-description", "j", "", "BackupJobDescription")
	cmd.Flags().StringSliceVarP(
		&destinationIdsSlice, "destination-ids", "d", []string{}, "BackupDestinationIds",
	)
	cmd.MarkFlagRequired("destination-ids")
	cmd.Flags().StringVarP(&retentionStrategyStr, "retention-strategy", "r", "", "RetentionStrategy")
	cmd.Flags().StringVarP(&backupScheduleStr, "backup-schedule", "s", "", "BackupSchedule")
	cmd.MarkFlagRequired("backup-schedule")
	cmd.Flags().StringVarP(
		&archiveCompressionFormatStr, "archive-compression-format", "c", "", "ArchiveCompressionFormat",
	)
	cmd.Flags().Uint64VarP(&timeoutSecsUint, "timeout-secs", "t", 0, "TimeoutSecs")
	cmd.Flags().Uint16VarP(
		&maxTaskRetentionCountUint, "max-task-retention-count", "M", 0, "MaxTaskRetentionCount",
	)
	cmd.Flags().Uint16VarP(
		&maxTaskRetentionDaysUint, "max-task-retention-days", "R", 0, "MaxTaskRetentionDays",
	)
	cmd.Flags().Uint16VarP(
		&maxConcurrentCpuCoresUint, "max-concurrent-cpu-cores", "C", 0, "MaxConcurrentCpuCores",
	)
	cmd.Flags().StringSliceVarP(
		&containerAccountIdsSlice, "container-account-ids", "u", []string{}, "ContainerAccountIds",
	)
	cmd.Flags().StringSliceVarP(
		&containerIdsSlice, "container-ids", "i", []string{}, "ContainerIds",
	)
	cmd.Flags().StringSliceVarP(
		&exceptContainerAccountIdsSlice, "except-container-account-ids", "U", []string{}, "ExceptContainerAccountIds",
	)
	cmd.Flags().StringSliceVarP(
		&exceptContainerIdsSlice, "except-container-ids", "I", []string{}, "ExceptContainerIds",
	)

	return cmd
}

func (controller *BackupController) UpdateJob() *cobra.Command {
	var jobIdUint, accountIdUint, timeoutSecsUint uint64
	var maxTaskRetentionCountUint, maxTaskRetentionDaysUint, maxConcurrentCpuCoresUint uint16
	var jobStatusBoolStr, jobDescriptionStr, backupScheduleStr string
	var destinationIdsSlice, containerAccountIdsSlice, exceptContainerAccountIdsSlice []string
	var containerIdsSlice, exceptContainerIdsSlice []string

	cmd := &cobra.Command{
		Use:   "update",
		Short: "UpdateBackupJob",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{
				"jobId":     jobIdUint,
				"accountId": accountIdUint,
				"jobStatus": jobStatusBoolStr,
			}

			if jobDescriptionStr != "" {
				requestBody["jobDescription"] = jobDescriptionStr
			}

			if len(destinationIdsSlice) > 0 {
				requestBody["destinationIds"] = sharedHelper.StringSliceValueObjectParser(
					destinationIdsSlice, valueObject.NewBackupDestinationId,
				)
			}

			if backupScheduleStr != "" {
				requestBody["backupSchedule"] = backupScheduleStr
			}

			if timeoutSecsUint != 0 {
				requestBody["timeoutSecs"] = timeoutSecsUint
			}

			if maxTaskRetentionCountUint != 0 {
				requestBody["maxTaskRetentionCount"] = maxTaskRetentionCountUint
			}

			if maxTaskRetentionDaysUint != 0 {
				requestBody["maxTaskRetentionDays"] = maxTaskRetentionDaysUint
			}

			if maxConcurrentCpuCoresUint != 0 {
				requestBody["maxConcurrentCpuCores"] = maxConcurrentCpuCoresUint
			}

			if len(containerAccountIdsSlice) > 0 {
				requestBody["containerAccountIds"] = sharedHelper.StringSliceValueObjectParser(
					containerAccountIdsSlice, valueObject.NewAccountId,
				)
			}

			if len(containerIdsSlice) > 0 {
				requestBody["containerIds"] = sharedHelper.StringSliceValueObjectParser(
					containerIdsSlice, valueObject.NewContainerId,
				)
			}

			if len(exceptContainerAccountIdsSlice) > 0 {
				requestBody["exceptContainerAccountIds"] = sharedHelper.StringSliceValueObjectParser(
					exceptContainerAccountIdsSlice, valueObject.NewAccountId,
				)
			}

			if len(exceptContainerIdsSlice) > 0 {
				requestBody["exceptContainerIds"] = sharedHelper.StringSliceValueObjectParser(
					exceptContainerIdsSlice, valueObject.NewContainerId,
				)
			}

			cliHelper.ServiceResponseWrapper(
				controller.backupService.UpdateJob(requestBody),
			)
		},
	}

	cmd.Flags().Uint64VarP(&jobIdUint, "job-id", "j", 0, "BackupJobId")
	cmd.MarkFlagRequired("job-id")
	cmd.Flags().Uint64VarP(&accountIdUint, "account-id", "a", 0, "AccountId")
	cmd.MarkFlagRequired("account-id")
	cmd.Flags().StringVarP(&jobStatusBoolStr, "job-status", "S", "true", "BackupJobStatus")
	cmd.Flags().StringVarP(&jobDescriptionStr, "job-description", "d", "", "BackupJobDescription")
	cmd.Flags().StringSliceVarP(
		&destinationIdsSlice, "destination-ids", "D", []string{}, "BackupDestinationIds",
	)
	cmd.Flags().StringVarP(&backupScheduleStr, "backup-schedule", "s", "", "BackupSchedule")
	cmd.Flags().Uint64VarP(&timeoutSecsUint, "timeout-secs", "t", 0, "TimeoutSecs")
	cmd.Flags().Uint16VarP(
		&maxTaskRetentionCountUint, "max-task-retention-count", "M", 0, "MaxTaskRetentionCount",
	)
	cmd.Flags().Uint16VarP(
		&maxTaskRetentionDaysUint, "max-task-retention-days", "R", 0, "MaxTaskRetentionDays",
	)
	cmd.Flags().Uint16VarP(
		&maxConcurrentCpuCoresUint, "max-concurrent-cpu-cores", "C", 0, "MaxConcurrentCpuCores",
	)
	cmd.Flags().StringSliceVarP(
		&containerAccountIdsSlice, "container-account-ids", "u", []string{}, "ContainerAccountIds",
	)
	cmd.Flags().StringSliceVarP(
		&containerIdsSlice, "container-ids", "i", []string{}, "ContainerIds",
	)
	cmd.Flags().StringSliceVarP(
		&exceptContainerAccountIdsSlice, "except-container-account-ids", "U", []string{}, "ExceptContainerAccountIds",
	)
	cmd.Flags().StringSliceVarP(
		&exceptContainerIdsSlice, "except-container-ids", "I", []string{}, "ExceptContainerIds",
	)

	return cmd
}

func (controller *BackupController) DeleteJob() *cobra.Command {
	var jobIdUint, accountIdUint uint64

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "DeleteBackupJob",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{
				"jobId":     jobIdUint,
				"accountId": accountIdUint,
			}

			cliHelper.ServiceResponseWrapper(
				controller.backupService.DeleteJob(requestBody),
			)
		},
	}

	cmd.Flags().Uint64VarP(&jobIdUint, "job-id", "j", 0, "BackupJobId")
	cmd.MarkFlagRequired("job-id")
	cmd.Flags().Uint64VarP(&accountIdUint, "account-id", "a", 0, "AccountId")
	cmd.MarkFlagRequired("account-id")
	return cmd
}

func (controller *BackupController) RunJob() *cobra.Command {
	var jobIdUint, accountIdUint uint64

	cmd := &cobra.Command{
		Use:   "run",
		Short: "RunBackupJob",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{
				"jobId":     jobIdUint,
				"accountId": accountIdUint,
			}

			cliHelper.ServiceResponseWrapper(
				controller.backupService.RunJob(requestBody),
			)
		},
	}

	cmd.Flags().Uint64VarP(&jobIdUint, "job-id", "j", 0, "BackupJobId")
	cmd.MarkFlagRequired("job-id")
	cmd.Flags().Uint64VarP(&accountIdUint, "account-id", "a", 0, "AccountId")
	cmd.MarkFlagRequired("account-id")
	return cmd
}

func (controller *BackupController) ReadTask() *cobra.Command {
	var taskIdUint, accountIdUint, jobIdUint, destinationIdUint uint64
	var taskStatusStr, retentionStrategyStr, containerIdStr string
	var paginationPageNumberUint32 uint32
	var paginationItemsPerPageUint16 uint16
	var paginationSortByStr, paginationSortDirectionStr string
	var paginationLastSeenIdStr string

	cmd := &cobra.Command{
		Use:   "get",
		Short: "ReadBackupTasks",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{}

			if taskIdUint != 0 {
				requestBody["taskId"] = taskIdUint
			}

			if accountIdUint != 0 {
				requestBody["accountId"] = accountIdUint
			}

			if jobIdUint != 0 {
				requestBody["jobId"] = jobIdUint
			}

			if destinationIdUint != 0 {
				requestBody["destinationId"] = destinationIdUint
			}

			if taskStatusStr != "" {
				requestBody["taskStatus"] = taskStatusStr
			}

			if retentionStrategyStr != "" {
				requestBody["retentionStrategy"] = retentionStrategyStr
			}

			if containerIdStr != "" {
				requestBody["containerId"] = containerIdStr
			}

			if paginationPageNumberUint32 != 0 {
				requestBody["pageNumber"] = paginationPageNumberUint32
			}
			if paginationItemsPerPageUint16 != 0 {
				requestBody["itemsPerPage"] = paginationItemsPerPageUint16
			}
			if paginationSortByStr != "" {
				requestBody["sortBy"] = paginationSortByStr
			}
			if paginationSortDirectionStr != "" {
				requestBody["sortDirection"] = paginationSortDirectionStr
			}
			if paginationLastSeenIdStr != "" {
				requestBody["lastSeenId"] = paginationLastSeenIdStr
			}

			cliHelper.ServiceResponseWrapper(
				controller.backupService.ReadTask(requestBody),
			)
		},
	}

	cmd.Flags().Uint64VarP(&taskIdUint, "task-id", "t", 0, "BackupTaskId")
	cmd.Flags().Uint64VarP(&jobIdUint, "job-id", "j", 0, "BackupJobId")
	cmd.Flags().Uint64VarP(&accountIdUint, "account-id", "a", 0, "BackupAccountId")
	cmd.Flags().Uint64VarP(&destinationIdUint, "destination-id", "d", 0, "BackupDestinationId")
	cmd.Flags().StringVarP(&taskStatusStr, "task-status", "s", "", "BackupTaskStatus")
	cmd.Flags().StringVarP(&retentionStrategyStr, "retention-strategy", "r", "", "RetentionStrategy")
	cmd.Flags().StringVarP(&containerIdStr, "container-id", "c", "", "ContainerId")
	cmd.Flags().Uint32VarP(
		&paginationPageNumberUint32, "page-number", "o", 0, "PageNumber (Pagination)",
	)
	cmd.Flags().Uint16VarP(
		&paginationItemsPerPageUint16, "items-per-page", "I", 0, "ItemsPerPage (Pagination)",
	)
	cmd.Flags().StringVarP(
		&paginationSortByStr, "sort-by", "y", "", "SortBy (Pagination)",
	)
	cmd.Flags().StringVarP(
		&paginationSortDirectionStr, "sort-direction", "x", "", "SortDirection (Pagination)",
	)
	cmd.Flags().StringVarP(
		&paginationLastSeenIdStr, "last-seen-id", "L", "", "LastSeenId (Pagination)",
	)

	return cmd
}

func (controller *BackupController) DeleteTask() *cobra.Command {
	var taskIdUint uint64
	var shouldDiscardFilesBoolStr string

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "DeleteBackupTask",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{
				"taskId":             taskIdUint,
				"shouldDiscardFiles": shouldDiscardFilesBoolStr,
			}

			cliHelper.ServiceResponseWrapper(
				controller.backupService.DeleteTask(requestBody),
			)
		},
	}

	cmd.Flags().Uint64VarP(&taskIdUint, "task-id", "t", 0, "BackupTaskId")
	cmd.MarkFlagRequired("task-id")
	cmd.Flags().StringVarP(
		&shouldDiscardFilesBoolStr, "should-discard-files", "s", "false", "ShouldDiscardFiles",
	)
	return cmd
}
