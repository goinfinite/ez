package cli

import (
	"fmt"

	"github.com/speedianet/control/src/infra/db"
	api "github.com/speedianet/control/src/presentation/api"
	cliController "github.com/speedianet/control/src/presentation/cli/controller"
	"github.com/spf13/cobra"
)

type Router struct {
	persistDbSvc *db.PersistentDatabaseService
}

func NewRouter(persistDbSvc *db.PersistentDatabaseService) Router {
	return Router{persistDbSvc: persistDbSvc}
}

func (router Router) accountRoutes() {
	var accountCmd = &cobra.Command{
		Use:   "account",
		Short: "AccountManagement",
	}

	accountController := cliController.NewAccountController(router.persistDbSvc)
	accountCmd.AddCommand(accountController.GetAccounts())
	accountCmd.AddCommand(accountController.AddAccount())
	accountCmd.AddCommand(accountController.UpdateAccount())
	accountCmd.AddCommand(accountController.DeleteAccount())
	rootCmd.AddCommand(accountCmd)
}

func (router Router) containerRoutes() {
	var containerCmd = &cobra.Command{
		Use:   "container",
		Short: "ContainerManagement",
	}

	containerController := cliController.NewContainerController(router.persistDbSvc)
	containerCmd.AddCommand(containerController.GetContainers())
	containerCmd.AddCommand(containerController.AddContainer())
	containerCmd.AddCommand(containerController.UpdateContainer())
	containerCmd.AddCommand(containerController.DeleteContainer())

	var containerProfileCmd = &cobra.Command{
		Use:   "profile",
		Short: "ContainerProfileManagement",
	}

	containerProfileController := cliController.NewContainerProfileController(router.persistDbSvc)
	containerProfileCmd.AddCommand(containerProfileController.GetContainerProfiles())
	containerProfileCmd.AddCommand(containerProfileController.AddContainerProfile())
	containerProfileCmd.AddCommand(containerProfileController.UpdateContainerProfile())
	containerProfileCmd.AddCommand(containerProfileController.DeleteContainerProfile())

	var containerRegistryCmd = &cobra.Command{
		Use:   "registry",
		Short: "ContainerRegistryManagement",
	}

	containerRegistryController := cliController.NewContainerRegistryController(router.persistDbSvc)
	containerRegistryCmd.AddCommand(containerRegistryController.GetRegistryImages())
	containerRegistryCmd.AddCommand(containerRegistryController.GetRegistryTaggedImage())

	containerCmd.AddCommand(containerProfileCmd)
	containerCmd.AddCommand(containerRegistryCmd)
	rootCmd.AddCommand(containerCmd)
}

func (router Router) licenseRoutes() {
	var licenseCmd = &cobra.Command{
		Use:   "license",
		Short: "LicenseManagement",
	}

	licenseController := cliController.NewLicenseController(router.persistDbSvc)
	licenseCmd.AddCommand(licenseController.GetLicenseInfo())
	rootCmd.AddCommand(licenseCmd)
}

func (router Router) mappingRoutes() {
	var mappingCmd = &cobra.Command{
		Use:   "mapping",
		Short: "MappingManagement",
	}

	mappingController := cliController.NewMappingController(router.persistDbSvc)
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

func (router Router) o11yRoutes() {
	var o11yCmd = &cobra.Command{
		Use:   "o11y",
		Short: "O11yManagement",
	}

	o11yController := cliController.O11yController{}
	o11yCmd.AddCommand(o11yController.GetO11yOverview())
	rootCmd.AddCommand(o11yCmd)
}

func (router Router) systemRoutes() {
	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "PrintVersion",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Speedia Control v0.0.1")
		},
	}

	var serveCmd = &cobra.Command{
		Use:   "serve",
		Short: "ServeApiDashboard",
		Run: func(cmd *cobra.Command, args []string) {
			api.ApiInit(router.persistDbSvc)
		},
	}

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(serveCmd)
}

func (router Router) RegisterRoutes() {
	router.accountRoutes()
	router.containerRoutes()
	router.licenseRoutes()
	router.mappingRoutes()
	router.o11yRoutes()
	router.systemRoutes()
}
