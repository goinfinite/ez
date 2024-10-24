package cliController

import (
	"github.com/goinfinite/ez/src/domain/useCase"
	"github.com/goinfinite/ez/src/infra"
	cliHelper "github.com/goinfinite/ez/src/presentation/cli/helper"
	"github.com/spf13/cobra"
)

type SysInstallController struct {
}

func (*SysInstallController) SysInstall() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sys-install",
		Short: "SysInstall",
		Run: func(cmd *cobra.Command, args []string) {
			sysInstallQueryRepo := infra.SysInstallQueryRepo{}
			sysInstallCmdRepo := infra.SysInstallCmdRepo{}
			serverCmdRepo := infra.ServerCmdRepo{}
			err := useCase.SysInstall(
				sysInstallQueryRepo,
				sysInstallCmdRepo,
				serverCmdRepo,
			)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, "InstallSuccess")
		},
	}

	return cmd
}
