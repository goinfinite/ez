package cli

import (
	"fmt"

	"github.com/speedianet/control/src/infra/db"
	infraEnvs "github.com/speedianet/control/src/infra/envs"
	"github.com/speedianet/control/src/presentation"
	cliController "github.com/speedianet/control/src/presentation/cli/controller"
	"github.com/spf13/cobra"
)

type Router struct {
	persistentDbSvc *db.PersistentDatabaseService
	transientDbSvc  *db.TransientDatabaseService
	trailDbSvc      *db.TrailDatabaseService
}

func NewRouter(
	persistentDbSvc *db.PersistentDatabaseService,
	transientDbSvc *db.TransientDatabaseService,
	trailDbSvc *db.TrailDatabaseService,
) Router {
	return Router{
		persistentDbSvc: persistentDbSvc,
		transientDbSvc:  transientDbSvc,
		trailDbSvc:      trailDbSvc,
	}
}

func (router *Router) accountRoutes() {
	var accountCmd = &cobra.Command{
		Use:   "account",
		Short: "AccountManagement",
	}

	accountController := cliController.NewAccountController(
		router.persistentDbSvc, router.trailDbSvc,
	)
	accountCmd.AddCommand(accountController.Read())
	accountCmd.AddCommand(accountController.Create())
	accountCmd.AddCommand(accountController.Update())
	accountCmd.AddCommand(accountController.Delete())
	rootCmd.AddCommand(accountCmd)
}

func (router *Router) containerRoutes() {
	var containerCmd = &cobra.Command{
		Use:   "container",
		Short: "ContainerManagement",
	}

	containerController := cliController.NewContainerController(router.persistentDbSvc)
	containerCmd.AddCommand(containerController.Read())
	containerCmd.AddCommand(containerController.ReadWithMetrics())
	containerCmd.AddCommand(containerController.Create())
	containerCmd.AddCommand(containerController.Update())
	containerCmd.AddCommand(containerController.Delete())

	var containerProfileCmd = &cobra.Command{
		Use:   "profile",
		Short: "ContainerProfileManagement",
	}

	containerProfileController := cliController.NewContainerProfileController(
		router.persistentDbSvc,
	)
	containerProfileCmd.AddCommand(containerProfileController.Read())
	containerProfileCmd.AddCommand(containerProfileController.Create())
	containerProfileCmd.AddCommand(containerProfileController.Update())
	containerProfileCmd.AddCommand(containerProfileController.Delete())

	var containerRegistryCmd = &cobra.Command{
		Use:   "registry",
		Short: "ContainerRegistryManagement",
	}

	containerRegistryController := cliController.NewContainerRegistryController(router.persistentDbSvc)
	containerRegistryCmd.AddCommand(containerRegistryController.ReadRegistryImages())
	containerRegistryCmd.AddCommand(containerRegistryController.ReadRegistryTaggedImage())

	var containerImageCmd = &cobra.Command{
		Use:   "image",
		Short: "ContainerImageManagement",
	}

	containerImageController := cliController.NewContainerImageController(
		router.persistentDbSvc, router.trailDbSvc,
	)
	containerImageCmd.AddCommand(containerImageController.Read())
	containerImageCmd.AddCommand(containerImageController.Delete())
	containerImageCmd.AddCommand(containerImageController.CreateSnapshot())

	var containerImageArchiveCmd = &cobra.Command{
		Use:   "archive",
		Short: "ContainerImageArchiveManagement",
	}
	containerImageArchiveCmd.AddCommand(containerImageController.ReadArchiveFiles())
	containerImageArchiveCmd.AddCommand(containerImageController.CreateArchiveFile())
	containerImageArchiveCmd.AddCommand(containerImageController.DeleteArchiveFile())
	containerImageCmd.AddCommand(containerImageArchiveCmd)

	containerCmd.AddCommand(containerProfileCmd)
	containerCmd.AddCommand(containerRegistryCmd)
	containerCmd.AddCommand(containerImageCmd)
	rootCmd.AddCommand(containerCmd)
}

func (router *Router) licenseRoutes() {
	var licenseCmd = &cobra.Command{
		Use:   "license",
		Short: "LicenseManagement",
	}

	licenseController := cliController.NewLicenseController(
		router.persistentDbSvc,
		router.transientDbSvc,
	)
	licenseCmd.AddCommand(licenseController.ReadLicenseInfo())
	licenseCmd.AddCommand(licenseController.RefreshLicense())
	rootCmd.AddCommand(licenseCmd)
}

func (router *Router) mappingRoutes() {
	var mappingCmd = &cobra.Command{
		Use:   "mapping",
		Short: "MappingManagement",
	}

	mappingController := cliController.NewMappingController(router.persistentDbSvc)
	mappingCmd.AddCommand(mappingController.Read())
	mappingCmd.AddCommand(mappingController.Create())
	mappingCmd.AddCommand(mappingController.Delete())

	var mappingTargetCmd = &cobra.Command{
		Use:   "target",
		Short: "MappingTargetManagement",
	}

	mappingTargetCmd.AddCommand(mappingController.CreateTarget())
	mappingTargetCmd.AddCommand(mappingController.DeleteTarget())

	mappingCmd.AddCommand(mappingTargetCmd)
	rootCmd.AddCommand(mappingCmd)
}

func (router *Router) o11yRoutes() {
	var o11yCmd = &cobra.Command{
		Use:   "o11y",
		Short: "O11yManagement",
	}

	o11yController := cliController.O11yController{}
	o11yCmd.AddCommand(o11yController.ReadO11yOverview())
	rootCmd.AddCommand(o11yCmd)
}

func (router *Router) scheduledTaskRoutes() {
	var scheduledTaskCmd = &cobra.Command{
		Use:   "task",
		Short: "ScheduledTaskManagement",
	}

	scheduledTaskController := cliController.NewScheduledTaskController(router.persistentDbSvc)
	scheduledTaskCmd.AddCommand(scheduledTaskController.Read())
	scheduledTaskCmd.AddCommand(scheduledTaskController.Update())
	rootCmd.AddCommand(scheduledTaskCmd)
}

func (router *Router) systemRoutes() {
	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "PrintVersion",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Speedia Control " + infraEnvs.SpeediaControlVersion)
		},
	}

	var serveCmd = &cobra.Command{
		Use:   "serve",
		Short: "ServeApiDashboard",
		Run: func(cmd *cobra.Command, args []string) {
			presentation.HttpServerInit(
				router.persistentDbSvc, router.transientDbSvc, router.trailDbSvc,
			)
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
	router.scheduledTaskRoutes()
	router.systemRoutes()
}
