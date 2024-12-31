package cli

import (
	"fmt"

	"github.com/goinfinite/ez/src/infra/db"
	infraEnvs "github.com/goinfinite/ez/src/infra/envs"
	"github.com/goinfinite/ez/src/presentation"
	cliController "github.com/goinfinite/ez/src/presentation/cli/controller"
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
	accountCmd.AddCommand(accountController.RefreshQuotas())
	accountCmd.AddCommand(accountController.Delete())
	rootCmd.AddCommand(accountCmd)
}

func (router *Router) backupRoutes() {
	var backupCmd = &cobra.Command{
		Use:   "backup",
		Short: "BackupManagement",
	}

	backupController := cliController.NewBackupController(
		router.persistentDbSvc, router.trailDbSvc,
	)

	var backupDestinationCmd = &cobra.Command{
		Use:   "destination",
		Short: "BackupDestinationManagement",
	}
	backupDestinationCmd.AddCommand(backupController.ReadDestination())
	backupDestinationCmd.AddCommand(backupController.CreateDestination())
	backupCmd.AddCommand(backupDestinationCmd)

	var backupJobCmd = &cobra.Command{
		Use:   "job",
		Short: "BackupJobManagement",
	}
	backupJobCmd.AddCommand(backupController.ReadJob())
	backupJobCmd.AddCommand(backupController.CreateJob())
	backupCmd.AddCommand(backupJobCmd)

	var backupTaskCmd = &cobra.Command{
		Use:   "task",
		Short: "BackupTaskManagement",
	}
	backupTaskCmd.AddCommand(backupController.ReadTask())
	backupCmd.AddCommand(backupTaskCmd)

	rootCmd.AddCommand(backupCmd)
}

func (router *Router) containerRoutes() {
	var containerCmd = &cobra.Command{
		Use:   "container",
		Short: "ContainerManagement",
	}

	containerController := cliController.NewContainerController(
		router.persistentDbSvc, router.trailDbSvc,
	)
	containerCmd.AddCommand(containerController.Read())
	containerCmd.AddCommand(containerController.Create())
	containerCmd.AddCommand(containerController.Update())
	containerCmd.AddCommand(containerController.Delete())

	var containerProfileCmd = &cobra.Command{
		Use:   "profile",
		Short: "ContainerProfileManagement",
	}

	containerProfileController := cliController.NewContainerProfileController(
		router.persistentDbSvc, router.trailDbSvc,
	)
	containerProfileCmd.AddCommand(containerProfileController.Read())
	containerProfileCmd.AddCommand(containerProfileController.Create())
	containerProfileCmd.AddCommand(containerProfileController.Update())
	containerProfileCmd.AddCommand(containerProfileController.Delete())

	var containerRegistryCmd = &cobra.Command{
		Use:   "registry",
		Short: "ContainerRegistryManagement",
	}

	containerRegistryController := cliController.NewContainerRegistryController(
		router.persistentDbSvc,
	)
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

func (router *Router) mappingRoutes() {
	var mappingCmd = &cobra.Command{
		Use:   "mapping",
		Short: "MappingManagement",
	}

	mappingController := cliController.NewMappingController(
		router.persistentDbSvc, router.trailDbSvc,
	)
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

	o11yController := cliController.NewO11yController(router.transientDbSvc)
	o11yCmd.AddCommand(o11yController.ReadOverview())
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
			fmt.Println("Infinite Ez " + infraEnvs.InfiniteEzVersion)
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
	router.backupRoutes()
	router.containerRoutes()
	router.mappingRoutes()
	router.o11yRoutes()
	router.scheduledTaskRoutes()
	router.systemRoutes()
}
