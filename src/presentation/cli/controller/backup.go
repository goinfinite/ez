package cliController

import (
	"github.com/goinfinite/ez/src/infra/db"
	cliHelper "github.com/goinfinite/ez/src/presentation/cli/helper"
	"github.com/goinfinite/ez/src/presentation/service"
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
	cmd.Flags().StringVarP(&remoteHostTypeStr, "remoteHostType", "r", "", "RemoteHostType")
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
