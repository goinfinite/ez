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

func o11yRoutes() {
	var o11yCmd = &cobra.Command{
		Use:   "o11y",
		Short: "O11yManagement",
	}

	rootCmd.AddCommand(o11yCmd)
	o11yCmd.AddCommand(cliController.GetO11yOverviewController())
}

func accountRoutes() {
	var accountCmd = &cobra.Command{
		Use:   "account",
		Short: "AccountManagement",
	}

	rootCmd.AddCommand(accountCmd)
	accountCmd.AddCommand(cliController.GetAccountsController())
	accountCmd.AddCommand(cliController.AddAccountController())
	accountCmd.AddCommand(cliController.DeleteAccountController())
	accountCmd.AddCommand(cliController.UpdateAccountController())
}

func containerRoutes() {
	var containerCmd = &cobra.Command{
		Use:   "container",
		Short: "ContainerManagement",
	}

	rootCmd.AddCommand(containerCmd)
	containerCmd.AddCommand(cliController.GetContainersController())
	containerCmd.AddCommand(cliController.AddContainerController())
	containerCmd.AddCommand(cliController.DeleteContainerController())
	containerCmd.AddCommand(cliController.UpdateContainerController())
}

func registerCliRoutes() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(cliController.SysInstallController())
	accountRoutes()
	containerRoutes()
	o11yRoutes()
}
