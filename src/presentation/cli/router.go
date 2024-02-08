package cli

import (
	"fmt"

	"github.com/speedianet/control/src/infra/db"
	api "github.com/speedianet/control/src/presentation/api"
	cliController "github.com/speedianet/control/src/presentation/cli/controller"
	"github.com/spf13/cobra"
)

type CliRouter struct {
	dbSvc *db.DatabaseService
}

func NewCliRouter(dbSvc *db.DatabaseService) CliRouter {
	return CliRouter{dbSvc: dbSvc}
}

func (router CliRouter) accountRoutes() {
	var accountCmd = &cobra.Command{
		Use:   "account",
		Short: "AccountManagement",
	}

	accountController := cliController.NewAccountController(router.dbSvc)
	accountCmd.AddCommand(accountController.GetAccounts())
	accountCmd.AddCommand(accountController.AddAccount())
	accountCmd.AddCommand(accountController.UpdateAccount())
	accountCmd.AddCommand(accountController.DeleteAccount())
	rootCmd.AddCommand(accountCmd)
}

func (router CliRouter) containerRoutes() {
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

	var containerRegistryCmd = &cobra.Command{
		Use:   "registry",
		Short: "ContainerRegistryManagement",
	}

	containerCmd.AddCommand(containerRegistryCmd)
	containerRegistryCmd.AddCommand(cliController.GetRegistryImagesController())
	containerRegistryCmd.AddCommand(cliController.GetRegistryTaggedImageController())
}

func (router CliRouter) licenseRoutes() {
	var licenseCmd = &cobra.Command{
		Use:   "license",
		Short: "LicenseManagement",
	}

	rootCmd.AddCommand(licenseCmd)
	licenseCmd.AddCommand(cliController.GetLicenseInfoController())
}

func (router CliRouter) mappingRoutes() {
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
	mappingTargetCmd.AddCommand(cliController.DeleteMappingTargetController())
}

func (router CliRouter) o11yRoutes() {
	var o11yCmd = &cobra.Command{
		Use:   "o11y",
		Short: "O11yManagement",
	}

	rootCmd.AddCommand(o11yCmd)
	o11yCmd.AddCommand(cliController.GetO11yOverviewController())
}

func (router CliRouter) systemRoutes() {
	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print software version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Speedia Control v0.0.1")
		},
	}

	var serveCmd = &cobra.Command{
		Use:   "serve",
		Short: "Start the CONTROL server (default to port 3141)",
		Run: func(cmd *cobra.Command, args []string) {
			api.ApiInit()
		},
	}

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(cliController.SysInstallController())
}

func (router CliRouter) RegisterRoutes() {
	router.accountRoutes()
	router.containerRoutes()
	router.licenseRoutes()
	router.mappingRoutes()
	router.o11yRoutes()
	router.systemRoutes()
}
