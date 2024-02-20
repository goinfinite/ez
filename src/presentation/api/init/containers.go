package apiInit

import (
	"github.com/speedianet/control/src/domain/useCase"
	"github.com/speedianet/control/src/infra"
	"github.com/speedianet/control/src/infra/db"
)

func BootContainers(persistentDbSvc *db.PersistentDatabaseService) {
	containerQueryRepo := infra.NewContainerQueryRepo(persistentDbSvc)
	containerCmdRepo := infra.NewContainerCmdRepo(persistentDbSvc)

	useCase.BootContainers(
		containerQueryRepo,
		containerCmdRepo,
	)
}
