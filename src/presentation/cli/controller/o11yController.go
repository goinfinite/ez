package cliController

import (
	"github.com/goinfinite/fleet/src/domain/useCase"
	"github.com/goinfinite/fleet/src/infra"
	cliHelper "github.com/goinfinite/fleet/src/presentation/cli/helper"
	"github.com/spf13/cobra"
)

func GetO11yOverviewController() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "overview",
		Short: "GetOverview",
		Run: func(cmd *cobra.Command, args []string) {
			o11yQueryRepo := infra.O11yQueryRepo{}
			o11yOverview, err := useCase.GetO11yOverview(o11yQueryRepo)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, o11yOverview)
		},
	}

	return cmd
}
