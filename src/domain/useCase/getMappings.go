package useCase

import (
	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/repository"
)

func GetMappings(
	mappingQueryRepo repository.MappingQueryRepo,
) ([]entity.Mapping, error) {
	return mappingQueryRepo.Get()
}
