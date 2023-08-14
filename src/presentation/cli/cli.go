package cli

import (
	"fmt"
	"os"
	"path/filepath"

	cliMiddleware "github.com/speedianet/sfm/src/presentation/cli/middleware"
	"github.com/speedianet/sfm/src/presentation/shared"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   filepath.Base(os.Args[0]),
	Short: "Speedia AppManager CLI",
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
}

func CliInit() {
	defer cliMiddleware.PanicHandler()
	cliMiddleware.PreventRootless()

	shared.CheckEnvs()

	registerCliRoutes()

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
