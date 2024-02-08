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

	containerController := cliController.NewContainerController(router.dbSvc)
	containerCmd.AddCommand(containerController.GetContainers())
	containerCmd.AddCommand(containerController.AddContainer())
	containerCmd.AddCommand(containerController.UpdateContainer())
	containerCmd.AddCommand(containerController.DeleteContainer())

	var containerProfileCmd = &cobra.Command{
		Use:   "profile",
		Short: "ContainerProfileManagement",
	}

	containerProfileController := cliController.NewContainerProfileController(router.dbSvc)
	containerProfileCmd.AddCommand(containerProfileController.GetContainerProfiles())
	containerProfileCmd.AddCommand(containerProfileController.AddContainerProfile())
	containerProfileCmd.AddCommand(containerProfileController.UpdateContainerProfile())
	containerProfileCmd.AddCommand(containerProfileController.DeleteContainerProfile())

	var containerRegistryCmd = &cobra.Command{
		Use:   "registry",
		Short: "ContainerRegistryManagement",
	}

	containerRegistryController := cliController.NewContainerRegistryController(router.dbSvc)
	containerRegistryCmd.AddCommand(containerRegistryController.GetRegistryImages())
	containerRegistryCmd.AddCommand(containerRegistryController.GetRegistryTaggedImage())

	containerCmd.AddCommand(containerProfileCmd)
	containerCmd.AddCommand(containerRegistryCmd)
	rootCmd.AddCommand(containerCmd)
}

func (router CliRouter) licenseRoutes() {
	var licenseCmd = &cobra.Command{
		Use:   "license",
		Short: "LicenseManagement",
	}

	licenseController := cliController.NewLicenseController(router.dbSvc)
	licenseCmd.AddCommand(licenseController.GetLicenseInfo())
	rootCmd.AddCommand(licenseCmd)
}

func (router CliRouter) mappingRoutes() {
	var mappingCmd = &cobra.Command{
		Use:   "mapping",
		Short: "MappingManagement",
	}

	mappingController := cliController.NewMappingController(router.dbSvc)
	mappingCmd.AddCommand(mappingController.GetMappings())
	mappingCmd.AddCommand(mappingController.AddMapping())
	mappingCmd.AddCommand(mappingController.DeleteMapping())

	var mappingTargetCmd = &cobra.Command{
		Use:   "target",
		Short: "MappingTargetManagement",
	}

	mappingTargetCmd.AddCommand(mappingController.AddMappingTarget())
	mappingTargetCmd.AddCommand(mappingController.DeleteMappingTarget())

	mappingCmd.AddCommand(mappingTargetCmd)
	rootCmd.AddCommand(mappingCmd)
}

func (router CliRouter) o11yRoutes() {
	var o11yCmd = &cobra.Command{
		Use:   "o11y",
		Short: "O11yManagement",
	}

	o11yController := cliController.O11yController{}
	o11yCmd.AddCommand(o11yController.GetO11yOverview())
	rootCmd.AddCommand(o11yCmd)
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
