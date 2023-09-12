package cliController

import (
	"log"
	"time"

	"github.com/speedianet/sfm/src/domain/dto"
	"github.com/speedianet/sfm/src/domain/useCase"
	"github.com/speedianet/sfm/src/domain/valueObject"
	"github.com/speedianet/sfm/src/infra"
	cliHelper "github.com/speedianet/sfm/src/presentation/cli/helper"
	"github.com/spf13/cobra"
)

func GetAccountsController() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "GetAccounts",
		Run: func(cmd *cobra.Command, args []string) {
			accQueryRepo := infra.AccQueryRepo{}
			accsList, err := useCase.GetAccounts(accQueryRepo)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, accsList)
		},
	}

	return cmd
}

func accQuotaFactory(
	cpuCores float64,
	memoryBytesUint uint64,
	diskBytesUint uint64,
	inodes uint64,
	withDefaults bool,
) (valueObject.AccountQuota, error) {
	accQuotaDefaults := valueObject.NewAccountQuotaWithDefaultValues()
	if !withDefaults {
		accQuotaDefaults = valueObject.NewAccountQuotaWithBlankValues()
	}

	cpuCoresCount := accQuotaDefaults.CpuCores
	if cpuCores != 0 {
		cpuCoresCount = valueObject.NewCpuCoresCountPanic(cpuCores)
	}

	memoryBytes := accQuotaDefaults.MemoryBytes
	if memoryBytesUint != 0 {
		memoryBytes = valueObject.NewBytePanic(memoryBytesUint)
	}

	diskBytes := accQuotaDefaults.DiskBytes
	if diskBytesUint != 0 {
		diskBytes = valueObject.NewBytePanic(diskBytesUint)
	}

	inodesCount := accQuotaDefaults.Inodes
	if inodes != 0 {
		inodesCount = valueObject.NewInodesCountPanic(inodes)
	}

	return valueObject.NewAccountQuota(
		cpuCoresCount,
		memoryBytes,
		diskBytes,
		inodesCount,
	), nil
}

func AddAccountController() *cobra.Command {
	var usernameStr string
	var passwordStr string
	var cpuCores float64
	var memoryBytesUint uint64
	var diskBytesUint uint64
	var inodes uint64

	cmd := &cobra.Command{
		Use:   "add",
		Short: "AddNewAccount",
		Run: func(cmd *cobra.Command, args []string) {
			username := valueObject.NewUsernamePanic(usernameStr)
			password := valueObject.NewPasswordPanic(passwordStr)

			quota, err := accQuotaFactory(
				cpuCores,
				memoryBytesUint,
				diskBytesUint,
				inodes,
				true,
			)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			addAccountDto := dto.NewAddAccount(
				username,
				password,
				quota,
			)

			accQueryRepo := infra.AccQueryRepo{}
			accCmdRepo := infra.AccCmdRepo{}

			err = useCase.AddAccount(
				accQueryRepo,
				accCmdRepo,
				addAccountDto,
			)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, "AccountAdded")
		},
	}

	cmd.Flags().StringVarP(&usernameStr, "username", "u", "", "Username")
	cmd.MarkFlagRequired("username")
	cmd.Flags().StringVarP(&passwordStr, "password", "p", "", "Password")
	cmd.MarkFlagRequired("password")
	cmd.Flags().Float64VarP(&cpuCores, "cpu", "c", 0, "CpuCores")
	cmd.Flags().Uint64VarP(&memoryBytesUint, "memory", "m", 0, "MemoryInBytes")
	cmd.Flags().Uint64VarP(&diskBytesUint, "disk", "d", 0, "DiskInBytes")
	cmd.Flags().Uint64VarP(&inodes, "inodes", "n", 0, "Inodes")
	return cmd
}

func DeleteAccountController() *cobra.Command {
	var accountIdStr string

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "DeleteAccount",
		Run: func(cmd *cobra.Command, args []string) {
			accountId := valueObject.NewAccountIdPanic(accountIdStr)

			accQueryRepo := infra.AccQueryRepo{}
			accCmdRepo := infra.AccCmdRepo{}

			err := useCase.DeleteAccount(
				accQueryRepo,
				accCmdRepo,
				accountId,
			)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, "AccountDeleted")
		},
	}

	cmd.Flags().StringVarP(&accountIdStr, "id", "i", "", "AccountId")
	cmd.MarkFlagRequired("id")
	return cmd
}

func UpdateAccountController() *cobra.Command {
	var accountIdStr string
	var passwordStr string
	shouldUpdateApiKeyBool := false
	var cpuCores float64
	var memoryBytesUint uint64
	var diskBytesUint uint64
	var inodes uint64

	cmd := &cobra.Command{
		Use:   "update",
		Short: "UpdateAccount (pass or apiKey)",
		Run: func(cmd *cobra.Command, args []string) {
			accountId := valueObject.NewAccountIdPanic(accountIdStr)

			var passPtr *valueObject.Password
			if passwordStr != "" {
				password := valueObject.NewPasswordPanic(passwordStr)
				passPtr = &password
			}

			var shouldUpdateApiKeyPtr *bool
			if shouldUpdateApiKeyBool {
				shouldUpdateApiKeyPtr = &shouldUpdateApiKeyBool
			}

			quota, err := accQuotaFactory(
				cpuCores,
				memoryBytesUint,
				diskBytesUint,
				inodes,
				false,
			)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			updateAccountDto := dto.NewUpdateAccount(
				accountId,
				passPtr,
				shouldUpdateApiKeyPtr,
				&quota,
			)

			accQueryRepo := infra.AccQueryRepo{}
			accCmdRepo := infra.AccCmdRepo{}

			if shouldUpdateApiKeyBool {
				newKey, err := useCase.UpdateAccountApiKey(
					accQueryRepo,
					accCmdRepo,
					updateAccountDto,
				)
				if err != nil {
					cliHelper.ResponseWrapper(false, err.Error())
				}

				cliHelper.ResponseWrapper(true, newKey)
			}

			err = useCase.UpdateAccount(
				accQueryRepo,
				accCmdRepo,
				updateAccountDto,
			)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}
		},
	}

	cmd.Flags().StringVarP(&accountIdStr, "id", "i", "", "AccountId")
	cmd.MarkFlagRequired("id")
	cmd.Flags().StringVarP(&passwordStr, "password", "p", "", "Password")
	cmd.Flags().BoolVarP(
		&shouldUpdateApiKeyBool,
		"update-api-key",
		"k",
		false,
		"ShouldUpdateApiKey",
	)
	cmd.Flags().Float64VarP(&cpuCores, "cpu", "c", 0, "CpuCores")
	cmd.Flags().Uint64VarP(&memoryBytesUint, "memory", "m", 0, "MemoryInBytes")
	cmd.Flags().Uint64VarP(&diskBytesUint, "disk", "d", 0, "DiskInBytes")
	cmd.Flags().Uint64VarP(&inodes, "inodes", "n", 0, "Inodes")
	return cmd
}

func UpdateAccountsQuotaUsageController() {
	taskInterval := time.Duration(15) * time.Minute
	timer := time.NewTicker(taskInterval)
	defer timer.Stop()

	accCmdRepo := infra.AccCmdRepo{}
	for range timer.C {
		err := accCmdRepo.UpdateQuotasUsage()
		if err != nil {
			log.Printf("UpdateQuotasUsageError: %v", err)
		}
	}
}
