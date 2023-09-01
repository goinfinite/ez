package cliController

import (
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

			if cpuCores == 0 {
				cpuCores = 1
			}
			if memoryBytesUint == 0 {
				memoryBytesUint = 1073741824
			}
			memoryBytes := valueObject.Byte(memoryBytesUint)
			if diskBytesUint == 0 {
				diskBytesUint = 5368709120
			}
			diskBytes := valueObject.Byte(diskBytesUint)
			if inodes == 0 {
				inodes = 500000
			}
			quota := valueObject.NewAccountQuota(
				cpuCores,
				memoryBytes,
				diskBytes,
				inodes,
			)

			addAccountDto := dto.NewAddAccount(
				username,
				password,
				quota,
			)

			accQueryRepo := infra.AccQueryRepo{}
			accCmdRepo := infra.AccCmdRepo{}

			err := useCase.AddAccount(
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
	cmd.Flags().Float64VarP(&cpuCores, "cpu-cores", "c", 0, "CpuCores")
	cmd.Flags().Uint64VarP(&memoryBytesUint, "memory-bytes", "m", 0, "MemoryBytes")
	cmd.Flags().Uint64VarP(&diskBytesUint, "disk-bytes", "d", 0, "DiskBytes")
	cmd.Flags().Uint64VarP(&inodes, "inodes", "i", 0, "Inodes")
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

	cmd.Flags().StringVarP(&accountIdStr, "account-id", "i", "", "AccountId")
	cmd.MarkFlagRequired("account-id")
	return cmd
}

func UpdateAccountController() *cobra.Command {
	var accountIdStr string
	var passwordStr string
	shouldUpdateApiKeyBool := false

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

			updateAccountDto := dto.NewUpdateAccount(
				accountId,
				passPtr,
				shouldUpdateApiKeyPtr,
			)

			accQueryRepo := infra.AccQueryRepo{}
			accCmdRepo := infra.AccCmdRepo{}

			if updateAccountDto.Password != nil {
				useCase.UpdateAccountPassword(
					accQueryRepo,
					accCmdRepo,
					updateAccountDto,
				)
			}

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
		},
	}

	cmd.Flags().StringVarP(&accountIdStr, "account-id", "i", "", "AccountId")
	cmd.MarkFlagRequired("account-id")
	cmd.Flags().StringVarP(&passwordStr, "password", "p", "", "Password")
	cmd.Flags().BoolVarP(
		&shouldUpdateApiKeyBool,
		"update-api-key",
		"k",
		false,
		"ShouldUpdateApiKey",
	)
	return cmd
}
