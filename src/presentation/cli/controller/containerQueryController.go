package cliController

import (
	"github.com/speedianet/sfm/src/domain/useCase"
	"github.com/speedianet/sfm/src/infra"
	cliHelper "github.com/speedianet/sfm/src/presentation/cli/helper"
	"github.com/spf13/cobra"
)

func GetContainersController() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "GetContainers",
		Run: func(cmd *cobra.Command, args []string) {
			containerQueryRepo := infra.ContainerQueryRepo{}
			containersList, err := useCase.GetContainers(containerQueryRepo)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, containersList)
		},
	}

	return cmd
}
