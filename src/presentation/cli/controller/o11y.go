package cliController

import (
	"github.com/speedianet/control/src/infra/db"
	cliHelper "github.com/speedianet/control/src/presentation/cli/helper"
	"github.com/speedianet/control/src/presentation/service"
	"github.com/spf13/cobra"
)

type O11yController struct {
	o11yService *service.O11yService
}

func NewO11yController(
	transientDbSvc *db.TransientDatabaseService,
) *O11yController {
	return &O11yController{o11yService: service.NewO11yService(transientDbSvc)}
}

func (controller *O11yController) ReadOverview() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "overview",
		Short: "ReadOverview",
		Run: func(cmd *cobra.Command, args []string) {
			cliHelper.ServiceResponseWrapper(controller.o11yService.ReadOverview())
		},
	}

	return cmd
}
