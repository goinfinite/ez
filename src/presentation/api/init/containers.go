package apiInit

import (
	"github.com/speedianet/control/src/domain/useCase"
	"github.com/speedianet/control/src/infra"
	"github.com/speedianet/control/src/infra/db"
)

func BootContainers(dbSvc *db.DatabaseService) {
	containerQueryRepo := infra.NewContainerQueryRepo(dbSvc)
	containerCmdRepo := infra.NewContainerCmdRepo(dbSvc)

	useCase.BootContainers(
		containerQueryRepo,
		containerCmdRepo,
	)
}
