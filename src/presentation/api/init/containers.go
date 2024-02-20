package apiInit

import (
	"github.com/speedianet/control/src/domain/useCase"
	"github.com/speedianet/control/src/infra"
	"github.com/speedianet/control/src/infra/db"
)

func BootContainers(persistDbSvc *db.PersistentDatabaseService) {
	containerQueryRepo := infra.NewContainerQueryRepo(persistDbSvc)
	containerCmdRepo := infra.NewContainerCmdRepo(persistDbSvc)

	useCase.BootContainers(
		containerQueryRepo,
		containerCmdRepo,
	)
}
