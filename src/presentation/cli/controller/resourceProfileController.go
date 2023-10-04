package cliController

import (
	"github.com/speedianet/sfm/src/domain/useCase"
	"github.com/speedianet/sfm/src/infra"
	cliHelper "github.com/speedianet/sfm/src/presentation/cli/helper"
	"github.com/spf13/cobra"
)

func GetResourceProfilesController() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "GetResourceProfiles",
		Run: func(cmd *cobra.Command, args []string) {
			resourceProfileQueryRepo := infra.ResourceProfileQueryRepo{}
			resourceProfilesList, err := useCase.GetResourceProfiles(
				resourceProfileQueryRepo,
			)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, resourceProfilesList)
		},
	}

	return cmd
}
