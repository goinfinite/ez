package cli

import (
	"fmt"
	"os"
	"path/filepath"

	cliMiddleware "github.com/speedianet/control/src/presentation/cli/middleware"
	sharedInit "github.com/speedianet/control/src/presentation/shared/init"
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

func CliInit() {
	defer cliMiddleware.PanicHandler()
	cliMiddleware.PreventRootless()

	sharedMiddleware.CheckEnvs()

	dbSvc := sharedInit.DatabaseService()

	cliMiddleware.SporadicLicenseValidation(dbSvc)

	cliRouter := NewCliRouter(dbSvc)
	cliRouter.RegisterRoutes()

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
