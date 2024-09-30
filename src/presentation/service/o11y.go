package service

import (
	"github.com/goinfinite/ez/src/domain/useCase"
	"github.com/goinfinite/ez/src/infra/db"
	o11yInfra "github.com/goinfinite/ez/src/infra/o11y"
)

type O11yService struct {
	transientDbSvc *db.TransientDatabaseService
}

func NewO11yService(
	transientDbSvc *db.TransientDatabaseService,
) *O11yService {
	return &O11yService{transientDbSvc: transientDbSvc}
}

func (service *O11yService) ReadOverview() ServiceOutput {
	o11yQueryRepo := o11yInfra.NewO11yQueryRepo(service.transientDbSvc)
	o11yOverview, err := useCase.ReadO11yOverview(o11yQueryRepo)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, o11yOverview)
}
