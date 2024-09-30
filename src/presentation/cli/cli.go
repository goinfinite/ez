package cli

import (
	"fmt"
	"os"
	"path/filepath"

	cliController "github.com/goinfinite/ez/src/presentation/cli/controller"
	cliInit "github.com/goinfinite/ez/src/presentation/cli/init"
	cliMiddleware "github.com/goinfinite/ez/src/presentation/cli/middleware"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   filepath.Base(os.Args[0]),
	Short: "Speedia Control CLI",
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
}

func RunRootCmd() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(0)
}

func CliInit() {
	defer cliMiddleware.PanicHandler()
	cliMiddleware.PreventRootless()

	cliMiddleware.CheckEnvs()
	cliMiddleware.LogHandler()

	isSystemInstall := false
	if len(os.Args) > 1 {
		isSystemInstall = os.Args[1] == "sys-install"
	}

	if isSystemInstall {
		sysInstallController := cliController.SysInstallController{}
		rootCmd.AddCommand(sysInstallController.SysInstall())
		RunRootCmd()
	}

	persistentDbSvc := cliInit.PersistentDatabaseService()
	transientDbSvc := cliInit.TransientDatabaseService()
	trailDbSvc := cliInit.TrailDatabaseService()

	isLicenseRefresh := false
	if len(os.Args) > 2 {
		isLicenseRefresh = os.Args[1] == "license" && os.Args[2] == "refresh"
	}

	if !isLicenseRefresh {
		cliMiddleware.InvalidLicenseBlocker(persistentDbSvc, transientDbSvc)
	}

	cliMiddleware.SporadicLicenseValidation(persistentDbSvc, transientDbSvc)

	router := NewRouter(persistentDbSvc, transientDbSvc, trailDbSvc)
	router.RegisterRoutes()

	RunRootCmd()
}
