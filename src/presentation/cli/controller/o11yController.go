package cliController

import (
	"github.com/speedianet/control/src/domain/useCase"
	"github.com/speedianet/control/src/infra/db"
	o11yInfra "github.com/speedianet/control/src/infra/o11y"
	cliHelper "github.com/speedianet/control/src/presentation/cli/helper"
	"github.com/spf13/cobra"
)

type O11yController struct {
	transientDbSvc *db.TransientDatabaseService
}

func NewO11yController(
	transientDbSvc *db.TransientDatabaseService,
) *O11yController {
	return &O11yController{
		transientDbSvc: transientDbSvc,
	}
}

func (repo *O11yController) GetO11yOverview() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "overview",
		Short: "GetOverview",
		Run: func(cmd *cobra.Command, args []string) {
			o11yQueryRepo := o11yInfra.NewO11yQueryRepo(repo.transientDbSvc)
			o11yOverview, err := useCase.GetO11yOverview(o11yQueryRepo)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, o11yOverview)
		},
	}

	return cmd
}
