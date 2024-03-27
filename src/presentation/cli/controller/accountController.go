package cliController

import (
	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/useCase"
	"github.com/speedianet/control/src/domain/valueObject"
	"github.com/speedianet/control/src/infra"
	"github.com/speedianet/control/src/infra/db"
	cliHelper "github.com/speedianet/control/src/presentation/cli/helper"
	"github.com/spf13/cobra"
)

type AccountController struct {
	persistentDbSvc *db.PersistentDatabaseService
}

func NewAccountController(
	persistentDbSvc *db.PersistentDatabaseService,
) *AccountController {
	return &AccountController{persistentDbSvc: persistentDbSvc}
}

func (controller *AccountController) GetAccounts() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "GetAccounts",
		Run: func(cmd *cobra.Command, args []string) {
			accQueryRepo := infra.NewAccQueryRepo(controller.persistentDbSvc)
			accsList, err := useCase.GetAccounts(accQueryRepo)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, accsList)
		},
	}

	return cmd
}

func (controller *AccountController) AddAccount() *cobra.Command {
	var usernameStr string
	var passwordStr string
	var quotaStr string

	cmd := &cobra.Command{
		Use:   "add",
		Short: "AddNewAccount",
		Run: func(cmd *cobra.Command, args []string) {
			username := valueObject.NewUsernamePanic(usernameStr)
			password := valueObject.NewPasswordPanic(passwordStr)

			var quotaPtr *valueObject.AccountQuota
			if quotaStr != "" {
				quota, err := valueObject.NewAccountQuotaFromString(quotaStr)
				if err != nil {
					cliHelper.ResponseWrapper(false, err.Error())
				}
				quotaPtr = &quota
			}

			addAccountDto := dto.NewAddAccount(
				username,
				password,
				quotaPtr,
			)

			accQueryRepo := infra.NewAccQueryRepo(controller.persistentDbSvc)
			accCmdRepo := infra.NewAccCmdRepo(controller.persistentDbSvc)

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
	cmd.Flags().StringVarP(&quotaStr, "quota", "q", "", "AccountQuota (cpu:memory:disk:inodes)")
	return cmd
}

func (controller *AccountController) UpdateAccount() *cobra.Command {
	var accountIdStr string
	var passwordStr string
	shouldUpdateApiKeyBool := false
	var quotaStr string

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

			var quotaPtr *valueObject.AccountQuota
			if quotaStr != "" {
				quota, err := valueObject.NewAccountQuotaFromString(quotaStr)
				if err != nil {
					cliHelper.ResponseWrapper(false, err.Error())
				}
				quotaPtr = &quota
			}

			updateAccountDto := dto.NewUpdateAccount(
				accountId,
				passPtr,
				shouldUpdateApiKeyPtr,
				quotaPtr,
			)

			accQueryRepo := infra.NewAccQueryRepo(controller.persistentDbSvc)
			accCmdRepo := infra.NewAccCmdRepo(controller.persistentDbSvc)

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

			err := useCase.UpdateAccount(
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
	cmd.Flags().StringVarP(&quotaStr, "quota", "q", "", "AccountQuota (cpu:memory:disk:inodes)")
	return cmd
}

func (controller *AccountController) DeleteAccount() *cobra.Command {
	var accountIdStr string

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "DeleteAccount",
		Run: func(cmd *cobra.Command, args []string) {
			accountId := valueObject.NewAccountIdPanic(accountIdStr)

			accQueryRepo := infra.NewAccQueryRepo(controller.persistentDbSvc)
			accCmdRepo := infra.NewAccCmdRepo(controller.persistentDbSvc)
			containerQueryRepo := infra.NewContainerQueryRepo(controller.persistentDbSvc)

			err := useCase.DeleteAccount(
				accQueryRepo,
				accCmdRepo,
				accountId,
				containerQueryRepo,
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
