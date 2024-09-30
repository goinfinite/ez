package apiInit

import (
	"github.com/goinfinite/ez/src/domain/useCase"
	"github.com/goinfinite/ez/src/infra"
	"github.com/goinfinite/ez/src/infra/db"
)

func BootContainers(persistentDbSvc *db.PersistentDatabaseService) {
	containerQueryRepo := infra.NewContainerQueryRepo(persistentDbSvc)
	containerCmdRepo := infra.NewContainerCmdRepo(persistentDbSvc)

	useCase.BootContainers(
		containerQueryRepo,
		containerCmdRepo,
	)
}
