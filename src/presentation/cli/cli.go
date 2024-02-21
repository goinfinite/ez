package cli

import (
	"fmt"
	"os"
	"path/filepath"

	cliController "github.com/speedianet/control/src/presentation/cli/controller"
	cliInit "github.com/speedianet/control/src/presentation/cli/init"
	cliMiddleware "github.com/speedianet/control/src/presentation/cli/middleware"
	sharedMiddleware "github.com/speedianet/control/src/presentation/shared/middleware"
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

	isSystemInstall := false
	if len(os.Args) > 1 {
		isSystemInstall = os.Args[1] == "sys-install"
	}

	if isSystemInstall {
		sysInstallController := cliController.SysInstallController{}
		rootCmd.AddCommand(sysInstallController.SysInstall())
		RunRootCmd()
	}

	sharedMiddleware.CheckEnvs()

	persistentDbSvc := cliInit.PersistentDatabaseService()
	transientDbSvc := cliInit.TransientDatabaseService()

	sharedMiddleware.InvalidLicenseBlocker(persistentDbSvc, transientDbSvc)
	cliMiddleware.SporadicLicenseValidation(persistentDbSvc, transientDbSvc)

	router := NewRouter(persistentDbSvc, transientDbSvc)
	router.RegisterRoutes()

	RunRootCmd()
}
