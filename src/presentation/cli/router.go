package cli

import (
	"fmt"

	api "github.com/speedianet/control/src/presentation/api"
	cliController "github.com/speedianet/control/src/presentation/cli/controller"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print software version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Speedia Control v0.0.1")
	},
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the CONTROL server (default to port 10001)",
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

	var containerProfileCmd = &cobra.Command{
		Use:   "profile",
		Short: "ContainerProfileManagement",
	}

	containerCmd.AddCommand(containerProfileCmd)
	containerProfileCmd.AddCommand(cliController.GetContainerProfilesController())
	containerProfileCmd.AddCommand(cliController.AddContainerProfileController())
	containerProfileCmd.AddCommand(cliController.UpdateContainerProfileController())
	containerProfileCmd.AddCommand(cliController.DeleteContainerProfileController())
}

func mappingRoutes() {
	var mappingCmd = &cobra.Command{
		Use:   "mapping",
		Short: "MappingManagement",
	}

	rootCmd.AddCommand(mappingCmd)
	mappingCmd.AddCommand(cliController.GetMappingsController())
	mappingCmd.AddCommand(cliController.AddMappingController())
	mappingCmd.AddCommand(cliController.DeleteMappingController())

	var mappingTargetCmd = &cobra.Command{
		Use:   "target",
		Short: "MappingTargetManagement",
	}

	mappingCmd.AddCommand(mappingTargetCmd)
	mappingTargetCmd.AddCommand(cliController.AddMappingTargetController())
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
	mappingRoutes()
	o11yRoutes()
}
