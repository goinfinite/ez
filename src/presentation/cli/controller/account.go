package cliController

import (
	"github.com/speedianet/control/src/domain/valueObject"
	"github.com/speedianet/control/src/infra/db"
	cliHelper "github.com/speedianet/control/src/presentation/cli/helper"
	"github.com/speedianet/control/src/presentation/service"
	"github.com/spf13/cobra"
)

type AccountController struct {
	accountService *service.AccountService
}

func NewAccountController(
	persistentDbSvc *db.PersistentDatabaseService,
	trailDbSvc *db.TrailDatabaseService,
) *AccountController {
	return &AccountController{
		accountService: service.NewAccountService(persistentDbSvc, trailDbSvc),
	}
}

func (controller *AccountController) Read() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "ReadAccounts",
		Run: func(cmd *cobra.Command, args []string) {
			cliHelper.ServiceResponseWrapper(controller.accountService.Read())
		},
	}

	return cmd
}

func (controller *AccountController) Create() *cobra.Command {
	var usernameStr string
	var passwordStr string
	var quotaStr string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "CreateAccount",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{
				"username":  usernameStr,
				"password":  passwordStr,
				"ipAddress": valueObject.NewLocalhostIpAddress().String(),
			}

			if quotaStr != "" {
				quota, err := valueObject.NewAccountQuotaFromString(quotaStr)
				if err != nil {
					cliHelper.ResponseWrapper(false, err.Error())
				}
				requestBody["quota"] = quota
			}

			cliHelper.ServiceResponseWrapper(
				controller.accountService.Create(requestBody),
			)
		},
	}

	cmd.Flags().StringVarP(&usernameStr, "username", "u", "", "Username")
	cmd.MarkFlagRequired("username")
	cmd.Flags().StringVarP(&passwordStr, "password", "p", "", "Password")
	cmd.MarkFlagRequired("password")
	cmd.Flags().StringVarP(&quotaStr, "quota", "q", "", "AccountQuota (cpu:memory:disk:inodes)")
	return cmd
}

func (controller *AccountController) Update() *cobra.Command {
	var accountIdStr string
	var passwordStr string
	shouldUpdateApiKeyBool := false
	var quotaStr string

	cmd := &cobra.Command{
		Use:   "update",
		Short: "UpdateAccount (pass or apiKey)",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{
				"accountId": accountIdStr,
				"ipAddress": valueObject.NewLocalhostIpAddress().String(),
			}

			if passwordStr != "" {
				requestBody["password"] = passwordStr
			}

			if shouldUpdateApiKeyBool {
				requestBody["shouldUpdateApiKey"] = true
			}

			if quotaStr != "" {
				quota, err := valueObject.NewAccountQuotaFromString(quotaStr)
				if err != nil {
					cliHelper.ResponseWrapper(false, err.Error())
				}
				requestBody["quota"] = quota
			}

			cliHelper.ServiceResponseWrapper(
				controller.accountService.Update(requestBody),
			)
		},
	}

	cmd.Flags().StringVarP(&accountIdStr, "id", "i", "", "AccountId")
	cmd.MarkFlagRequired("id")
	cmd.Flags().StringVarP(&passwordStr, "password", "p", "", "Password")
	cmd.Flags().BoolVarP(
		&shouldUpdateApiKeyBool, "update-api-key", "k", false, "ShouldUpdateApiKey",
	)
	cmd.Flags().StringVarP(&quotaStr, "quota", "q", "", "AccountQuota (cpu:memory:disk:inodes)")
	return cmd
}

func (controller *AccountController) Delete() *cobra.Command {
	var accountIdStr string

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "DeleteAccount",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{
				"accountId": accountIdStr,
				"ipAddress": valueObject.NewLocalhostIpAddress().String(),
			}

			cliHelper.ServiceResponseWrapper(
				controller.accountService.Delete(requestBody),
			)
		},
	}

	cmd.Flags().StringVarP(&accountIdStr, "id", "i", "", "AccountId")
	cmd.MarkFlagRequired("id")
	return cmd
}
