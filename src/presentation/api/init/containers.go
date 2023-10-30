package apiInit

import (
	"github.com/goinfinite/fleet/src/domain/useCase"
	"github.com/goinfinite/fleet/src/infra"
	"github.com/goinfinite/fleet/src/infra/db"
)

func BootContainers(dbSvc *db.DatabaseService) {
	containerQueryRepo := infra.NewContainerQueryRepo(dbSvc)
	containerCmdRepo := infra.NewContainerCmdRepo(dbSvc)

	useCase.BootContainers(
		containerQueryRepo,
		containerCmdRepo,
	)
}
