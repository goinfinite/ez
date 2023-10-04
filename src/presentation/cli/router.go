package cli

import (
	"fmt"

	api "github.com/speedianet/sfm/src/presentation/api"
	cliController "github.com/speedianet/sfm/src/presentation/cli/controller"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print software version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Speedia FleetManager v0.0.1")
	},
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the SFM server (default to port 10001)",
	Run: func(cmd *cobra.Command, args []string) {
		api.ApiInit()
	},
}

func accountRoutes() {
	var accountCmd = &cobra.Command{
		Use:   "account",
		Short: "AccountManagement",
	}

	rootCmd.AddCommand(accountCmd)
	accountCmd.AddCommand(cliController.GetAccountsController())
	accountCmd.AddCommand(cliController.AddAccountController())
	accountCmd.AddCommand(cliController.UpdateAccountController())
	accountCmd.AddCommand(cliController.DeleteAccountController())
}

func containerRoutes() {
	var containerCmd = &cobra.Command{
		Use:   "container",
		Short: "ContainerManagement",
	}

	rootCmd.AddCommand(containerCmd)
	containerCmd.AddCommand(cliController.GetContainersController())
	containerCmd.AddCommand(cliController.AddContainerController())
	containerCmd.AddCommand(cliController.UpdateContainerController())
	containerCmd.AddCommand(cliController.DeleteContainerController())

	var resourceProfileCmd = &cobra.Command{
		Use:   "resource-profile",
		Short: "ResourceProfileManagement",
	}

	containerCmd.AddCommand(resourceProfileCmd)
	resourceProfileCmd.AddCommand(cliController.GetResourceProfilesController())
	resourceProfileCmd.AddCommand(cliController.AddResourceProfileController())
	resourceProfileCmd.AddCommand(cliController.UpdateResourceProfileController())
	resourceProfileCmd.AddCommand(cliController.DeleteResourceProfileController())
}

func o11yRoutes() {
	var o11yCmd = &cobra.Command{
		Use:   "o11y",
		Short: "O11yManagement",
	}

	rootCmd.AddCommand(o11yCmd)
	o11yCmd.AddCommand(cliController.GetO11yOverviewController())
}

func registerCliRoutes() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(cliController.SysInstallController())
	accountRoutes()
	containerRoutes()
	o11yRoutes()
}
